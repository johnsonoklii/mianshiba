package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mianshiba/domain/agent/agent/resume"
	"mianshiba/domain/interview/dal"
	"mianshiba/domain/interview/dal/model"
	"mianshiba/domain/interview/repository"
	"mianshiba/infra/contract/storage"
	mjson "mianshiba/pkg/json"
	"mianshiba/pkg/pdf"
	"os"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// ResumeParseResult 简历解析结果
type ResumeParseResult struct {
	BasicInfo struct {
		Name      string `json:"name"`
		WorkYears string `json:"work_years"`
		Contact   string `json:"contact"`
	} `json:"basic_info"`
	Education []struct {
		School         string `json:"school"`
		Major          string `json:"major"`
		Degree         string `json:"degree"`
		GraduationYear string `json:"graduation_year"`
	} `json:"education"`
	WorkExperience []struct {
		Company          string `json:"company"`
		Position         string `json:"position"`
		Duration         string `json:"duration"`
		Responsibilities string `json:"responsibilities"`
	} `json:"work_experience"`
	TechStack                   []string      `json:"tech_stack"`
	Projects                    []interface{} `json:"projects"`
	Skills                      []string      `json:"skills"`
	Certifications              []string      `json:"certifications"`
	Strengths                   string        `json:"strengths"`
	PotentialWeaknesses         string        `json:"potential_weaknesses"`
	RecommendedDifficulty       string        `json:"recommended_difficulty"`
	InterviewFocusAreas         []string      `json:"interview_focus_areas"`
	SuggestedQuestionDirections []string      `json:"suggested_questions_directions"`
}

type ParseResumeRequest struct {
	FileKey  string // 文件唯一标识
	FileID   int64  // 文件ID
	UserID   int64  // 用户ID
	Filename string // 文件名
	Filetype string // 文件类型
	Filesize int64  // 文件大小
}

type ResumeAgentComponents struct {
	OSSClient  storage.Storage
	ResumeRepo repository.ResumeRepository
}

type ResumeAgent interface {
	ParseResumeAndSave(ctx context.Context, req *ParseResumeRequest) error
}

func NewResumeAgent(components *ResumeAgentComponents) ResumeAgent {
	return &resumeAgentImpl{
		ResumeAgentComponents: components,
	}
}

type resumeAgentImpl struct {
	*ResumeAgentComponents
}

// ParseResumeAndSave 调用简历解析智能体解析简历，并将结果保存到数据库
func (r *resumeAgentImpl) ParseResumeAndSave(ctx context.Context, req *ParseResumeRequest) error {
	// 添加 120 秒超时
	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// 创建简历解析智能体
	agent, err := resume.NewResumeParserAgent()
	if err != nil {
		log.Printf("[ParseResumeAndSave] 创建简历解析智能体失败: %v", err)
		return err
	}

	// 创建 runner
	runner := adk.NewRunner(timeoutCtx, adk.RunnerConfig{
		Agent: agent,
	})

	bytes, err := r.GetResumeObject(timeoutCtx, req)
	if err != nil {
		log.Printf("[ParseResumeAndSave] 获取简历文件内容失败: %v", err)
		return err
	}

	fmt.Printf("[ParseResumeAndSave] 原始文件内容: %s", bytes)

	resumeContent, err := pdf.TryParsePDFWithMultipleEncodings(bytes)
	if err != nil {
		log.Printf("[ParseResumeAndSave] 解析PDF内容失败: %v", err)
		return err
	}
	fmt.Printf("[ParseResumeAndSave] 解析后的简历内容: %s", resumeContent)

	// 构建查询消息，包含简历文件路径
	query := fmt.Sprintf(`【重要】请立即解析以下简历文件并提取关键信息：

简历内容：%s

【必须执行的步骤】：
1. 【第一步】从解析的简历文本中提取所有关键信息（姓名、工作年限、联系方式、教育背景、工作经历、技术栈、项目经验、技能、证书等）
2. 【第二步】分析候选人的背景特点和核心竞争力
3. 【第三步】生成面试建议和推荐难度

【重要提示】：
- 必须从简历内容中提取真实的信息，不要返回空数据
- 所有JSON字段都必须填充实际内容
- 只返回JSON格式，不要返回其他文本

请返回完整的 JSON 格式结果。`, resumeContent)

	// 创建用户消息
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: query,
	}

	messages := []adk.Message{
		userMsg,
	}

	// 运行智能体
	iter := runner.Run(timeoutCtx, messages)

	var lastMessage string
	for {
		select {
		case <-timeoutCtx.Done():
			log.Printf("[ParseResumeAndSave] 超时：等待智能体响应超过 120 秒")
			return fmt.Errorf("timeout waiting for resume parsing (120s)")
		default:
		}

		fmt.Printf("[ParseResumeAndSave]......\n")

		event, ok := iter.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			log.Printf("[ParseResumeAndSave] 错误: %v", event.Err)
			return fmt.Errorf("error during resume parsing: %w", event.Err)
		}

		// 收集最后一条消息
		if event.Output != nil && event.Output.MessageOutput != nil {
			lastMessage = event.Output.MessageOutput.Message.Content
		}
	}

	// 解析智能体响应
	if lastMessage == "" {
		log.Printf("[ParseResumeAndSave] 智能体未返回任何响应")
		return fmt.Errorf("agent returned empty response")
	}

	log.Printf("[ParseResumeAndSave] 智能体响应内容: %s", lastMessage)
	parseResult := parseResumeResponse(lastMessage)
	if parseResult == nil {
		log.Printf("[ParseResumeAndSave] 无法解析简历响应")
		return fmt.Errorf("failed to parse resume response")
	}

	// 验证解析结果是否有效（不能全是空数据）
	if !isValidResumeResult(parseResult) {
		log.Printf("[ParseResumeAndSave] 解析结果无效（全是空数据），请检查简历文件是否正确")
		return fmt.Errorf("resume parsing result is empty or invalid")
	}

	// 将解析结果保存到数据库
	err = r.saveResumeToDatabase(ctx, req, parseResult)
	if err != nil {
		log.Printf("[ParseResumeAndSave] 保存简历失败: %v", err)
		return fmt.Errorf("failed to save resume: %w", err)
	}

	log.Printf("[ParseResumeAndSave] 简历解析成功，简历ID: %d", req.FileID)
	return nil
}

func (r *resumeAgentImpl) GetResumeObject(ctx context.Context, req *ParseResumeRequest) (bytes []byte, err error) {
	// 从minio获取获取文件信息
	bytes, err = r.OSSClient.GetObject(ctx, req.FileKey)
	if err != nil {
		log.Printf("[GetResumeObject] 从minio获取文件失败: %v", err)
		return nil, fmt.Errorf("failed to get resume object from minio: %w", err)
	}

	return bytes, nil
}

func (r *resumeAgentImpl) DownloadResumeObject(ctx context.Context, fileKey, localPath string) error {

	// 2. 检查对象是否存在
	// _, err = r.OSSClient.StatObject(context.Background(), bucketName, objectPath, minio.StatObjectOptions{})
	// if err != nil {
	// 	return fmt.Errorf("检查对象存在失败: %v", err)
	// }

	// 3. 下载对象
	object, err := r.OSSClient.GetObject(context.Background(), fileKey)
	if err != nil {
		return fmt.Errorf("获取对象失败: %v", err)
	}

	// 4. 创建本地文件
	localFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %v", err)
	}
	defer localFile.Close()

	reader := bytes.NewReader(object)
	// 5. 复制数据
	if _, err = io.Copy(localFile, reader); err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}

	return nil
}

// parseResumeResponse 从智能体响应解析简历数据
func parseResumeResponse(agentResponse string) *ResumeParseResult {
	result := &ResumeParseResult{}

	// 尝试直接解析 JSON
	if err := json.Unmarshal([]byte(agentResponse), result); err != nil {
		log.Printf("[parseResumeResponse] 直接解析 JSON 失败: %v，尝试提取 JSON", err)
		// 尝试从文本中提取 JSON
		jsonStr := mjson.ExtractJSONFromResponse(agentResponse)
		if jsonStr == "" {
			log.Printf("[parseResumeResponse] 无法提取 JSON，原始响应: %s", agentResponse)
			return nil
		}

		log.Printf("[parseResumeResponse] 提取的 JSON: %s", jsonStr)
		// 尝试解析提取的 JSON
		if err := json.Unmarshal([]byte(jsonStr), result); err != nil {
			log.Printf("[parseResumeResponse] 解析提取的 JSON 失败: %v", err)
			return nil
		}
	}

	return result
}

// isValidResumeResult 检查解析结果是否有效（不能全是空数据）
func isValidResumeResult(result *ResumeParseResult) bool {
	if result == nil {
		return false
	}

	// 检查基本信息是否有内容
	if result.BasicInfo.Name != "" || result.BasicInfo.WorkYears != "" || result.BasicInfo.Contact != "" {
		return true
	}

	// 检查教育背景
	if len(result.Education) > 0 {
		return true
	}

	// 检查工作经历
	if len(result.WorkExperience) > 0 {
		return true
	}

	// 检查技术栈
	if len(result.TechStack) > 0 {
		return true
	}

	// 检查项目经验
	if len(result.Projects) > 0 {
		return true
	}

	// 检查技能
	if len(result.Skills) > 0 {
		return true
	}

	// 检查证书
	if len(result.Certifications) > 0 {
		return true
	}

	// 检查其他字段
	if result.Strengths != "" || result.PotentialWeaknesses != "" || result.RecommendedDifficulty != "" {
		return true
	}

	// 检查面试关注领域
	if len(result.InterviewFocusAreas) > 0 {
		return true
	}

	// 检查建议的提问方向
	if len(result.SuggestedQuestionDirections) > 0 {
		return true
	}

	// 如果所有字段都是空的，返回 false
	return false
}

func (r *resumeAgentImpl) saveResumeToDatabase(ctx context.Context, req *ParseResumeRequest, parseResult *ResumeParseResult) error {
	contentJSON, err := json.Marshal(parseResult)
	if err != nil {
		log.Printf("[saveResumeToDatabase] 序列化简历数据失败: %v", err)
		return fmt.Errorf("failed to marshal resume data: %w", err)
	}

	fileName := req.Filename
	log.Printf("[saveResumeToDatabase], 提取的文件名: %s", fileName)

	// 解析结果保存数据库
	err = r.ResumeRepo.UpdateResume(ctx, req.FileID, &model.Resume{
		LlmParseContent: string(contentJSON),
		Status:          dal.StatusParseSuccess,
	})
	if err != nil {
		log.Printf("[saveResumeToDatabase] 创建简历记录失败: %v", err)
		return fmt.Errorf("failed to create resume record: %w", err)
	}

	log.Printf("[saveResumeToDatabase] 简历记录已保存，ID: %d, 用户ID: %d, 文件Key: %s, 文件名: %s", req.FileID, req.UserID, req.FileKey, fileName)
	return nil
}

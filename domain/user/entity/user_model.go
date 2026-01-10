package entity

type UserModel struct {
	ID            int64
	UserID        int64  // 用户ID
	Name          string // 模型显示名称（用户维度唯一）
	ModelKey      string // 模型标识（doubao-1.5-vision-lite-250315）
	Protocol      string // 协议类型（openai/ark/claude/gemini/deepseek/ollama/qwen/ernie）
	BaseURL       string // API 基础地址
	ConfigJSON    string // 额外配置（如区域、访问密钥等）
	SecretHint    string // 密钥脱敏提示（如显示末尾4位）
	ProviderName  string // 提供商名称（如 OpenAI、Ark、DeepSeek）
	MetaID        int64  // 关联全局 model_meta.id（继承能力/图标）
	DefaultParams string // 默认参数（如 temperature、max_tokens）
	Scope         int32  // 使用范围（位掩码：1=智能体, 2=应用, 4=工作流）
	Status        int32  // 状态（0=禁用, 1=启用）
	IsDefault     int32  // 是否为默认（0=不是, 1=是）
	CreatedAt     int64  // 创建时间（毫秒时间戳）
	UpdatedAt     int64  // 更新时间（毫秒时间戳）
}

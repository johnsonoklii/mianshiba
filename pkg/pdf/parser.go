package pdf

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func IsPDF(data []byte) bool {
	// PDF文件以 "%PDF-" 开头
	if len(data) < 5 {
		return false
	}
	return string(data[:5]) == "%PDF-"
}

// 检查文本是否可读
func isReadable(text string) bool {
	if len(text) == 0 {
		return false
	}

	// 计算可打印字符比例
	printable := 0
	total := 0
	for _, r := range text {
		total++
		if r >= 32 && r <= 126 || // ASCII可打印
			r >= 0x4E00 && r <= 0x9FFF { // 常用汉字
			printable++
		}
	}

	ratio := float64(printable) / float64(total)
	fmt.Printf("可读字符比例: %.2f%%\n", ratio*100)

	return ratio > 0.3 // 30%以上可读字符
}

func TryParsePDFWithMultipleEncodings(pdfBytes []byte) (string, error) {
	// 尝试不同的编码
	encodings := []string{"", "gbk", "gb2312", "utf-8"}

	var lastErr error
	for _, encoding := range encodings {
		text, err := ParsePDFContentWithEncoding(pdfBytes, encoding)
		if err == nil && isReadable(text) {
			fmt.Printf("✅ 使用编码 %s 解析成功\n", encoding)
			return text, nil
		}
		lastErr = err
	}

	return "", fmt.Errorf("所有编码尝试失败: %v", lastErr)
}

// ParsePDFContentWithEncoding 支持编码转换的PDF解析
func ParsePDFContentWithEncoding(pdfBytes []byte, encoding string) (string, error) {
	if len(pdfBytes) == 0 {
		return "", fmt.Errorf("PDF字节数据为空")
	}

	// 创建PDF阅读器
	reader := bytes.NewReader(pdfBytes)
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return "", fmt.Errorf("创建PDF阅读器失败: %v", err)
	}

	// 检查加密
	isEncrypted, err := pdfReader.IsEncrypted()
	if err == nil && isEncrypted {
		success, err := pdfReader.Decrypt([]byte(""))
		if err != nil || !success {
			return "", fmt.Errorf("PDF已加密，无法解密")
		}
	}

	// 获取页数
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("获取页数失败: %v", err)
	}

	var allText strings.Builder
	allText.WriteString(fmt.Sprintf("PDF总页数: %d\n\n", numPages))

	// 逐页提取文本
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			continue // 跳过错误页
		}

		extractor, err := extractor.New(page)
		if err != nil {
			continue
		}

		text, err := extractor.ExtractText()
		if err != nil {
			continue
		}

		// 编码转换
		if encoding != "" {
			text, err = convertEncoding(text, encoding)
			if err != nil {
				fmt.Printf("编码转换失败(第%d页): %v\n", pageNum, err)
			}
		}

		allText.WriteString(fmt.Sprintf("=== 第 %d 页 ===\n", pageNum))
		allText.WriteString(text)
		allText.WriteString("\n\n")
	}

	return allText.String(), nil
}

// 编码转换
func convertEncoding(text, encoding string) (string, error) {
	switch strings.ToLower(encoding) {
	case "gbk", "gb2312":
		reader := transform.NewReader(
			strings.NewReader(text),
			simplifiedchinese.GBK.NewDecoder(),
		)
		decoded, err := io.ReadAll(reader)
		if err != nil {
			return text, err
		}
		return string(decoded), nil

	case "utf-8":
		return text, nil // 已经是UTF-8

	default:
		return text, nil
	}
}

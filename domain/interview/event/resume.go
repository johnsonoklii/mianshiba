package event

import "time"

// ResumeParseEvent 简历解析事件
type ResumeParseEvent struct {
	FileKey   string    `json:"file_key"`   // 文件唯一标识
	FileID    int64     `json:"file_id"`    // 文件ID
	UserID    int64     `json:"user_id"`    // 用户ID
	Filename  string    `json:"filename"`   // 文件名
	Filetype  string    `json:"filetype"`   // 文件类型
	Filesize  int64     `json:"filesize"`   // 文件大小
	CreatedAt time.Time `json:"created_at"` // 事件创建时间
}

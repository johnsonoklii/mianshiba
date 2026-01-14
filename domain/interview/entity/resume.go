package entity

type Resume struct {
	ID          int64  // 主键ID
	UserID      int64  // 用户ID
	FileKey     string // 对象存储中的文件唯一标识
	Filename    string // 原始文件名
	Filetype    string // 文件类型，如 pdf/docx
	Filesize    int64  // 文件大小（字节）
	Status      int32  // 简历状态：1已上传 2解析中 3已解析 4已删除 5失败
	ParseStatus int32  // 解析状态：0未开始 1解析中 2成功 3失败
	ParseError  string // 解析失败原因摘要
	UploadAt    int64  // 更新时间
}

// ResumeMsg Kafka消息结构
type ResumeMsg struct {
	FileKey  string `json:"file_key"`
	FileID   int64  `json:"file_id"`
	Filename string `json:"filename"`
	Filetype string `json:"filetype"`
	Filesize int64  `json:"filesize"`
	UserID   int64  `json:"user_id"`
}

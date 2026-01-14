namespace go interview

// ==================== 1. 简历上传相关 ====================

// 获取简历上传URL请求
struct ResumeUploadUrlRequest {
    1: string filename (api.query="filename")    // 文件名
    2: string filetype (api.query="filetype")    // 文件类型
}

// 获取简历上传URL响应
struct ResumeUploadUrlResponse {
    1: required string upload_url                          // 上传URL
    2: required string file_key                            // 文件标识
    3: required i64 file_id                                // 文件ID
    
    253: required i32 code
    254: required string msg
}

// 简历元信息请求（上传完成后发送）
struct ResumeMetaInfoRequest {
    1: required string file_key (api.form="file_key")      // 文件标识
    2: required i64 file_id (api.form="file_id")           // 文件ID
    3: required string filename (api.form="filename")      // 文件名
    4: required string filetype (api.form="filetype")      // 文件类型
    5: required i64 filesize (api.form="filesize")         // 文件大小（字节）
}

// 简历元信息响应
struct ResumeMetaInfoResponse {
    1: required ResumeInfo data
    
    253: required i32 code
    254: required string msg
}

// ==================== 2. 简历下载相关 ====================

// 获取简历下载URL请求
struct ResumeDownloadUrlRequest {
    1: required string file_key (api.form="file_key")      // 文件标识
}

// 获取简历下载URL响应
struct ResumeDownloadUrlResponse {
    1: required string download_url                        // 下载URL
    
    253: required i32 code
    254: required string msg
}

// ==================== 3. 简历删除相关 ====================

// 获取简历删除URL请求
struct ResumeDeleteUrlRequest {
    1: required string file_key (api.form="file_key")      // 文件标识
}

// 获取简历删除URL响应
struct ResumeDeleteUrlResponse {
    1: required string delete_url                          // 删除URL
    
    253: required i32 code
    254: required string msg
}

// 记录简历删除信息请求
struct ResumeDeleteInfoRequest {
    1: required string file_key (api.form="file_key")      // 文件标识
    2: optional string reason (api.form="reason")          // 删除原因
}

// 记录简历删除信息响应
struct ResumeDeleteInfoResponse {
    253: required i32 code
    254: required string msg
}

// ==================== 4. 简历信息相关 ====================

// 简历信息结构体
struct ResumeInfo {
    1: required i64 id                                     // 简历ID
    2: required string file_key                            // 文件标识
    3: required string filename                            // 文件名
    4: required string filetype                            // 文件类型
    5: required i64 filesize                               // 文件大小（字节）
    6: required i64 upload_at                              // 上传时间戳
    7: required i32 status                                 // 状态（uploading, uploaded, deleted）
    8: required i64 user_id                                // 用户ID
}

// 获取简历列表请求
struct ResumeListRequest {
    1: optional string status (api.query="status")         // 状态筛选
    2: optional string keyword (api.query="keyword")       // 关键词搜索
    3: optional i32 page (api.query="page", api.vd="$>=1") // 页码，默认 1
    4: optional i32 size (api.query="size", api.vd="$>=1&&$<=100")  // 每页数量，默认 20
}

// 获取简历列表响应
struct ResumeListResponse {
    1: required list<ResumeInfo> list
    2: required i64 total
    3: required i32 page
    4: required i32 size
    
    253: required i32 code
    254: required string msg
}

// 获取简历详情请求
struct ResumeDetailRequest {
    1: required string file_key (api.path="file_key")      // 文件标识
}

// 获取简历详情响应
struct ResumeDetailResponse {
    1: required ResumeInfo data
    
    253: required i32 code
    254: required string msg
}

// ==================== 5. 基础响应结构体 ====================

struct EmptyRequest {}

struct EmptyResponse {}

struct BaseResponse {
    253: required i32 code
    254: required string msg
}

// ==================== 6. 服务定义 ====================

// 面试服务定义
service InterviewService {
    // 1. 获取简历上传URL
    ResumeUploadUrlResponse GetResumeUploadUrl(1: ResumeUploadUrlRequest request) (
        api.get="/api/interview/resume/upload/url",
        api.category="interview",
        api.gen_path="interview"
    )

    // 2. 保存简历元信息
    ResumeMetaInfoResponse SaveResumeMetaInfo(1: ResumeMetaInfoRequest request) (
        api.post="/api/interview/resume/meta/save",
        api.category="interview",
        api.gen_path="interview"
    )

    // 3. 获取简历下载URL
    ResumeDownloadUrlResponse GetResumeDownloadUrl(1: ResumeDownloadUrlRequest request) (
        api.get="/api/interview/resume/download/url",
        api.category="interview",
        api.gen_path="interview"
    )

    // 4. 获取简历删除URL
    ResumeDeleteUrlResponse GetResumeDeleteUrl(1: ResumeDeleteUrlRequest request) (
        api.get="/api/interview/resume/delete/url",
        api.category="interview",
        api.gen_path="interview"
    )

    // 5. 记录简历删除信息
    ResumeDeleteInfoResponse RecordResumeDeleteInfo(1: ResumeDeleteInfoRequest request) (
        api.post="/api/interview/resume/delete/record",
        api.category="interview",
        api.gen_path="interview"
    )

    // 6. 获取简历列表
    ResumeListResponse GetResumeList(1: ResumeListRequest request) (
        api.get="/api/interview/resume/list",
        api.category="interview",
        api.gen_path="interview"
    )

    // 7. 获取简历详情
    ResumeDetailResponse GetResumeDetail(1: ResumeDetailRequest request) (
        api.get="/api/interview/resume/detail/:file_key",
        api.category="interview",
        api.gen_path="interview"
    )
}
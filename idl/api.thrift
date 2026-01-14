include "./user/user.thrift"
include "./interview/interview.thrift"
include "./app/developer_api.thrift"

namespace go mianshiba

service UserService extends user.UserService {}
service InterviewService extends interview.InterviewService {}
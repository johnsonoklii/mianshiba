package errno

// User: 700 000 000 ~ 700 999 999
const (
	ErrUserAuthenticationFailed = 700012006

	ErrUserEmailAlreadyExistCode      = 700000001
	ErrUserUniqueNameAlreadyExistCode = 700000002
	ErrUserInfoInvalidateCode         = 700000003
	ErrUserSessionInvalidateCode      = 700000004
	ErrUserResourceNotFound           = 700000005
	ErrUserInvalidParamCode           = 700000006
	ErrUserPermissionCode             = 700000007
	ErrNotAllowedRegisterCode         = 700000008
)

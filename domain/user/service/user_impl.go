package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"mianshiba/domain/user/dal/model"
	"mianshiba/domain/user/entity"
	"mianshiba/domain/user/repository"
	"mianshiba/infra/contract/cache"
	"mianshiba/infra/contract/idgen"
	"mianshiba/pkg/errorx"
	"mianshiba/pkg/jwt"
	"mianshiba/pkg/logs"
	"mianshiba/types/consts"
	"mianshiba/types/errno"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

type UserComponents struct {
	CacheCli cache.Cmdable
	IDGen    idgen.IDGenerator
	UserRepo repository.UserRepository
}

func NewUserDomain(ctx context.Context, c *UserComponents) User {
	return &userImpl{
		UserComponents: c,
	}
}

type userImpl struct {
	*UserComponents
}

func (u *userImpl) Create(ctx context.Context, req *CreateUserRequest) (user *entity.User, err error) {
	exist, err := u.UserRepo.CheckEmailExist(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if exist {
		return nil, errorx.New(errno.ErrUserEmailAlreadyExistCode, errorx.KV("email", req.Email))
	}

	// Hashing passwords using the Argon2id algorithm
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userName := req.UserName
	if userName == "" {
		userName = strings.Split(req.Email, "@")[0]
	}

	userID, err := u.IDGen.GenID(ctx)
	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	newUser := &model.User{
		ID:        userID,
		Username:  userName,
		Email:     req.Email,
		Password:  hashedPassword,
		Avatar:    "", // TODO
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Deleted:   false,
	}

	err = u.UserRepo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("insert user failed: %w", err)
	}

	return userPo2Do(newUser), nil
}

func (u *userImpl) Login(ctx context.Context, email, password string) (user *entity.User, jwtToken string, err error) {
	userModel, exist, err := u.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	if !exist {
		return nil, "", errorx.New(errno.ErrUserInfoInvalidateCode)
	}

	// Verify the password using the Argon2id algorithm
	valid, err := verifyPassword(password, userModel.Password)
	if err != nil {
		return nil, "", err
	}
	if !valid {
		return nil, "", errorx.New(errno.ErrUserInfoInvalidateCode)
	}

	jwtToken, err = jwt.GenerateToken(userModel.ID, userModel.Username, userModel.Role)
	if err != nil {
		return nil, "", err
	}

	err = u.SavejwtToken(ctx, jwtToken, userModel.ID)
	if err != nil {
		return nil, "", err
	}

	return userPo2Do(userModel), jwtToken, nil
}

func (u *userImpl) Logout(ctx context.Context, userID int64) (err error) {
	err = u.DeletejwtToken(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *userImpl) ForgotPassword(ctx context.Context, email string) (err error) {
	userModel, exist, err := u.UserRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if !exist {
		return nil
	}

	resetToken, err := generateResetToken()
	if err != nil {
		return fmt.Errorf("generate reset token error: %w", err)
	}

	tokenKey := fmt.Sprintf("reset_password:%s", resetToken)
	err = u.CacheCli.Set(ctx, tokenKey, userModel.ID, time.Hour).Err()
	if err != nil {
		return fmt.Errorf("store reset token error: %w", err)
	}

	// TODO: Send reset password email to the user
	// This would typically involve generating a reset URL with the token
	// and sending it via email using an email service

	logs.Infof("[ForgotPassword] Reset token generated for user ID %d: %s", userModel.ID, resetToken)

	return nil
}

func (u *userImpl) ResetPassword(ctx context.Context, token, password string) (err error) {
	tokenKey := fmt.Sprintf("reset_password:%s", token)
	userID, err := u.CacheCli.Get(ctx, tokenKey).Int64()
	if err != nil {
		return errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "Invalid or expired reset token"))
	}

	err = u.CacheCli.Del(ctx, tokenKey).Err()
	if err != nil {
		logs.Errorf("[ResetPassword] Failed to delete reset token: %v", err)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password error: %w", err)
	}

	userModel, err := u.UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by ID error: %w", err)
	}

	if userModel == nil {
		return errorx.New(errno.ErrUserInvalidParamCode, errorx.KV("msg", "User not found"))
	}

	err = u.UserRepo.UpdatePassword(ctx, userModel.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("update password error: %w", err)
	}

	logs.Infof("[ResetPassword] Password reset successful for user ID %d", userID)

	return nil
}

func (u *userImpl) GetUser(ctx context.Context, userID int64) (user *entity.User, err error) {
	return nil, nil
}

func (u *userImpl) GetUserList(ctx context.Context) (userList []*entity.User, err error) {
	return nil, nil
}

func (u *userImpl) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (err error) {
	return nil
}

func (u *userImpl) SavejwtToken(ctx context.Context, jwtToken string, userID int64) error {
	return u.CacheCli.Set(ctx, consts.JwtCacheKey(userID), jwtToken, consts.DefaultSessionDuration).Err()
}

func (u *userImpl) DeletejwtToken(ctx context.Context, userID int64) error {
	return u.CacheCli.Del(ctx, consts.JwtCacheKey(userID)).Err()
}

func (u *userImpl) GetJwtToken(ctx context.Context, userID int64) (jwtToken string, err error) {
	token, err := u.CacheCli.Get(ctx, consts.JwtCacheKey(userID)).Result()
	if err != nil {
		return "", err
	}
	return token, nil
}

// Argon2id parameter
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var defaultArgon2Params = &argon2Params{
	memory:      64 * 1024, // 64MB
	iterations:  3,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

// Hashing passwords using the Argon2id algorithm
func hashPassword(password string) (string, error) {
	p := defaultArgon2Params

	// Generate random salt values
	salt := make([]byte, p.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Calculate the hash value using the Argon2id algorithm
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// Encoding to base64 format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id $v = 19 $m = 65536, t = 3, p = 4 $< salt > $< hash >
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// generateResetToken generates a random reset token for password reset
func generateResetToken() (string, error) {
	// Generate a random 32-byte token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the token to base64 URL safe format
	resetToken := base64.URLEncoding.EncodeToString(tokenBytes)

	return resetToken, nil
}

func verifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash string
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var p argon2Params
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	p.saltLength = uint32(len(salt))

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	p.keyLength = uint32(len(decodedHash))

	// Calculate the hash value using the same parameters and salt values
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	return subtle.ConstantTimeCompare(decodedHash, computedHash) == 1, nil
}

func userPo2Do(model *model.User) *entity.User {
	return &entity.User{
		ID:        model.ID,
		Username:  model.Username,
		Email:     model.Email,
		Role:      model.Role,
		Avatar:    model.Avatar,
		CreatedAt: model.CreatedAt.UnixMilli(),
		UpdatedAt: model.UpdatedAt.UnixMilli(),
	}
}

package encrypt

import (
	"fmt"
	"mianshiba/conf"

	"github.com/coze-dev/coze-studio/backend/domain/plugin/encrypt"
)

const DefaultAPIKeySecret = "default-secret"

func getSecret() string {
	secret := conf.Global.API.APIKeySecret
	if secret == "" {
		secret = DefaultAPIKeySecret
	}
	return secret
}

// EncryptAPIKey 加密API密钥
// AES-CBC 解密实现
func EncryptAPIKey(apiKey string) (string, error) {
	secret := getSecret()
	fmt.Printf("secret: %s\n", secret)
	return encrypt.EncryptByAES([]byte(apiKey), secret)
}

// DecryptAPIKey 解密API密钥
// AES-CBC 解密实现
func DecryptAPIKey(encrypted string) (string, error) {
	secret := getSecret()
	data, err := encrypt.DecryptByAES(encrypted, secret)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

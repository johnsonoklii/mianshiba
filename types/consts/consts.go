package consts

import (
	"fmt"
	"time"
)

const (
	SessionMaxAgeSecond    = 30 * 24 * 60 * 60
	DefaultSessionDuration = SessionMaxAgeSecond * time.Second
	JwtTokenPrefix         = "jwt"
)

const (
	UseSSL      = "USE_SSL"
	SSLCertFile = "SSL_CERT_FILE"
	SSLKeyFile  = "SSL_KEY_FILE"

	HostKeyInCtx          = "HOST_KEY_IN_CTX"
	RequestSchemeKeyInCtx = "REQUEST_SCHEME_IN_CTX"

	SessionDataKeyInCtx = "SESSION_DATA_KEY_IN_CTX"
	OpenapiAuthKeyInCtx = "OPENAPI_AUTH_KEY_IN_CTX"

	APIConnectorID = int64(1024)

	StorageType        = "STORAGE_TYPE"
	MinIOAPIHost       = "MINIO_API_HOST"
	MinIOProxyEndpoint = "MINIO_PROXY_ENDPOINT"
	TOSBucketEndpoint  = "TOS_BUCKET_ENDPOINT"
	S3BucketEndpoint   = "S3_BUCKET_ENDPOINT"
)

const (
	ApplyUploadActionURI = "/api/common/upload/apply_upload_action"
	UploadURI            = "/api/common/upload"
)

func JwtCacheKey(userID int64) string {
	return fmt.Sprintf("%s:%d", JwtTokenPrefix, userID)
}

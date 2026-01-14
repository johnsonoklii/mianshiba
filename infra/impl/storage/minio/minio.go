package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mianshiba/conf"
	"mianshiba/infra/contract/storage"
	"mianshiba/infra/impl/storage/proxy"
	"mianshiba/pkg/ctxcache"
	"mianshiba/types/consts"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioClient struct {
	host            string
	client          *minio.Client
	accessKeyID     string
	secretAccessKey string
	bucketName      string
	endpoint        string
}

func getMinioClient(_ context.Context, endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*minioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client failed %v", err)
	}

	m := &minioClient{
		client:          client,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		bucketName:      bucketName,
		endpoint:        endpoint,
	}

	err = m.createBucketIfNeed(context.Background(), client, bucketName, "cn-north-1")
	if err != nil {
		return nil, fmt.Errorf("init minio client failed %v", err)
	}
	return m, nil
}

func New(ctx context.Context) (storage.Storage, error) {
	m, err := getMinioClient(ctx, conf.Global.MinIO.Endpoint, conf.Global.MinIO.AccessKey, conf.Global.MinIO.SecretKey, conf.Global.MinIO.Bucket, conf.Global.MinIO.UseSSL)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *minioClient) createBucketIfNeed(ctx context.Context, client *minio.Client, bucketName, region string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("check bucket %s exist failed %v", bucketName, err)
	}

	if exists {
		return nil
	}

	err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
	if err != nil {
		return fmt.Errorf("create bucket %s failed %v", bucketName, err)
	}

	return nil
}

func (m *minioClient) test() {
	ctx := context.Background()
	objectName := fmt.Sprintf("test-file-%d.txt", rand.Int())

	err := m.PutObject(ctx, objectName, []byte("hello content"), storage.WithContentType("text/plain"))
	if err != nil {
		log.Fatalf("upload file failed: %v", err)
	}
	log.Printf("upload file success")

	url, err := m.GetObjectUrl(ctx, objectName)
	if err != nil {
		log.Fatalf("get file url failed: %v", err)
	}

	log.Printf("get file url success, url: %s", url)

	content, err := m.GetObject(ctx, objectName)
	if err != nil {
		log.Fatalf("download file failed: %v", err)
	}

	log.Printf("download file success, content: %s", string(content))

	err = m.DeleteObject(ctx, objectName)
	if err != nil {
		log.Fatalf("delete object failed: %v", err)
	}

	log.Printf("delete object success")
}

func (m *minioClient) PutObject(ctx context.Context, objectKey string, content []byte, opts ...storage.PutOptFn) error {
	option := storage.PutOption{}
	for _, opt := range opts {
		opt(&option)
	}

	minioOpts := minio.PutObjectOptions{}
	if option.ContentType != nil {
		minioOpts.ContentType = *option.ContentType
	}

	if option.ContentEncoding != nil {
		minioOpts.ContentEncoding = *option.ContentEncoding
	}

	if option.ContentDisposition != nil {
		minioOpts.ContentDisposition = *option.ContentDisposition
	}

	if option.ContentLanguage != nil {
		minioOpts.ContentLanguage = *option.ContentLanguage
	}

	if option.Expires != nil {
		minioOpts.Expires = *option.Expires
	}

	_, err := m.client.PutObject(ctx, m.bucketName, objectKey,
		bytes.NewReader(content), int64(len(content)), minioOpts)
	if err != nil {
		return fmt.Errorf("PutObject failed: %v", err)
	}
	return nil
}

func (m *minioClient) GetObject(ctx context.Context, objectKey string) ([]byte, error) {
	obj, err := m.client.GetObject(ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("GetObject failed: %v", err)
	}
	defer obj.Close()
	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("ReadObject failed: %v", err)
	}
	return data, nil
}

func (m *minioClient) DeleteObject(ctx context.Context, objectKey string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("DeleteObject failed: %v", err)
	}
	return nil
}

func (m *minioClient) GetObjectUrl(ctx context.Context, objectKey string, opts ...storage.GetOptFn) (string, error) {
	option := storage.GetOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.Expire == 0 {
		option.Expire = 3600 * 24 * 7
	}

	reqParams := make(url.Values)
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectKey, time.Duration(option.Expire)*time.Second, reqParams)
	if err != nil {
		return "", fmt.Errorf("GetObjectUrl failed: %v", err)
	}

	// logs.CtxDebugf(ctx, "[GetObjectUrl] origin presignedURL.String = %s", presignedURL.String())
	ok, proxyURL := proxy.CheckIfNeedReplaceHost(ctx, presignedURL.String())
	if ok {
		return proxyURL, nil
	}

	return presignedURL.String(), nil
}

func (m *minioClient) GetUploadUrl(ctx context.Context, objectKey string, opts ...storage.GetOptFn) (string, error) {
	option := storage.GetOption{}
	for _, opt := range opts {
		opt(&option)
	}

	if option.Expire == 0 {
		option.Expire = 3600 // 默认1小时过期
	}

	presignedURL, err := m.client.PresignedPutObject(ctx, m.bucketName, objectKey, time.Duration(option.Expire)*time.Second)
	if err != nil {
		return "", fmt.Errorf("[GetUploadUrl] failed: %v", err.Error())
	}

	// logs.CtxDebugf(ctx, "[GetUploadUrl] origin presignedURL.String = %s", presignedURL.String())
	ok, proxyURL := proxy.CheckIfNeedReplaceHost(ctx, presignedURL.String())
	if ok {
		return proxyURL, nil
	}

	return presignedURL.String(), nil
}

func (m *minioClient) GetUploadHost(ctx context.Context) string {
	currentHost, ok := ctxcache.Get[string](ctx, consts.HostKeyInCtx)
	if !ok {
		return ""
	}
	return currentHost + consts.ApplyUploadActionURI
}

func (m *minioClient) GetServerID() string {
	return ""
}

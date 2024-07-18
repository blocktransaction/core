package oss

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

type OssClient struct {
	Bucket      *oss.Bucket
	basePath    string
	maxFileSize int
}

func NewOssClient(option *Option) *OssClient {
	client, err := oss.New(option.Endpoint, option.AccessKeyId, option.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	bucket, err := client.Bucket(option.BucketName)
	if err != nil {
		panic(err)
	}
	return &OssClient{
		Bucket:      bucket,
		basePath:    option.BasePath,
		maxFileSize: option.MaxFileSize,
	}
}

// UploadFilePath 上传本地文件路径
func (o *OssClient) UploadFilePath(toPath, filePath string) error {
	return o.Bucket.PutObjectFromFile(toPath, filePath)
}

// UploadFileStream 上传文件流
func (o *OssClient) UploadFileStream(toPath string, stream io.Reader) error {
	return o.Bucket.PutObject(toPath, stream)
}

// 上传图片
func (o *OssClient) UploadImage(filename, stream string) (string, error) {
	fileContentPosition := strings.Index(stream, ",")
	uploadBaseString := stream[fileContentPosition+1:]
	uploadString, err := base64.StdEncoding.DecodeString(uploadBaseString)
	if err != nil {
		return "", err
	}

	if len(uploadString) > o.maxFileSize {
		return "", errors.New("too large")
	}
	extName, isImage := IsImage(filename)
	if !isImage {
		return "", errors.New("not image")
	}

	filename = uuid.New().String() + extName
	err = o.Bucket.PutObject(o.basePath+filename, strings.NewReader(string(uploadString)))
	if err != nil {
		return "", err
	}

	return "https://" + o.Bucket.BucketName + "." + o.Bucket.Client.Config.Endpoint + "/" + o.basePath + filename, nil
}

// 获取文件类型
func GetFileContentType(out *os.File) (string, error) {
	// 只需要前 512 个字节就可以了
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

// 是否为图片
func IsImage(filename string) (string, bool) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp":
		return ext, true
	}
	return "", false
}

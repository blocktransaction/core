package oss

type Option struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	BasePath        string //基础路径
	MaxFileSize     int    //最大文件大小
}

// 选项
func NewOption(endpoint, accessKeyId, accessKeySecret, bucketName string, maxFileSize int) *Option {
	return &Option{
		Endpoint:        endpoint,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		BucketName:      bucketName,
		MaxFileSize:     maxFileSize,
	}
}

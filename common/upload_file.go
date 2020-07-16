package common

import "bytes"

// UploadFile 上传文件信息
type UploadFile struct {
	FileName string
	Data     *bytes.Buffer
}

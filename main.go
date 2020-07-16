package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/beyondblog/any-get/common"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const (
	tempDir = "./tmp"
)

func zipFiles(filename string, files []*common.UploadFile) error {

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, os.ModePerm)
	}

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err = appendFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func appendFileToZip(zipWriter *zip.Writer, file *common.UploadFile) error {

	writer, err := zipWriter.Create(file.FileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file.Data)
	return err
}

func uploadHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(200000)
	if err != nil {
		log.Fatal(err)
	}
	var uploadFiles []*common.UploadFile

	for file := range c.Request.MultipartForm.File {
		f, h, _ := c.Request.FormFile(file)

		log.Println("处理上传的文件: " + h.Filename)
		uploadFile := &common.UploadFile{FileName: h.Filename,
			Data: bytes.NewBuffer([]byte{}),
		}
		//保存文件到内存buffer 中
		io.Copy(uploadFile.Data, f)
		uploadFiles = append(uploadFiles, uploadFile)
	}
	name := uuid.Must(uuid.NewV4())
	filename := fmt.Sprintf("%s/%s.zip", tempDir, name)
	err = zipFiles(filename, uploadFiles)
	if err != nil {
		log.Println(err)
	}
	message := fmt.Sprintf("文件上传完毕! 下载地址: %s/download/%s.zip\n", baseUrl, name)
	c.String(http.StatusOK, message)
}

var baseUrl string

func main() {
	flag.StringVar(&baseUrl, "baseUrl", "http://127.0.0.1:8080", "baseUrl such as https://you_domian")
	flag.Parse()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "")
	})

	r.Static("/download", tempDir)
	r.POST("/", uploadHandler)
	r.Run("0.0.0.0:8080")
}

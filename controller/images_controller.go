package controller

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"

	"github.com/gin-gonic/gin"
)

// 上传方法
func GetImageStoragePath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dataDir := filepath.Join(dir, "img")
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err := os.MkdirAll(dataDir, 0755)
		if err != nil {
			panic(err)
		}
	}
	return dataDir
}

func storeImage(fileName string, reader io.Reader) (string, error) {
	imagePath := filepath.Join(GetImageStoragePath(), fileName)
	f, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if _, err := io.Copy(f, reader); err != nil {
		return "", err
	}
	return imagePath, nil
}

// helper to create a random suffix
func randSuffix(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// UploadFile 处理文件上传，字段名为 "file"，可选表单字段 content_id
func UploadFile(c *gin.Context) {
	log.Printf("=== UPLOAD DEBUG START ===")
	defer log.Printf("=== UPLOAD DEBUG END ===")

	// 打印所有请求头
	log.Printf("All Headers:")
	for name, values := range c.Request.Header {
		log.Printf("  %s: %v", name, values)
	}

	log.Printf("Method: %s", c.Request.Method)
	log.Printf("ContentLength: %d", c.Request.ContentLength)
	log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))

	// 方法1：直接尝试获取文件
	log.Printf("Attempting c.FormFile('file')...")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Printf("c.FormFile ERROR: %v", err)
		log.Printf("Error type: %T", err)
	} else {
		log.Printf("c.FormFile SUCCESS: %s (%d bytes)", fileHeader.Filename, fileHeader.Size)
		processUpload(c, fileHeader)
		return
	}

	// 方法2：手动解析 multipart form
	log.Printf("Attempting ParseMultipartForm...")
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("ParseMultipartForm ERROR: %v", err)
		log.Printf("Error type: %T", err)

		if err == http.ErrNotMultipart {
			constants.SendResponse(c, constants.UserBadRequest, gin.H{
				"error":                 "request is not multipart/form-data",
				"received_content_type": c.GetHeader("Content-Type"),
			})
			return
		}

		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error":  "failed to parse multipart form",
			"detail": err.Error(),
		})
		return
	}

	log.Printf("ParseMultipartForm SUCCESS")

	// 检查解析结果
	if c.Request.MultipartForm == nil {
		log.Printf("MultipartForm is nil after successful parse - this should not happen!")
		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error": "multipart form is nil",
		})
		return
	}

	log.Printf("MultipartForm Value fields count: %d", len(c.Request.MultipartForm.Value))
	log.Printf("MultipartForm File fields count: %d", len(c.Request.MultipartForm.File))

	// 打印所有字段
	for fieldName, values := range c.Request.MultipartForm.Value {
		log.Printf("Value field: %s = %v", fieldName, values)
	}

	for fieldName, files := range c.Request.MultipartForm.File {
		log.Printf("File field: %s, file count: %d", fieldName, len(files))
		for i, file := range files {
			log.Printf("  File[%d]: %s (%d bytes)", i, file.Filename, file.Size)
		}
	}

	// 查找文件字段
	if files, exists := c.Request.MultipartForm.File["file"]; exists && len(files) > 0 {
		fileHeader = files[0]
		log.Printf("Found file in MultipartForm.File: %s (%d bytes)", fileHeader.Filename, fileHeader.Size)
		processUpload(c, fileHeader)
		return
	}

	// 如果没有找到文件，列出所有可用的文件字段
	availableFiles := []string{}
	for name := range c.Request.MultipartForm.File {
		availableFiles = append(availableFiles, name)
	}

	log.Printf("No 'file' field found. Available file fields: %v", availableFiles)

	constants.SendResponse(c, constants.UserBadRequest, gin.H{
		"error":            "file field 'file' not found",
		"available_fields": availableFiles,
		"received_headers": map[string]string{
			"Content-Type": c.GetHeader("Content-Type"),
		},
		"hint": "use multipart/form-data with field named 'file' of type File",
	})
}

func processUpload(c *gin.Context, fileHeader *multipart.FileHeader) {
	log.Printf("Starting file processing for: %s", fileHeader.Filename)

	src, err := fileHeader.Open()
	if err != nil {
		log.Printf("Error opening file: %v", err)
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "open file failed"})
		return
	}
	defer src.Close()

	// 这里继续你原有的文件处理逻辑
	ext := filepath.Ext(fileHeader.Filename)
	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(ext)
	}

	storageName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randSuffix(6), ext)

	if _, err := storeImage(storageName, src); err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "save file failed", "detail": err.Error()})
		return
	}

	// 处理 content_id 等其他逻辑...
	var contentID *uint = nil
	if cid := c.PostForm("content_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			tmp := uint(parsed)
			contentID = &tmp
		}
	}

	filePath := filepath.Join(GetImageStoragePath(), storageName)
	fileURL := "/img/" + storageName
	fr := Model.FileRecord{
		OriginalName: fileHeader.Filename,
		StorageName:  storageName,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileHeader.Size,
		ContentID:    contentID,
		UploadTime:   time.Now(),
		Status:       "active",
	}

	if err := database.DB.Create(&fr).Error; err != nil {
		_ = os.Remove(filePath)
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "store record failed", "detail": err.Error()})
		return
	}

	log.Printf("File upload completed successfully: %s", storageName)
	constants.SendResponse(c, constants.Success, fr)
}

// DownloadImage 安全地提供已上传图片的下载/预览，支持 Range 请求
// 路径参数采用 :name （例如 /img/:name）
// 会校验路径不越界、存在性，并使用 http.ServeContent 提供高效的文件服务
func DownloadImage(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": "missing file name"})
		return
	}

	// 构造并校验路径，防止目录穿越
	storageDir := GetImageStoragePath()
	cleanName := filepath.Clean(name) // 清理路径段
	// disallow absolute paths
	if filepath.IsAbs(cleanName) {
		cleanName = filepath.Base(cleanName)
	}
	fullPath := filepath.Join(storageDir, cleanName)

	rel, err := filepath.Rel(storageDir, fullPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": "invalid file name"})
		return
	}

	fi, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "file not found"})
			return
		}
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "stat file failed", "detail": err.Error()})
		return
	}
	if fi.IsDir() {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": "not a file"})
		return
	}

	f, err := os.Open(fullPath)
	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": "open file failed", "detail": err.Error()})
		return
	}
	defer f.Close()

	// 推断 Content-Type
	ext := filepath.Ext(fullPath)
	ctype := mime.TypeByExtension(ext)
	if ctype == "" {
		// 尝试从文件读取前 512 字节判断
		var buf [512]byte
		n, _ := f.Read(buf[:])
		ctype = http.DetectContentType(buf[:n])
		// 重置文件偏移，ServeContent 会从头读
		_, _ = f.Seek(0, io.SeekStart)
	}

	// 设置响应头（Content-Disposition可按需改为attachment）
	c.Header("Content-Type", ctype)
	c.Header("Content-Length", strconv.FormatInt(fi.Size(), 10))
	c.Header("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))

	// 使用标准库提供的 ServeContent 支持部分请求（Range）和缓存语义
	http.ServeContent(c.Writer, c.Request, fi.Name(), fi.ModTime(), f)
}

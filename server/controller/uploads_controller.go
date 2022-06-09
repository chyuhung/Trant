package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func UploadsController(context *gin.Context) {
	// 获取url中的path参数值
	if urlpath := context.Param("path"); urlpath != "" {
		target := filepath.Join(getUploadsDir(), urlpath)
		context.Header("Content-Description", "File Transfer")
		context.Header("Content-Transfer-Encoding", "binary")
		context.Header("Content-Description", "attachment;filename="+urlpath)
		context.Header("Content-Type", "application/octet-stream")
		context.File(target)
	} else {
		context.Status(http.StatusNotFound)
	}
}

func getUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return uploads
}

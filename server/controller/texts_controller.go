package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func TextsController(context *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() // 获取程序运行绝对路径
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe) // 获取程序运行家目录绝对路径
		if err != nil {
			log.Fatal(err)
		}
		uploads := filepath.Join(dir, "uploads") //构造新目录绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  //创建目录，赋予权限
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()                   //构造新文件名称
		fullpath := path.Join("uploads", filename+".txt") //构造新文件相对路径
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		context.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}

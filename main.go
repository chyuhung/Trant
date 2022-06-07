package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func TextsController(context *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe)
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()
		uploads := filepath.Join(dir, "uploads")
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		context.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.POST("/api/v1/texts", TextsController)
		router.NoRoute(func(context *gin.Context) {
			path := context.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				context.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				context.Status(http.StatusNotFound)
			}
		})
		err := router.Run(":8080")
		if err != nil {
			fmt.Println("Gin 启动失败，程序退出！")
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	url := "http://127.0.0.1:8080/static/index.html"
	chromePath := "C:\\Program Files (x86)\\Microsoft\\EdgeCore\\102.0.1245.30\\msedge.exe"
	cmd := exec.Command(chromePath, "--app="+url)
	err := cmd.Start()
	if err != nil {
		fmt.Println("浏览器启动失败，请手动访问：", url)
		fmt.Println(err)
	}

	chanSignal := make(chan os.Signal, 1)
	signal.Notify(chanSignal, os.Interrupt)
	<-chanSignal
	err = cmd.Process.Kill()
	if err != nil {
		fmt.Println("浏览器关闭失败，请手动关闭！")
		fmt.Println(err)
	}
	//fmt.Println("hah")

}

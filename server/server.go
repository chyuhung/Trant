package server

import (
	"Trant/server/controller"
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

// Run给main调用，首字母大写
func Run() {
	port := "27149"
	//启动gin
	gin.SetMode(gin.DebugMode) // 设置gin为debug模式
	router := gin.Default()    // 创建一个gin引擎示例
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/static", http.FS(staticFiles))
	router.POST("/api/v1/texts", controller.TextsController)
	router.GET("/uploads/:path", controller.UploadsController)
	router.GET("/api/v1/addresses", controller.AddressesController)
	router.GET("/api/v1/qrcodes", controller.QrcodesController)
	router.POST("/api/v1/files", controller.FilesController)
	router.NoRoute(func(context *gin.Context) {
		urlpath := context.Request.URL.Path
		if strings.HasPrefix(urlpath, "/static/") {
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
	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Gin 启动失败，程序退出！")
		fmt.Println(err)
		os.Exit(1)
	}
}

package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
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

func AddressesController(context *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		// 检查ip地址，排除回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				//排除169开头地址
				if ip4[0] != 169 {
					result = append(result, ip4.String())
				}

			}
		}
	}
	context.JSON(http.StatusOK, gin.H{"addresses": result})
}
func GetUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return uploads
}

func UploadsController(context *gin.Context) {
	// 获取url中的path参数值
	if urlpath := context.Param("path"); urlpath != "" {
		target := filepath.Join(GetUploadsDir(), urlpath)
		context.Header("Content-Description", "File Transfer")
		context.Header("Content-Transfer-Encoding", "binary")
		context.Header("Content-Description", "attachment;filename="+urlpath)
		context.Header("Content-Type", "application/octet-stream")
		context.File(target)
	} else {
		context.Status(http.StatusNotFound)
	}
}

func QrcodesController(context *gin.Context) {
	if content := context.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}
		context.Data(http.StatusOK, "image/png", png)
	} else {
		context.Status(http.StatusBadRequest)
	}

}

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode) // 设置gin为debug模式
		router := gin.Default()    // 创建一个gin引擎示例
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.POST("/api/v1/texts", TextsController)
		router.GET("/uploads/:path", UploadsController)
		router.GET("/api/v1/addresses", AddressesController)
		router.GET("/api/v1/qrcodes", QrcodesController)
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

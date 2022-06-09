package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
	"log"
	"net/http"
)

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

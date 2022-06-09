package controller

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

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

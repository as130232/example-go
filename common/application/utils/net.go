package utils

import (
	"fmt"
	"linebot-go/common/global"
	"linebot-go/common/infrastructure/consts/contextKey"
	"linebot-go/common/infrastructure/consts/errorType"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-http-utils/headers"
)

func GetClientIp(c *gin.Context) string {
	clientIp := c.Request.Header.Get(headers.XForwardedFor) // >1.7有bug 代理取header X-Real-IP或X-Forwarded-IP
	if clientIp == "" {
		clientIp = c.ClientIP() // 取客戶端ip
	}
	if clientIp == "::1" {
		clientIp = "127.0.0.1"
	}

	return clientIp
}

func GetRemoteIp(c *gin.Context) string {
	remoteIp := c.RemoteIP() // 先回代理ip,沒有就取客戶端ip
	if remoteIp == "" {
		remoteIp = GetClientIp(c)
	}
	if remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp
}

func CheckIpAddressInCIDR(c *gin.Context, ipAddress string) {
	if global.ServerConfig.SysApiCidr == nil {
		panic(GenErrorMsg(http.StatusBadRequest, errorType.StatusForbidden, "access denied"))
	}

	if global.ServerConfig.SysApiCidr.IpNet == "" {
		panic(GenErrorMsg(http.StatusBadRequest, errorType.StatusForbidden, "access denied"))
	}

	_, allowedNet, _ := net.ParseCIDR(global.ServerConfig.SysApiCidr.IpNet)
	SetActionLog(c, contextKey.AllowedNet, fmt.Sprintf("%+v", allowedNet))

	ip := net.ParseIP(ipAddress)
	SetActionLog(c, contextKey.IpCIDR, fmt.Sprintf("%+v", ip))

	if !allowedNet.Contains(ip) {
		panic(GenErrorMsg(http.StatusBadRequest, errorType.StatusForbidden, "access denied"))
	}
}

func GetAllInterfaceNameAndIp() map[string][]net.IP {
	var interfaceNameAndIpMap = make(map[string][]net.IP)

	interfaces, _ := net.Interfaces()
	for _, itf := range interfaces {
		addrArr, err := itf.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrArr {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipArr := interfaceNameAndIpMap[itf.Name]
			if ipArr == nil {
				ipArr = make([]net.IP, 0)
			}
			ipArr = append(ipArr, ip)
			interfaceNameAndIpMap[itf.Name] = ipArr
		}
	}

	return interfaceNameAndIpMap
}

func GetHostName(interfaceNameAndIp map[string][]net.IP) string {
	hostName, _ := os.Hostname()
	if interfaceNameAndIp != nil && len(interfaceNameAndIp) > 0 {
		eth0AddrArr := interfaceNameAndIp["eth0"] // docker interface
		if eth0AddrArr != nil && len(eth0AddrArr) > 0 {
			for _, eth0Addr := range eth0AddrArr {
				if strings.HasPrefix(eth0Addr.String(), "10.") { // docker ip
					hostName = fmt.Sprintf("%s-%s", hostName, eth0Addr.String())
					break
				}
			}
		}
	}

	return hostName
}

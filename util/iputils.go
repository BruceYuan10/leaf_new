package util

import (
	"github.com/BruceYuan10/leaf_new/log"
	"net/http"
	"strings"
)

func GetIp(r *http.Request) string {
	ip := r.Header.Get("x-forwarded-for")
	log.Debug("x-forwarded-for ip: %v", ip)
	if len(ip) == 0 || "unknown" == strings.ToLower(ip) {
		ip = r.Header.Get("Proxy-Client-IP")
	}
	log.Debug("Proxy-Client-IP ip: %v", ip)
	if len(ip) == 0 || "unknown" == strings.ToLower(ip) {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	log.Debug("WL-Proxy-Client-IP ip: %v", ip)
	if len(ip) == 0 || "unknown" == strings.ToLower(ip) {
		ip = r.RemoteAddr
	}
	log.Debug("Get$Ip ip: %v", ip)
	if strings.Contains(ip, ",") {
		return strings.Split(ip, ",")[0]
	}
	log.Debug("GetIp$ip$result: %v", ip)
	return ip
}

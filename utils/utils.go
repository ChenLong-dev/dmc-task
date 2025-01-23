package utils

import (
	"encoding/json"
	"math/rand"
	"net"
	"strings"
	"time"
)

func GetUTCTime() time.Time {
	//t, _ := time.Parse(time.DateTime, "1970-01-01 00:00:00Z")
	return time.Now().UTC()
}

func GetUTCTime2(duration time.Duration) time.Time {
	return GetUTCTime().Add(duration)
}

func GetLocalTime() time.Time {
	return time.Now().Local()
}

func GetLocalTime2(duration time.Duration) time.Time {
	return GetLocalTime().Add(duration)
}

func GetTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

func GetTimestamp(t time.Time) int64 {
	return t.Unix()
}

func GetTimeStr(t time.Time) string {
	return t.Format(time.DateTime)
}

func GetRandInt(min, max int) int {
	// 使用当前时间作为随机数种子
	rand.NewSource(time.Now().UnixNano())
	// 生成随机数：min <= rand <= max
	return rand.Intn(max) + min
}

func MarshalByJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func UnmarshalByJson(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func GetLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	ips := make([]string, 0)
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // 不是IPv4地址
			}
			ips = append(ips, ip.String())
		}
	}
	return strings.Join(ips, "-")
}

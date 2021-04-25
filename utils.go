package yu

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
)

// Logf 固定格式打印信息
func Logf(fomart string, v ...interface{}) {
	f, l := GetCaller(2)
	args := append([]interface{}{f, l}, v...)
	log.Printf("[INFO %s:%d] - "+fomart, args...)
}

// Warnf 固定格式打印警告信息
func Warnf(fomart string, v ...interface{}) {
	f, l := GetCaller(2)
	args := append([]interface{}{f, l}, v...)
	log.Printf("[WARN %s:%d] - "+fomart, args...)
}

// Errf 固定格式打印错误信息
func Errf(fomart string, v ...interface{}) {
	f, l := GetCaller(2)
	args := append([]interface{}{f, l}, v...)
	log.Printf("[ERR %s:%d] - "+fomart, args...)
}

// LogErr 固定格式打印error
func LogErr(err error, msg string) {
	if err == nil {
		return
	}
	f, l := GetCaller(2)
	log.Printf("[ERR %s:%d] - %s, err: %v", f, l, msg, err)
}

// GetCaller 获取调用方法名
func GetCaller(d int) (string, int) {
	pc, _, line, ok := runtime.Caller(d)
	if !ok {
		return "unknown", -1
	}
	src := runtime.FuncForPC(pc).Name()
	p := strings.LastIndex(src, "/")
	return src[p+1:], line
}

// Trace 耗时检测
func Trace(msg string) func() {
	start := time.Now()
	return func() { log.Printf("[INFO] - %s exec %s", msg, time.Since(start)) }
}

// CreateFolder 创建path路径的文件夹
func CreateFolder(path string) error {
	return os.MkdirAll(path, os.ModeDir|os.ModePerm)
}

// HeadSeparator 确保path以/开头
func HeadSeparator(path string) string {
	if path == "" {
		return string(os.PathSeparator)
	}
	if strings.HasPrefix(path, string(os.PathSeparator)) {
		return path
	}
	return string(os.PathSeparator) + path
}

// TailSeparator 确保path以/结尾
func TailSeparator(path string) string {
	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return path
	}
	return path + string(os.PathSeparator)
}

// UUID 生成没有-的UUID
func UUID() string {
	u := uuid.NewV4()
	buf := make([]byte, 32)
	hex.Encode(buf, u[:])
	return string(buf)
}

// MD5 获取字符串的MD5
func MD5(src string) []byte {
	h := md5.New()
	h.Write([]byte(src))
	return h.Sum(nil)
}

const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandomString 生成n长度的字符串
func RandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GetRealIP 获取真实IP
func GetRealIP(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}
	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}
	return remoteAddr
}

// StrTime 将时间戳转成字符串
func StrTime(timestamp int64, timeLayout string) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp, 0).Format(timeLayout)
}

// UnixTime 将字符串转成时间戳
func UnixTime(timeStr, layout string) int64 {
	date, err := time.ParseInLocation(layout, timeStr, time.Local)
	if err != nil {
		return 0
	}
	return date.Unix()
}

// ToJsonStr 将结构体转成json字符串
func ToJsonStr(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

// LoadJSON 解析name.json
func LoadJSON(name string, data interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(data)
}

// SaveJSON 保存到name.json
func SaveJSON(name string, data interface{}) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(data)
}

// LoadOrSaveJSON 解析name.json，不存在保存
func LoadOrSaveJSON(name string, res interface{}, def func() interface{}) (err error) {
	if err = LoadJSON(name, res); err == nil {
		return
	}
	if !os.IsNotExist(err) {
		return
	}
	// 如果文件不存在，err不为nil需要手动设置为nil
	err = nil
	if def != nil {
		err = SaveJSON(name, def())
	}
	return
}

// LoadToml 解析name.toml
func LoadToml(name string, data interface{}) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return toml.NewDecoder(file).Decode(data)
}

// SaveToml 保存到name.toml
func SaveToml(name string, data interface{}) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()
	return toml.NewEncoder(file).Encode(data)
}

// LoadOrSaveToml 解析name.toml，不存在保存
func LoadOrSaveToml(name string, data interface{}, def func() interface{}) (err error) {
	if err = LoadToml(name, data); err == nil {
		return
	}
	if !os.IsNotExist(err) {
		return
	}
	// 如果文件不存在，err不为nil需要手动设置为nil
	err = nil
	if def != nil {
		err = SaveToml(name, def())
	}
	return
}

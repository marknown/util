package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
)

// 基础 Json 数据类型
type JsonBase struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Code    string `json:"code"`
}

// JsonExit 是 JsonResult 的结构体
type JsonExit struct {
	JsonBase
	Data interface{} `json:"data"`
}

// EnableProxy 代理到数据分析工具
func EnableProxy(proxy string) {
	os.Setenv("HTTP_PROXY", proxy)
}

// DisableProxy 清除代理
func DisableProxy() {
	os.Unsetenv("HTTP_PROXY")
}

// JsonResult 返回 JsonExit 的 JSON 字符串
func JsonResult(message string, success bool, data interface{}, jsonpCallback string, code string) string {
	var obj JsonExit
	obj.Message = message
	obj.Success = success
	obj.Data = data
	obj.Code = code
	if code == "" {
		if success {
			obj.Code = "ok"
		} else {
			obj.Code = "error"
		}
	}

	jsonBytes, err := json.Marshal(obj)
	if nil != err {
		log.Fatal(jsonBytes)
	}

	jsonString := string(jsonBytes)

	if "" != jsonpCallback {
		jsonString = jsonpCallback + "(" + jsonString + ")"
	}

	return jsonString
}

// JSONToStruct 把 Json 字符串解析成指定结构体
func JSONToStruct(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

// JSONPToStruct 把 Jsonp 字符串解析成指定结构体
func JSONPToStruct(s string, v interface{}) error {
	pos1 := strings.Index(s, "(")
	pos2 := strings.LastIndex(s, ")")

	if -1 < pos1 && -1 < pos2 {
		s = s[pos1+1 : pos2]
		return JSONToStruct(s, v)
	}

	return errors.New("input is not valid jsonp string")
}

// PrintMapString 打印 map[string]string 的值
func PrintMapString(val map[string]string) {
	for k, v := range val {
		fmt.Println(k + ":" + v)
	}
}

// PrintVal 打印任意类型的值
func PrintVal(val interface{}) {
	fmt.Printf("%#v\n", val)
}

// PrintError 打印错误信息
func PrintError(err error) {
	if nil != err {
		fmt.Println(err.Error())
	}
}

// PanicError 打印错误信息
func PanicError(err error) {
	if nil != err {
		panic(err.Error())
	}
}

// RandomN 返回一个[0,n)的随机数，不包含n
func RandomN(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

// URLEncode 把字符串编码
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// URLDecode 把字符串解码
func URLDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

// SubString 截取字符串
func SubString(s string, start int, length int) string {
	var r = []rune(s)
	len := len(r)

	end := start + length

	if start < 0 {
		start = 0
	}
	if end > len {
		end = len
	}

	return string(r[start:end])
}

// RegexReplace 根据正则来匹配
func RegexReplace(src string, replace string, regex string) string {
	re := regexp.MustCompile(regex)
	return re.ReplaceAllString(src, replace)
}

// GetWorkDirectory 获取当前工作目录
func GetWorkDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		PrintError(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// GetCurrentDirectory 获取可执行文件的当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		PrintError(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// GetTimestamp 返回自1970年开始的时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// GetTimestampString 返回自1970年开始的时间戳字符串
func GetTimestampString() string {
	return IntToStr(GetTimestamp())
}

// GetTimestampMs 返回自1970年开始的时间戳 毫秒级别
func GetTimestampMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// GetTimestampMsString 返回自1970年开始的时间戳字符串 毫秒级别
func GetTimestampMsString() string {
	return IntToStr(GetTimestampMs())
}

// GetTimeString 返回 2006-01-02 15:04:05 格式的时间字符串
func GetTimeString(t time.Time) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return t.In(cstZone).Format("2006-01-02 15:04:05")
}

// GetTimeMsString 返回 2006-01-02 15:04:05.000 格式的时间字符串(到毫秒)
func GetTimeMsString(ms int64) string {
	f := "2006-01-02 15:04:05.999"
	s := GetTimeFromatString(time.Unix(ms/1000, ms%1000*1000000), f)
	if 19 == len(s) {
		s += "."
	}
	s += strings.Repeat("0", len(f)-len(s))

	return s
}

// GetTimeMsStringFromTime 返回 2006-01-02 15:04:05.000 格式的时间字符串(到毫秒)
func GetTimeMsStringFromTime(t time.Time) string {
	f := "2006-01-02 15:04:05.999"
	s := GetTimeFromatString(t, f)
	if 19 == len(s) {
		s += "."
	}
	s += strings.Repeat("0", len(f)-len(s))

	return s
}

// GetTimeIntervalString 返回 2006-01-02 15:04:05 格式的时间字符串
func GetTimeIntervalString(t time.Time, intervalFormat string) string {
	interval, _ := time.ParseDuration(intervalFormat)
	newTime := t.Add(interval)
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return newTime.In(cstZone).Format("2006-01-02 15:04:05")
}

// GetTimeFromatString 返回指定格式的时间字符串
func GetTimeFromatString(t time.Time, format string) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return t.In(cstZone).Format(format)
}

// ParseTimeFromString 把字符串日期解析成 time.Time 类型
func ParseTimeFromString(t string) (time.Time, error) {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	return time.ParseInLocation("2006-01-02 15:04:05", t, cstZone)
}

// ToUpperFirst 把字符串首个字母转成大写
func ToUpperFirst(s string) string {
	if s[0] >= 97 && s[0] <= 122 {
		return fmt.Sprintf("%c", s[0]-32) + s[1:]
	}
	return s
}

// SliceContains 判断字符串 slice 里是否含有指定 value
func SliceContains(src []string, value string) bool {
	isContain := false
	for _, srcValue := range src {
		if srcValue == value {
			isContain = true
			break
		}
	}
	return isContain
}

// SliceContainsInt 判断int数组 slice 里是否含有指定 value
func SliceContainsInt(src []int, value int) bool {
	isContain := false
	for _, srcValue := range src {
		if srcValue == value {
			isContain = true
			break
		}
	}
	return isContain
}

// CeilFloat32 要进一法取整的值
func CeilFloat32(f float32) float32 {
	return float32(math.Ceil(float64(f)))
}

// FloorFloat32 舍去法取整
func FloorFloat32(f float32) float32 {
	return float32(math.Floor(float64(f)))
}

// CompareFloat32 按指定精度比较两个float32类型
// f1 > f2 返回 1
// f1 = f2 返回 0
// f1 < f2 返回 -1
func CompareFloat32(f1, f2 float32, precision int) int {
	if precision < 0 {
		precision = 0
	}
	min := math.Pow10(-precision)
	diff := math.Abs(float64(f1) - float64(f2))

	if math.Max(float64(f1), float64(f2)) == float64(f1) && diff > min {
		return 1
	} else if math.Max(float64(f1), float64(f2)) == float64(f2) && diff > min {
		return -1
	}

	return 0
}

// EncodingTo 字符串从 from 编码转到 to 编码
func EncodingTo(s string, from string, to string) string {
	Decoder := mahonia.NewDecoder(from)
	Encoder := mahonia.NewEncoder(to)
	return Encoder.ConvertString(Decoder.ConvertString(s))
}

// GBKToUTF 把 GBK 字符串转成 UTF8 字符串
func GBKToUTF(s string) string {
	return EncodingTo(s, "GBK", "UTF8")
}

// FileRead 读取文件内容
func FileRead(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if nil != err {
		return ""
	}

	return string(bytes)
}

// FileAppend 把字符串追加到指定文件
func FileAppend(fileName string, s string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	buf := []byte(s)
	fd.Write(buf)
	fd.Close()
}

// FileWrite 把字符串追写到指定文件，之前的内容会被清空
func FileWrite(fileName string, s string) {
	fd, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	buf := []byte(s)
	fd.Write(buf)
	fd.Close()
}

// FileExist 判断路径或者文件是否存在
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// StrToInt 字符串转成 int64
func StrToInt(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return int64(i)
}

// IntToStr int64 转成字符串
func IntToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Int16ToStr int 转成字符串
func Int16ToStr(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

// StrToFloat 字符串转成 float64
func StrToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return f
}

// FloatToStr float64 字符串转成字符串，prec为小数精度
func FloatToStr(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// commands 不同平台的打开命令
var commands = map[string]string{
	"windows": "cmd",
	"darwin":  "open",
	"linux":   "xdg-open",
}

// Open calls the OS default program for uri
func Open(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	if "windows" == runtime.GOOS {
		uri = " /c start " + uri
	}
	cmd := exec.Command(run, uri)
	err := cmd.Start()
	if nil != err {
		return err
	}

	// bytes, errOutput := cmd.Output()
	_, errOutput := cmd.Output()
	if nil != errOutput {
		return errOutput
	}

	// fmt.Println("Output:" + string(bytes))

	return nil
}

// Exec 执行系统命令，并返回输出
func Exec(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	err := cmd.Run()
	if nil != err {
		return "", err
	}

	return output.String(), nil
}

// ExecAsync 异步执行系统命令
func ExecAsync(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	err := cmd.Start()
	if nil != err {
		return err
	}

	return nil
}

// Md5 函数
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

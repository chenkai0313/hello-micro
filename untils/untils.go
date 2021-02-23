package untils

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/blinkbean/dingtalk"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// sha256加密
func Sha256(t string) string {
	h := sha256.New()
	h.Write([]byte(t))
	return hex.EncodeToString(h.Sum(nil))
}

//float64 保留小数
func Decimal(value float64, d string) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%."+d+"f", value), 64)
	return value
}

// 表单形式发送http的post请求
func HttpFormPost(url, data string) (string, error) {
	return HttpRequest("POST", url, data, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
}

// json形式发送http的post请求
func HttpJsonPost(url, data string) (string, error) {
	return HttpRequest("POST", url, data, map[string]string{"Content-Type": "application/json;charset=UTF-8"})
}

// 发送http请求
func HttpRequest(method string, url string, data string, headers map[string]string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	for key, content := range headers {
		req.Header.Set(key, content)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body[:]), nil
}

//时间戳转字符串
func TimestampToString(timeStamp int64) string {
	if timeStamp == 0 {
		return ""
	}
	timeTemplate := "2006-01-02 15:04:05"
	return time.Unix(timeStamp, 0).Format(timeTemplate)
}

//字符串时间转int64
func TimeStringToInt64(stringTime string) int64 {
	loc, _ := time.LoadLocation("Local")
	the_time, _ := time.ParseInLocation("2006-01-02 15:04:05", stringTime, loc)
	return the_time.Unix()
}

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	aaa := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, aaa[r.Intn(len(aaa))])
	}
	return string(result)
}

//生成随机数字
func GetRandomNumber(length int) string {
	str := "0123456789"
	aaa := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, aaa[r.Intn(len(aaa))])
	}
	return string(result)
}

//获取ip
func GetIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected  to the network fail")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

//ip string 转int
func StringIpToInt(ipstring string) int {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt int = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

//ip int 转string
func IpIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var len int = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < len; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[len-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < len; i++ {
		buffer.WriteString(ipSegs[i])
		if i < len-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}

//解压zip
func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//zip压缩文件
func Zip(srcFile string, destZip string) error {
	zipfile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
		// header.Name = path
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}

//判断是否是zip压缩包
func IsZip(Contenttype string) bool {
	s := strings.Split(Contenttype, "/")
	fileType := s[len(s)-1]
	if fileType != "zip" && fileType != "x-zip-compressed" {
		return false
	}
	return true
}

//获取文件类型
func GetFileType(Contenttype string) string {
	s := strings.Split(Contenttype, "/")
	fileType := s[len(s)-1]
	return fileType
}

//判断文件文件夹是否存在
func IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

//时间字符串转时间戳，空字符串转0
func TimeStrToInt64(timeStr string) int64 {
	var res int64
	if timeStr != "" {
		res = TimeStringToInt64(timeStr)
	} else {
		res = 0
	}
	return res
}

//字符串转int,出错返回0
func StrToInt64(str string) int64 {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		res = 0
	}
	return res
}

//叮叮通知
func Ding(msg string, token string) {
	//单个机器人有单位时间内消息条数的限制，如果有需要可以初始化多个token，发消息时随机发给其中一个机器人。
	var dingToken = []string{token}
	cli := dingtalk.InitDingTalk(dingToken, "")
	_ = cli.SendTextMessage(msg)
}

func InArrayInt32(item int32, items []int32) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func InArrayInt64(item int64, items []int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func InArrayString(item string, items []string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

//表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

//表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

//下载文件资源
func DownloadImg(imgUrl string, imgType string) (imgPath string, err error) {
	fileName := GetRandomString(16) + "." + imgType
	if imgType == "" {
		fileName = path.Base(imgUrl)
	}
	filePath := "./files/" + fileName
	out, err := os.Create(filePath)
	defer out.Close()
	resp, err := http.Get(imgUrl)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(out, bytes.NewReader(pix))
	if err != nil {
		return "", err
	}
	return filePath, nil
}

//base64转图片
func Base64ToImage(base64Byte []byte) (imgPath string, err error) {
	fileName := GetRandomString(16) + ".jpg"
	imgPathRes := "./files/" + fileName
	// 解压
	encodeString := base64.StdEncoding.EncodeToString(base64Byte)
	dist, _ := base64.StdEncoding.DecodeString(encodeString)
	// 写入新文件
	f, err := os.OpenFile(imgPathRes, os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer f.Close()
	f.Write(dist)
	if err != nil {
		return "", err
	}
	return imgPathRes, nil
}

//图片转base64
func ImagesToBase64(str_images string) ([]byte, error) {
	//读原图片
	ff, _ := os.Open(str_images)
	defer ff.Close()
	sourcebuffer := make([]byte, 500000)
	n, err := ff.Read(sourcebuffer)
	if err != nil {
		return nil, err
	}
	//base64压缩
	sourcestring := base64.StdEncoding.EncodeToString(sourcebuffer[:n])
	return []byte(sourcestring), nil
}

//ArrayDiff 模拟PHP array_diff函数  差集
func ArrayDiff(array1 []interface{}, othersParams ...[]interface{}) ([]interface{}, error) {
	if len(array1) == 0 {
		return []interface{}{}, nil
	}
	if len(array1) > 0 && len(othersParams) == 0 {
		return array1, nil
	}
	var tmp = make(map[interface{}]int, len(array1))
	for _, v := range array1 {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]interface{}, 0, len(tmp))
	for k, v := range tmp {
		if v == 1 {
			res = append(res, k)
		}
	}
	return res, nil
}

//ArrayIntersect 模拟PHP array_intersect函数  交集
func ArrayIntersect(array1 []interface{}, othersParams ...[]interface{}) ([]interface{}, error) {
	if len(array1) == 0 {
		return []interface{}{}, nil
	}
	if len(array1) > 0 && len(othersParams) == 0 {
		return array1, nil
	}
	var tmp = make(map[interface{}]int, len(array1))
	for _, v := range array1 {
		tmp[v] = 1
	}
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] != 0 {
				tmp[arg]++
			}
		}
	}
	var res = make([]interface{}, 0, len(tmp))
	for k, v := range tmp {
		if v > 1 {
			res = append(res, k)
		}
	}
	return res, nil
}

//获取excel 列数据
// filePath excel 路径
//sheetName
//cols 列名称
func ExcelValueLoc(filePath, sheetName string, cols []string) ([][]string, error) {
	xlsx, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows := xlsx.GetRows(sheetName)
	if len(rows) == 0 {
		return nil, nil
	}
	colIndex := make([]int, len(cols))

	// 获取每个col的所在序列号
	for index, row := range rows {
		if index == 0 {
			num := 0
			for _, col := range cols {
				for key, colCell := range row {
					if colCell == col {
						colIndex[num] = key + 1
						num++
					}
				}
			}
		}
	}

	//	对存在的量进行重新矫正，以解决初始变量长度问题
	res_len := 0
	for _, coli := range colIndex {
		if coli-1 >= 0 {
			res_len++
		}
	}

	// 获取数据
	res_data := make([][]string, len(rows)-1)
	res_index := 0
	for index, row := range rows {
		if index != 0 {
			data := make([]string, res_len)
			for i, colindex := range colIndex {
				for key, colCell := range row {
					if key == colindex-1 {
						data[i] = colCell
					}
				}
			}
			res_data[res_index] = data
			res_index++
		}
	}
	return res_data, nil
}

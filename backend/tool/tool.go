package tool

import (
	"errors"
	"github.com/gen2brain/go-unarr"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

//打开文件和读内容 利用io/ioutil
func ReadAll(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	//对内容的操作
	//ReadFile返回的是[]byte字节切片，要用string()方法转变成字符串
	//去除内容结尾的换行符
	str := strings.TrimRight(string(content), "\n")
	return str, nil
}

//文件写入 先清空再写入 利用ioutil
func WriteFast(filePath string, content string) error {
	dir, _ := path.Split(filePath)
	exist, err := IsFileExisted(dir)
	if err != nil {
		return err
	} else if exist == false {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//判断文件/文件夹是否存在
func IsFileExisted(path string) (bool, error) {
	//返回 true, nil = 存在
	//返回 false, nil = 不存在
	//返回 _, !nil = 位置错误，无法判断
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//利用HTTP Get请求获得数据json
func GetHttpData(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	//body, err := resp.Js
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	_ = resp.Body.Close()

	return string(data), nil
}

//下载文件 (下载地址，存放位置)
func DownloadFile(url string, location string) error {
	//利用HTTP下载文件并读取内容给data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		errorInfo := "http failed, check if file exists, HTTP Status Code:" + strconv.Itoa(resp.StatusCode)
		return errors.New(errorInfo)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	//确保下载位置存在
	_, fileName := path.Split(url)
	ok, err := IsFileExisted(location)
	if err != nil {
		return err
	} else if ok == false {
		err := os.Mkdir(location, os.ModePerm)
		if err != nil {
			return err
		}
	}
	//文件写入 先清空再写入 利用ioutil
	err = ioutil.WriteFile(location+"/"+fileName, data, 0666)
	if err != nil {
		return err
	} else {
		return nil
	}
}

//判断是不是non-ASCII
func IsNonASCII(str string) bool {
	re := regexp.MustCompile("[[:^ascii:]]")
	return re.MatchString(str)
	//var count int
	//for _, v := range str {
	//	if unicode.Is(unicode.Han, v) {
	//		count++
	//		break
	//	}
	//}
	//return count > 0
}

//解压zip 7z rar tar
func Decompress(from string, to string) error {
	a, err := unarr.NewArchive(from)
	if err != nil {
		return err
	}
	defer a.Close()

	_, err = a.Extract(to)
	if err != nil {
		return err
	}

	return nil
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//规格化路径
func FormatPath(s string) string {
	if strings.HasPrefix(s, "."){
		s = strings.Replace(s, ".", getCurrentDirectory(), 1)
	}
	s = strings.Replace(s, "\\", "/", -1)

	return strings.TrimRight(s, "\\")
}

//复制文件夹
func CopyDir(from string, to string) error {
	from = FormatPath(from)
	to = FormatPath(to)

	//确保目标路径存在，否则复制报错exit status 4
	exist, err := IsFileExisted(to)
	if err != nil {
		return err
	} else if exist == false {
		err := os.Mkdir(to, os.ModePerm)
		if err != nil {
			return err
		}
	}
	var out string
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("xcopy", from, to, "/I", "/E", "/Y", "/R")
	} else {
		cmd = exec.Command("cp", "-R", from, to)
	}

	//if runtime.GOOS == "windows" {
	//	out, err = Cmd("xcopy /I /E /Y " + strconv.Quote(from) + " " + strconv.Quote(to))
	//} else {
	//	out, err = Cmd("cp -R " + from + " " + to)
	//}
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println(out, err)
	}
	return err
}

//执行一次command指令
func Cmd(command string) (string, error) {
	var out []byte
	var err error
	if runtime.GOOS == "windows" {
		c := exec.Command("cmd.exe", "/c", command)
		out, err = c.CombinedOutput()
	} else {
		c := exec.Command("/bin/bash", "-c", command)
		out, err = c.CombinedOutput()
	}
	//cmd.Args = a
	return string(out), err
}
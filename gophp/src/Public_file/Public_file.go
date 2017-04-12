package Public_file

import (
	"fmt"
	//	"io"
	"os"
	"os/exec"
	"strings"
	//	"time"
	"io/ioutil"
	"path/filepath"

	"github.com/axgle/mahonia"
)

var (
	Cgi_ip   = "" //IP
	Cgi_port = "" //端口

)

func Cmdexec(cmd string, system string) string {
	//	defer Try_Err() //异常处理
	var c *exec.Cmd
	var data string
	if system == "windows" {
		argArray := strings.Split("/c "+cmd, " ")
		c = exec.Command("cmd", argArray...)
	} else {
		c = exec.Command("/bin/sh", "-c", cmd)
	}
	out, _ := c.Output()
	data = string(out)
	if system == "windows" {
		dec := mahonia.NewDecoder("gbk")
		data = dec.ConvertString(data)
	}
	return data
}

func PathExists(path string) (bool, string) { //判断文件是否存在
	_, err := os.Stat(path)
	if err == nil {
		return true, ""
	}
	if os.IsNotExist(err) {
		return false, "file does not exist"
	}
	return false, fmt.Sprintf("%s", err)
}

func ReadAll(path string) string { //读取文件
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	// fmt.Println(string(fd))
	return string(fd)
}

func Get_CurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	//fmt.Println("file:", file)
	path, _ := filepath.Abs(file)
	//fmt.Println("path:", path)
	splitstring := strings.Split(path, "\\")
	size := len(splitstring)
	splitstring = strings.Split(path, splitstring[size-1])
	//fmt.Println("path:", splitstring[0])
	return splitstring[0]
	//ret := strings.Replace(splitstring[0], "\\", "/", size-1)
	//return ret
}

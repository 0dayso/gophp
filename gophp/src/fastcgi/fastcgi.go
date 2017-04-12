package fastcgi

//执行PHP代码
import (
	"fmt"
	//"os/exec"

	"Public_file"
	//	"encoding/json"
	"net"
	//"reflect"
	//	"strconv"
	//"reflect"
	"net/http"
	"strings"
	"time"
	"unsafe"
	//	"github.com/axgle/mahonia"
)

func Run_cgi(ip, port string) { //启动PHP  CGI
	//	ip := "127.0.0.1"
	//	port := "9000"
	//Public_file.Cmdexec("taskkill /f /t /im GOPHP.exe", "windows")
	Public_file.Cmdexec("taskkill /f /t /im cmd.exe", "windows")
	Public_file.Cmdexec("taskkill /f /t /im php-cgi.exe", "windows")

	for { //死循环
		process_name := Public_file.Cmdexec("tasklist", "windows")
		if strings.Contains(process_name, "php-cgi.exe") == false {
			fmt.Printf("run cgi\n")
			ip_port := fmt.Sprintf(".\\php\\php-cgi.exe -b %s:%s -c php.ini", ip, port)
			Public_file.Cmdexec(ip_port, "windows")
		}
		time.Sleep(10 * time.Second)
	}
}

func rand_data(rand_data string) string { //去除前3行无效数据
	data_list := strings.Split(rand_data, "\n")
	//fmt.Println(data_list[3:])
	rand_datax := ""
	for index, value := range data_list {
		if index >= 3 {
			if len(value) >= 1 {
				//fmt.Printf("arr[%d]=%s \n", index, value)
				rand_datax += value + "\n"
			}
		}
	}
	return rand_datax //strings.Replace(rand_datax, " ", "", 1)
}

func BytesString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func B2S(bs []int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func Request(ip, port string, params map[string]string) (bool, string) { //请求数据
	//TCP链接
	addr := ip + ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return false, fmt.Sprintf("ResolveTCPAddr err:%s", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return false, fmt.Sprintf("DialTCP err:%s", err)
	}
	//=====================================================

	requestpkghead := []byte{1, 1, 0, 1, 0, 8, 0, 0}
	requestpkgbody := []byte{0, 1, 0, 0, 0, 0, 0, 0}
	paramsBody := []byte{}
	//	for k, v := range args {
	//		//log.Println(k, v)
	//		fmt.Println(k, v)
	//	}
	for k, v := range params {
		//fmt.Println(reflect.TypeOf(v))
		if v != "" {
			paramsBody = append(paramsBody, byte(len([]byte(k))))
			paramsBody = append(paramsBody, byte(len([]byte(v))))
			paramsBody = append(paramsBody, []byte(k)...)
			paramsBody = append(paramsBody, []byte(v)...)
		}
	}
	paramsHead := []byte{1, 4, 0, 1, 0, byte(len(paramsBody)), 0, 0}
	endpkghead := []byte{1, 3, 0, 1, 0, 0, 0, 0}
	request := append(requestpkghead, requestpkgbody...)
	request = append(request, paramsHead...)
	request = append(request, paramsBody...)
	request = append(request, endpkghead...)
	if _, err := conn.Write(request); err != nil {
		defer conn.Close()
		return false, fmt.Sprintf("Write err:%s", err)
	} else {
		response := make([]byte, 8192)
		if _, err := conn.Read(response); err != nil {
			defer conn.Close()
			return false, fmt.Sprintf("Read err:%s", err)
		} else {
			//fmt.Println("XXXXX", len(string(response[8:])))
			//http://blog.csdn.net/small_qch/article/details/19562661
			//http://andylin02.iteye.com/blog/648412?spm=5176.100239.blogcont58999.6.6KeZAX
			//Fastcgi前8个字节怎么解析啊遇到点问题
			//return true, rand_data(string(response[8:]))
			//fmt.Println("111111:", string(response[0:]))
			defer conn.Close()
			return true, rand_data(BytesString(response)) //rand_data(BytesString(response))
		}

	}
	defer conn.Close()
	return false, ""
}

func HTTP_Version(r *http.Request) map[string]string { //构造HTTP头
	env := make(map[string]string) //map[string]string
	//fmt.Println(reflect.TypeOf(env))
	//http://php.net/manual/zh/reserved.variables.server.php
	//CGI/1.1 标准中。只有下列变量定义在其中： AUTH_TYPE， CONTENT_LENGTH, CONTENT_TYPE, GATEWAY_INTERFACE, PATH_INFO,
	//PATH_TRANSLATED, QUERY_STRING, REMOTE_ADDR, REMOTE_HOST, REMOTE_IDENT, REMOTE_USER, REQUEST_METHOD, SCRIPT_NAME,
	//SERVER_NAME, SERVER_PORT, SERVER_PROTOCOL 和 SERVER_SOFTWARE。其它的变量均作为“供应商扩展（vendor extensions）”来对待。
	//	//env["PHP_SELF"] = "\\1.php" ///server_indices.php
	//当前执行脚本的绝对路径名。
	path_name := Public_file.Get_CurrentPath() + "www" + strings.Replace(r.URL.String(), "/", "\\", -1) //路径  函数名
	env["SCRIPT_FILENAME"] = path_name                                                                  //E:/web/server_indices.php
	//当前运行脚本所在的文档根目录。在服务器配置文件中定义。
	env["DOCUMENT_ROOT"] = Public_file.Get_CurrentPath() + "www\\" //E:/web/
	//访问页面时的请求方法。例如：“GET”、“HEAD”，“POST”，“PUT”。
	env["REQUEST_METHOD"] = r.Method //GET
	//post提交数据
	//if r.Method == "POST" {
	//env["PHP_VALUE"] = "allow_url_include = On\ndisable_functions = \nsafe_mode = Off\nauto_prepend_file = php://input"
	//}
	env["HTTP_HOST"] = ""   //    localhost
	env["SERVER_ADDR"] = "" //127.0.0.1:9004
	env["SERVER_PORT"] = "" //80
	//当前运行脚本所在服务器主机的名称。
	env["SERVER_NAME"] = "" //localhost
	//服务器使用的 CGI 规范的版本。例如，“CGI/1.1”。
	env["GATEWAY_INTERFACE"] = "CGI/1.1" //CGI/1.1
	//服务器标识的字串，在响应请求时的头部中给出。
	env["SERVER_SOFTWARE"] = "C++ / fcgiclient" //Apache/2.2.22 (Win64) PHP/5.3.13
	//请求页面时通信协议的名称和版本。例如，“HTTP/1.0”。
	env["SERVER_PROTOCOL"] = r.Proto //HTTP/1.1

	//传不进去！！！！  string(r.Header["User-Agent"][0])
	//env["HTTP_USER_AGENT"] = "" //    Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.17 (KHTML, like Gecko)
	//查询(query)的字符串。
	env["QUERY_STRING"] = r.URL.RawQuery
	env["DOCUMENT_URI"] = ""
	env["HTTPS"] = "" //    -

	//正在浏览当前页面用户的 IP 地址。
	env["REMOTE_ADDR"] = "" //127.0.0.1  127.0.0.1:8070  不知道为何传进去就出错了
	//用户连接到服务器时所使用的端口。
	env["REMOTE_PORT"] = "" //65037

	//访问此页面所需的 URI。例如，“/index.html”。
	env["REQUEST_URI"] = "" //  /server_indices.php   r.URL.Path + "?" + r.URL.RawQuery
	env["SCRIPT_NAME"] = "" //  /server_indices.php
	//env["PATH_INFO"] = r.URL.Path
	env["CONTENT_LENGTH"] = "" //文件大小
	env["CONTENT_TYPE"] = ""
	env["REQUEST_SCHEME"] = ""

	//	//env["argv"] = ""                    //    -
	//	//env["argc"] = ""                    //    -
	//	env["REQUEST_TIME"] = "1361542579"          //    1361542579
	//	//env["REQUEST_TIME_FLOAT"] = ""
	//	//env["HTTP_ACCEPT"] = ""                                                   //    text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
	//	//env["HTTP_ACCEPT_CHARSET"] = ""                                                          //    ISO-8859-1,utf-8;q=0.7,*;q=0.3
	//	//env["HTTP_ACCEPT_ENCODING"] = ""                                                         //    gzip,deflate,sdch
	//	//env["HTTP_ACCEPT_LANGUAGE"] = ""                                                         //    fr-FR,fr;q=0.8,en-US;q=0.6,en;q=0.4
	//	//env["HTTP_CONNECTION"] = ""                                                              //    keep-alive
	//
	//	env["HTTP_REFERER"] = ""                                                                 //    http://localhost/
	//	env["REMOTE_HOST"] = ""                                                                  //    -
	//	env["REMOTE_USER"] = ""                                                                  //    -
	//	env["REDIRECT_REMOTE_USER"] = ""                                                         //    -
	//env["SERVER_ADMIN"] = ""                                                                 //    myemail@personal.us
	//	env["SERVER_SIGNATURE"] = ""                                                             //
	//	env["PATH_TRANSLATED"] = ""                                                              //    -
	//
	//	env["PHP_AUTH_DIGEST"] = ""                                                              //    -
	//	env["PHP_AUTH_USER"] = ""                                                                //    -
	//	env["PHP_AUTH_PW"] = ""                                                                  //    -
	//	env["AUTH_TYPE"] = ""                                                                    //    -
	//	                                                                   //    -
	//	env["ORIG_PATH_INFO"] = ""                                                               //    -

	return env
}

#golangPHPcgi   GOphp--GO解析PHP源码并实现一个miniPHP服务起器<br>
#by<br>
golang php cgi   github:https://github.com/webxscan/gophp<br>
BLOG： http://blog.csdn.net/webxscan/ <br>
BY：斗转星移 QQ:29295842 <br>
####软件目的<br>
实现一个本地PHP解析器，不用使用阿帕奇或者IIS。<br>
这样就可以实现很多自定义扩展。<br>
软件目前写了4天，还有很多不完美的地方还希望大家予以纠正。<br>
<br>
<br>
###代码<br>
```
package main

import (
	"fmt"
	"net"
	"strconv"
	"Public_file" //公用文件
	"fastcgi"
	"time"
	"www"
)

func main() {
	fmt.Println("PHP CGI WEB Server Example text V:1.0")
	fmt.Println("BY:29295842@qq.com")
	//	fmt.Printf(Cmdexec(".\\php\\php546x161220011555\\php.exe", "windows"))
	//	//fmt.Printf(Cmdexec(".\\php\\php546x161220011555\\php.exe .\\php\\php546x161220011555\\1.php", "windows"))
	ip := "127.0.0.1"
	//port := "9002"
	port := 9000
	url := ""
	for {
		url = fmt.Sprintf("%s:%d", ip, port)          //127.0.0.1
		tcpAddr, _ := net.ResolveTCPAddr("tcp4", url) //转换IP格式
		_, err := net.DialTCP("tcp", nil, tcpAddr)    //查看是否连接成功
		if err != nil {
			break //#跳出
		}
		port++
		time.Sleep(1 * time.Second)
	}
	Public_file.Cgi_ip = ip                    //记录
	Public_file.Cgi_port = strconv.Itoa(port)  //记录
	go fastcgi.Run_cgi(ip, strconv.Itoa(port)) //PHP cgi

	http_url := "127.0.0.1:8070"
	go www.Www_root(http_url) //启动WEB网站

	time.Sleep(1 * time.Second)

	for { //死循环
		time.Sleep(10 * time.Second)
	}

	//make一个chan用于阻塞主线程,避免程序退出
	//	blockMainRoutine := make(chan bool)
	//	<-blockMainRoutine
}
```
<br>
<br>
```
CGI
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
```
<br>
<br>
#测试<br>
<img src="http://img.blog.csdn.net/20170412105813245?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvd2VieHNjYW4=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast"  alt="pyqteval" />
<img src="http://img.blog.csdn.net/20170412105704354?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvd2VieHNjYW4=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast"  alt="pyqteval" />

<img src="http://img.blog.csdn.net/20170412105827855?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvd2VieHNjYW4=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast"  alt="pyqteval" />
<img src="http://img.blog.csdn.net/20170412105849072?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvd2VieHNjYW4=/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast"  alt="pyqteval" />

```
D:/go32/bin/go.exe build -i [C:/Users/Administrator/Desktop/gotophp/GOPHP]
成功: 进程退出代码 0.
C:/Users/Administrator/Desktop/gotophp/GOPHP/GOPHP.exe  [C:/Users/Administrator/Desktop/gotophp/GOPHP]
PHP CGI WEB Server Example text V:1.0
BY:29295842@qq.com
http WEB Server run  http://127.0.0.1:8070
run cgi
path_name: C:\Users\Administrator\Desktop\gotophp\GOPHP\www\1.php
MIME: text/html
xxxxxxxx:C:\Users\Administrator\Desktop\gotophp\GOPHP\www\1.php 
path_name: C:\Users\Administrator\Desktop\gotophp\GOPHP\www\favicon.ico
path_name: C:\Users\Administrator\Desktop\gotophp\GOPHP\www\1.jpg
MIME: image/jpeg
path_name: C:\Users\Administrator\Desktop\gotophp\GOPHP\www\favicon.ico

```


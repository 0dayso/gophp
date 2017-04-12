package main

import (
	"Public_file" //公用文件
	"fastcgi"
	"fmt"
	"net"
	"strconv"
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

package www

import (
	//"LogX"
	"Public_file"
	"fastcgi"
	"fmt"
	"io"
	"net/http"
	//	"reflect"
	"bufio"
	"os"
	"strings"
)

var (
	file_MIME = make(map[string]string) //保存网页格式解析
)

func Www_root(url string) {
	open_file_MIME() //文件格式解析
	//建立HTTP  UI控制
	//url := "127.0.0.1:8070"
	data := fmt.Sprintf("http WEB Server run  http://%s", url)
	fmt.Println(data)
	//http.HandleFunc("www.a.com/", handlerFunc)
	http.HandleFunc("/", Url_path)
	err := http.ListenAndServe(url, nil)
	if err != nil {
		fmt.Printf("log_error--%s\n", err)
	}

}

func open_file_MIME() { //文件格式解析
	f, err := os.Open("MIME.txt")
	if err != nil {
		//return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		//fmt.Println(line)
		s := strings.Split(line, "|")
		if len(s) >= 2 {
			file_MIME[s[0]] = s[1]
		}

		//		handler(line)
		if err != nil {
			if err == io.EOF {
				break //#跳出
				//return
			}
			break //#跳出
			//return
		}
	}
}

func find_file_MIME(url_path string) string { //查找文件解析方式
	path_name_x := strings.ToLower(url_path) //转换为小写
	for k, v := range file_MIME {
		if strings.Contains(path_name_x, k) {
			//fmt.Printf("%s---%s \n", k, v)
			return v
		}
	}
	return "application/octet-stream"
}

func Url_path(web http.ResponseWriter, r *http.Request) { //路径解析
	path_name := Public_file.Get_CurrentPath() + "www" + strings.Replace(r.URL.String(), "/", "\\", -1) //路径  函数名
	//path_name := r.URL.String() //路径  函数名
	fmt.Println("path_name:", path_name)
	//fmt.Println("rrrrrr:", r.URL.String())
	BOOL, err := Public_file.PathExists(path_name)
	if BOOL == false { //文件不存在
		err_404_file := fmt.Sprintf("404  %s", err)
		io.WriteString(web, err_404_file)
		return
	}
	cookie := http.Cookie{Name: "UserName", Value: "12345678", Path: "/", MaxAge: 86400}
	http.SetCookie(web, &cookie)
	//文件返回格式头
	MIME := find_file_MIME(r.URL.String()) //查找文件解析方式
	fmt.Println("MIME:", MIME)
	web.Header().Set("Content-Type", MIME)

	path_name_x := strings.ToLower(path_name)
	//PHP解析
	if strings.Contains(path_name_x, ".php") {
		fmt.Printf("xxxxxxxx:%s \n", path_name)
		//fmt.Printf(r)
		BOOL, rand_data := fastcgi.Request(Public_file.Cgi_ip, Public_file.Cgi_port, fastcgi.HTTP_Version(r))
		if BOOL {
			//fmt.Println(rand_data)
			io.WriteString(web, string(rand_data))
			return
		} else {
			http_rand := fmt.Sprintf("php except---%s", rand_data)
			fmt.Println(http_rand)
			io.WriteString(web, string(http_rand))
			return
		}
	}
	//直接返回源文件
	rand_data := Public_file.ReadAll(path_name)
	io.WriteString(web, rand_data)
	return

	//直接返回源文件
	//	arr := [...]string{".txt", ".htm", ".html"}
	//	for index, value := range arr {
	//		if strings.Contains(path_name_x, value) {
	//			fmt.Printf("arr[%d]=%s \n", index, value)
	//			rand_data := Public_file.ReadAll(path_name)
	//			io.WriteString(web, rand_data)
	//			return
	//		}
	//	}
}

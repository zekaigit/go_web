package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	//fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}
func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                    //解析参数，默认是不会解析的
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("login.html")
		t.Execute(w, token)
		//log.Println(t.Execute(w, nil))
	} else {
		//请求的是登录数据，那么执行登录的逻辑判断

		token := r.Form.Get("token")
		if token != "" {
			//验证token合法性
		} else {
			//不存在token报错
		}
		fmt.Println("username length:", len(r.Form["username"][0]))
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
		/*
			fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username"))) //输出到服务器端
			fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
			template.HTMLEscape(w, []byte(r.Form.Get("username"))) //输出到客户端
		*/
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(curtime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

}

func main() {
	//http.HandleFunc("/", sayhelloName)       //设置访问的路由
	http.HandleFunc("/login", login)         //设置访问的路由
	http.HandleFunc("/upload", upload)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

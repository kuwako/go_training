package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	//オプションを解析
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!")
}

func main() {
	//アクセスのルーティングを設定
	http.HandleFunc("/", sayhelloName)
	//監視するポートを設定
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"./configs"
	"./controllers"
	"./models"
)

func init() {
	models.DbInit()
}

func main() {
	addr := configs.HOST + ":" + configs.PORT

	router := controllers.NewRouter()
	http.Handle("/", router)

	// 由于使用的是相对路径，程序需要运行在本文件所在目录下
	fsh := http.FileServer(http.Dir("../public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fsh))

	go func() {
		fmt.Println("pprof")
		http.ListenAndServe("localhost:12345", nil)
	}()

	fmt.Printf("start server at %v...\n", addr)
	http.ListenAndServe(addr, nil)
}

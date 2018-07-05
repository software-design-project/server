package main

import (
	"fmt"
	"net/http"

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

	fmt.Printf("start server at %v...\n", addr)
	http.ListenAndServe(addr, nil)
}

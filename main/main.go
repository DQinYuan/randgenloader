package main

import (
	"fmt"
	"github.com/DQinYuan/randgenloader"
	"net/http"
)

/*
export CONFPATH=/home/dqyuan/language/Go/projects/randgenloader
export RMPATH=/home/dqyuan/language/Mysql/randgenx
export RESULTPATH=/home/dqyuan/language/Go/projects/randgenloader
 */

func main() {
	startServer()
}

func startServer() {
	http.HandleFunc("/init", randgenloader.Init)
	http.HandleFunc("/loaddata", randgenloader.LoadData)
	http.HandleFunc("/query", randgenloader.Query)
	http.HandleFunc("/Compare", randgenloader.Compare)
	err := http.ListenAndServe(":9080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

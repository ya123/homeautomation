package main

import (
	"fmt"
	. "github.com/yannis/home_automation"
	"net/http"
)

var (
	device1 = NewDevice("Device1", NewLedPlugger("red", 0.5))
	device2 = NewDevice("Device2", NewLedPlugger("green", 0.7))
)

func main() {
	http.HandleFunc("/", List)
	http.HandleFunc("/details", Details)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/change", Change)

	fmt.Println("started")

	err := http.ListenAndServe("127.0.0.1:80", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

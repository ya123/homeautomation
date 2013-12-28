package main

import (
	"fmt"
	. "github.com/ya123/homeautomation"
	"net/http"
	"os"
	"path"
	"strconv"

//	"time"
)

var static = path.Join(os.Getenv("GOPATH"), "src", "github.com", "ya123", "homeautomation", "static", "")

func dim(rw http.ResponseWriter, req *http.Request) {
	val := req.URL.Query().Get("value")

	//time.Sleep(time.Second * 2)
	rw.Header().Set("Content-Type", "application/json")
	i, err := strconv.Atoi(val)
	if err != nil {
		rw.Write([]byte(`{"result": 0}`))
		return
	}
	rw.Write([]byte(fmt.Sprintf(`{"result": %d}`, i+10)))
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Fehler: Bitte Server-Adresse angeben, z.B.\n\tlocalhost:80 oder \n\t192.168.220.134:80\n")
		os.Exit(1)
	}
	address := os.Args[1]

	Load()

	http.HandleFunc("/details", Details)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/change", Change)
	http.HandleFunc("/on", On)
	http.HandleFunc("/off", Off)
	http.HandleFunc("/dim", dim)
	fileserver := http.FileServer(
		//http.Dir(`D:\GO\GOPATH\src\github.com\yannis\home_automation\static\`),
		http.Dir(static),
	)
	http.Handle("/css/", fileserver)
	http.Handle("/bootstrap-3.0.0/", fileserver)
	http.Handle("/jquery-ui-1.10.3/", fileserver)
	http.Handle("/js/", fileserver)
	http.HandleFunc("/", List)

	fmt.Printf("started on %s\n", address)

	err := http.ListenAndServe(address, nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

package main

import (
	"encoding/json"
	"frp-api/frpc"
	"log"
	"net/http"
	"os"
)

//var FRP_PATH = os.Getenv("FRP_DIR")
var FRP_WORK_DIRECTORY = getEnv("FRP", "./")
var PORT = getEnv("PORT", "5600")

var client *frpc.Frpc


func apiFrpcAddHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.RequestURI)
	switch r.Method {

	case "GET":
		w.Write([]byte("hi"))

	case "POST":
		d := json.NewDecoder(r.Body)
		p := &frpc.Proxy{}
		if err := d.Decode(p); err != nil {
			log.Println("failed to decode json")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Println("Add proxy:", p.String())

	default:
		http.Error(w, "unsupported method: " + r.Method, http.StatusBadRequest)

	}
}

func main() {
	client = frpc.NewFrpc(FRP_WORK_DIRECTORY)
	http.HandleFunc("/api/frpc", apiFrpcAddHandler)
	err := http.ListenAndServe(":" + PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
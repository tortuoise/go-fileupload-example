package main

import (
	_ "fmt"
	"io"
	"log"
	"net/http"
)

func handleUpload(w http.ResponseWriter, r *http.Request) {

	log.Println(r.Header)
	log.Println(r.URL)
	rdr, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
	}
	for {
		np, err := rdr.NextPart()
		if err == io.EOF {
			break
		}
		bs := make([]byte, 10000)
		n, err := np.Read(bs)
		log.Println(n, " bytes ", string(bs))
	}
}

func main() {

	http.HandleFunc("/", handleUpload)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

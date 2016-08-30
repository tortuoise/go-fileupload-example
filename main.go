package main

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Streams upload directly from file -> mime/multipart -> pipe -> http-request
func streamingUploadFile(params map[string]string, paramName, path string, w *io.PipeWriter, file *os.File) {
	defer file.Close()
	defer w.Close()
	writer := multipart.NewWriter(w)
	_ = writer.SetBoundary("doofus")
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()
	go streamingUploadFile(params, paramName, path, w, file)
	return http.NewRequest("POST", uri, r)
}

func main() {
	path, _ := os.Getwd()
	path += "/test.pdf"
	extraParams := map[string]string{
		"title":       "My Document",
		"author":      "Matt Aimonetti",
		"description": "A document with all the Go programming language secrets",
	}
	request, err := newfileUploadRequest("http://localhost:8080", extraParams, "file", "products_no_header.csv")
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "multipart/form-data; boundary=doofus")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		_, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
	}
}

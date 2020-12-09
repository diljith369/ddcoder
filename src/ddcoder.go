package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var ddcodertemplate *template.Template
var result string

func init() {
	ddcodertemplate = template.Must(template.ParseFiles("template/ddcoder.html"))

}
func main() {
	startserver()
}

func ddocoderhomepage(httpw http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		err := ddcodertemplate.Execute(httpw, nil)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := req.ParseForm()
		if err != nil {
			fmt.Println(err)
		}
		userval := req.Form.Get("userstring")
		fmt.Println(userval)

		choice := req.Form.Get("userselection")
		fmt.Println(choice)
		if choice == "Base64 Encode" {
			result = operations(userval, "b64")
		} else if choice == "Base64 Decode" {
			result = operations(userval, "decb64")
		} else if choice == "Compress and Encode (GZip/Base64)" {
			result = operations(userval, "cmpenc")
		} else if choice == "Decompress and Decode (GZip/Base64)" {
			result = operations(userval, "decompdec")
		}

		fmt.Println(result)
		err = ddcodertemplate.Execute(httpw, result)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func operations(userval, optype string) string {
	var result string
	if optype == "b64" {
		result = base64.StdEncoding.EncodeToString([]byte(userval))
	} else if optype == "decb64" {
		dcodedval, _ := base64.StdEncoding.DecodeString(userval)
		result = string(dcodedval)
	} else if optype == "cmpenc" {
		result = compressandcode(userval)
	} else if optype == "decompdec" {
		result = decodeanddecompress(userval)
	}
	return result
}

func decodeanddecompress(userval string) string {
	dcodedval, _ := base64.StdEncoding.DecodeString(userval)
	zipbuffer := bytes.NewReader(dcodedval)
	decompressedbytes, _ := gzip.NewReader(zipbuffer)
	originalval, _ := ioutil.ReadAll(decompressedbytes)
	return string(originalval)
}

func compressandcode(userval string) string {
	var zipbuffer bytes.Buffer
	gzwriter := gzip.NewWriter(&zipbuffer)
	if _, err := gzwriter.Write([]byte(userval)); err != nil {
		log.Fatal(err)
	}
	if err := gzwriter.Flush(); err != nil {
		log.Fatal(err)
	}
	if err := gzwriter.Close(); err != nil {
		log.Fatal(err)
	}
	result := base64.StdEncoding.EncodeToString(zipbuffer.Bytes())
	return result
}
func startserver() {
	router := mux.NewRouter()
	router.HandleFunc("/", ddocoderhomepage)

	router.PathPrefix("/static/css/").Handler(http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css/"))))

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8085",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 180 * time.Second,
		ReadTimeout:  180 * time.Second,
	}
	srv.ListenAndServe()
}

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
	/assets/css/custom.js
	/assets/img/logo.png
	/assets/js/custom.js
	/assets/js/jquery-3.4.1.min.js
	/assets/js/jquery.mask.min.js
	/assets/js/js.cookie-2.2.1.min.js
	/assets/js/popper.min.js
	/assets/plugins/bootstrap-table/bootstrap-table.min.css
	/assets/plugins/bootstrap-table/bootstrap-table.min.js
	/assets/plugins/bootstrap/bootstrap.min.css
	/assets/plugins/bootstrap/bootstrap.min.js
	/assets/plugins/moment/moment-timezone-with-data.min.js
	/assets/plugins/moment/moment-timezone.min.js
	/assets/plugins/moment/moment.min.js
*/
var assetMap map[string][]byte

func main() {

	compressed, err := decode(encoded)
	if err != nil {
		panic(err)
	}

	b, err := decompress(compressed)
	if err != nil {
		panic(err)
	}

	var m map[string][]byte
	err = json.Unmarshal(b, &m)

	//spew.Dump(m)

	//http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	//
	//m, err := decode(encoded)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//
	for key, src := range m {
		fmt.Println(key)
		http.HandleFunc(key, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, string(src))
		})

	}
	log.Fatal(http.ListenAndServe(":8080", nil))

	//spew.Dump(m)
	//fmt.Println(m["assets/css/bootstrap.min.css"])
	//    // fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	//})
	//
	//log.Fatal(http.ListenAndServe(":8080", nil))
}

func decompress(b []byte) ([]byte, error) {
	buf := bytes.NewBuffer(b)

	var r io.Reader
	var err error
	r, err = gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return resB.Bytes(), nil

}

func decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
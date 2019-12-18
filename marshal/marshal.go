package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

var dataDir = "data"

func main() {

	list := []string{
		"assets/css/bootstrap.min.css",
		"assets/js/bootstrap.min.js",
		"assets/js/jquery-3.4.1.min.js",
		"assets/js/popper.min.js",
	}

	m, err := encode(list)
	if err != nil {
		fmt.Println(err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}

	str := base64.StdEncoding.EncodeToString(b)
	fmt.Printf(str)

}

func encode(files []string) (map[string][]byte, error) {
	m := make(map[string][]byte)
	for _, path := range files {
		p := filepath.Join(dataDir, path)
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		m[path] = b
	}
	return m, nil
}

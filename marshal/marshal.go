package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var resourceDir = "data"

var extMap = map[string]bool{
	".js":  true,
	".png": true,
	".css": true,
}

func main() {

	list := getAssets()
	if len(list) < 1 {
		return
	}
	//spew.Dump(list)

	//fmt.Println("")
	//fmt.Println("")
	//fmt.Println("")

	encodedStr, m := encodeAssets(list)
	f, err := ioutil.TempFile(".", "result")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("/*\n")
	//fmt.Println("Keys---------------------------------")

	keys := make([]string, 0)
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		f.WriteString("\t" + k + "\n")
	}
	f.WriteString("*/\n")
	f.WriteString(encodedStr)

	fmt.Printf("output: %s\n", f.Name())
}

func encodeAssets(list []string) (string, map[string][]byte) {
	m, err := encode(list)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(b), m

}

func getAssets() []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(resourceDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if _, exists := extMap[filepath.Ext(path)]; !exists {
			return nil
		}

		list = append(list, filepath.ToSlash(path))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return list

}

func encode(files []string) (map[string][]byte, error) {
	m := make(map[string][]byte)
	for _, path := range files {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		m[strings.TrimPrefix(path, resourceDir)] = b
		fmt.Printf("[%s] is encoded successfully. len=%d\n", filepath.Base(path), len(b))
	}
	return m, nil
}

package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Hashing(str string) string {
	hash := md5.Sum([]byte(str))
	hashString := hex.EncodeToString(hash[:])
	return hashString
}

func LoadPage(w http.ResponseWriter, str string, data any) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("/Users/artemlukmanov/GolandProjects/warehouse-application/pkg/pages/%s.html", str))
	if err != nil {
		log.Println(err)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

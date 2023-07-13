package helper

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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

func Unmarshal(r *http.Request, structura any) *any {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	if err := json.Unmarshal(bytes, &structura); err != nil {
		log.Println(err)
		return nil
	}

	return &structura
}

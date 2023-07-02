package page

import (
	"fmt"
	"net/http"
	"os"
	"warehouse-application/pkg/errors"
)

func loadPage(title string) []byte {
	fileName := fmt.Sprintf("/Users/kare/GolandProjects/warehouse-application/page/pages/%s.html", title)
	body, err := os.ReadFile(fileName)
	if err != nil {
		return []byte("Error")
	}
	return body
}

func ShowPage(w http.ResponseWriter, r *http.Request, title string) {
	page := loadPage(title)
	_, err := w.Write(page)
	errors.CheckWarning(err)
}

package server

import (
	"fmt"
	"net/http"
)

type input struct {
	username string `json:"username"`
	password string `json:"password"`
}

func StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<!DOCTYPE html>\n<html>\n  <head>\n    <meta charset=\"UTF-8\">\n    <title>Авторизация</title>\n  </head>\n  <body>\n    <h1>Авторизация</h1>\n    <form id=\"login-form\">\n      <label for=\"username\">Логин:</label>\n      <input type=\"text\" id=\"username\" name=\"username\"><br>\n      <label for=\"password\">Пароль:</label>\n      <input type=\"password\" id=\"password\" name=\"password\"><br>\n      <button type=\"submit\">Войти</button>\n    </form>\n    <br>\n    <form id=\"register-form\">\n      <label for=\"username\">Логин:</label>\n      <input type=\"text\" id=\"username\" name=\"username\"><br>\n      <label for=\"password\">Пароль:</label>\n      <input type=\"password\" id=\"password\" name=\"password\"><br>\n      <button type=\"submit\">Зарегистрироваться</button>\n    </form>\n    <script>\n      const registerForm = document.querySelector('#register-form');\n      registerForm.addEventListener('submit', (event) => {\n        event.preventDefault();\n        const formData = new FormData(registerForm);\n        const xhr = new XMLHttpRequest();\n        xhr.open('POST', 'http://localhost:8080/users');\n        xhr.setRequestHeader('Content-Type', 'application/json');\n        xhr.onload = () => {\n          if (xhr.status === 200) {\n            console.log(xhr.responseText);\n            // обработка успешного ответа от сервера\n          } else {\n            console.error(xhr.statusText);\n            // обработка ошибки\n          }\n        };\n        const data = {\n          username: formData.get('username'),\n          password: formData.get('password')\n        };\n        xhr.send(JSON.stringify(data));\n      });\n    </script>\n  </body>\n</html>")
	})
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("post")
		}
		err := r.ParseForm()
		if err != nil {
			fmt.Println("error")
			w.WriteHeader(http.StatusBadRequest)
		}
		fmt.Println(r.Form.Get("username"))
		fmt.Println(r.Form.Get("password"))
	})
	http.ListenAndServe("localhost:8080", nil)
}

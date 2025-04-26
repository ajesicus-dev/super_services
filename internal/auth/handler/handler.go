package handler

import "net/http"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Login endpoint"))
}
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// success - call the next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth to provide auth layer for a particular handler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	fmt.Println("length: ", len(segs))
	if len(segs) != 4 {
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

// google client ID
// 516301198069-io33qbi2edp7360rf7h55tgbp2sa379q.apps.googleusercontent.com
// google client secret
// TBKKQ6TZeEN4byopoICIxkdS

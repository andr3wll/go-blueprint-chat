package main

import (
	"blueprint/chat_01/trace"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	fbID, exists := os.LookupEnv("FB_ID")
	if !exists {
		log.Print("FB_ID not found")
	}
	fbSecret, exixts := os.LookupEnv("FB_SECRET")
	if !exixts {
		log.Print("FB_SECRET not found")
	}
	githubID, exists := os.LookupEnv("GITHUB_ID")
	if !exists {
		log.Print("GITHUB_ID not found")
	}
	githubSecret, exists := os.LookupEnv("GITHUB_SECRET")
	if !exists {
		log.Print("GITHUB_SECRET not found")
	}
	googleID, exists := os.LookupEnv("GOOGLE_ID")
	if !exists {
		log.Print("GOOGLE_ID not found")
	}
	googleSecret, exists := os.LookupEnv("GOOGLE_SECRET")
	if !exists {
		log.Print("GOOGLE_SECRET not found")
	}

	// setup gomniauth
	gomniauth.SetSecurityKey("#^dDUR5/NSdSs/k)")
	gomniauth.WithProviders(
		facebook.New(fbID, fbSecret,
			"http://localhost:8080/auth/callback/facebook"),
		github.New(githubID, githubSecret,
			"http://localhost:8080/auth/callback/github"),
		google.New(googleID, googleSecret,
			"http://localhost:8080/auth/callback/google"),
	)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	})
	// room
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/room", r)
	// get the room going
	go r.run()
	// start the web server
	log.Println("Starting web service on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

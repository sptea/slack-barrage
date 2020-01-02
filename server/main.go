package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/sptea/slack-barrage/server/entity"
)

var listenPort string
var logger *log.Logger
var store *sessions.CookieStore
var session *sessions.Session
var upgrader = websocket.Upgrader{}

const (
	sessionName = "sid"
)

func clientHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
		// TODO Fatal is not suitable to here (should return error)
	}
	defer conn.Close()

	entity.ClientHandler(conn, logger)
}

func authFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, sessionName)
		logger.Println("kokoInauthFilter")
		logger.Println(session.Values["userInfo"])
		logger.Println(session.Values)

		if session.Values["userInfo"] == nil {
			logger.Println("redirect")
			http.Redirect(w, r, "/auth", 302)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	url := AuthMethod(session)

	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, 302)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, sessionName)

	err := callbackMethod(session, r.FormValue("state"), r.FormValue("code"), logger)
	if err != nil {
		logger.Println(err)
	}

	session.Save(r, w)
}

func init() {
	logger = log.New(os.Stdout, "[slack-barrage]", log.LstdFlags)

	logger.Println("init")
	err := godotenv.Load()
	if err != nil {
		logger.Panic("Error loading .env file")
	}

	listenPort = os.Getenv("PORT")
	if listenPort == "" {
		logger.Panic("Faild to load env: PORT")
	}

	InitAuthConfig(logger)
}

func main() {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		logger.Panic("Faild to load env: SESSION_KEY")
	}
	store = sessions.NewCookieStore([]byte(sessionKey))
	session = sessions.NewSession(store, sessionName)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", testHandler)
	r.Get("/auth", authHandler)
	r.Get("/callback", callbackHandler)

	r.Get("/ws", clientHandler)

	r.Route("/admin", func(r chi.Router) {
		r.Use(authFilter)
		r.Get("/", testHandler)
	})

	entity.MessageRoutine(logger)

	logger.Printf("Started to listen: " + listenPort)
	http.ListenAndServe(":"+listenPort, r)
}

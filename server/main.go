package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/sptea/slack-barrage/server/domain"
	v2 "google.golang.org/api/oauth2/v2"
)

var listenPort string
var logger *log.Logger
var store *sessions.CookieStore
var session *sessions.Session

func clientHandler(w http.ResponseWriter, r *http.Request) {
	domain.ClientHandler(w, r, logger)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	state, _ := uuid.NewV4()

	config := domain.GetConfig()

	url := config.AuthCodeURL(state.String())
	http.Redirect(w, r, url, 302)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	config := domain.GetConfig()

	context := context.Background()
	code := r.FormValue("code")

	tok, err := config.Exchange(context, code)
	if err != nil {
		panic(err)
	}

	if tok.Valid() == false {
		panic(errors.New("vaild token"))
	}

	service, _ := v2.New(config.Client(context, tok))
	tokenInfo, _ := service.Tokeninfo().AccessToken(tok.AccessToken).Context(context).Do()

	session, _ := store.Get(r, "session-name")
	session.Values["email"] = tokenInfo.Email
	fmt.Println(tokenInfo.Email)
	session.Save(r, w)

	json.NewEncoder(w).Encode(tokenInfo)
}

func init() {
	logger = log.New(os.Stdout, "[slack-barrage]", log.LstdFlags)

	err := godotenv.Load()
	if err != nil {
		logger.Panic("Error loading .env file")
	}

	listenPort = os.Getenv("PORT")
	if listenPort == "" {
		logger.Panic("Faild to load env: PORT")
	}

	domain.InitAuthConfig(logger)
}

func main() {
	b := make([]byte, 48)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}
	str := strings.TrimRight(base32.StdEncoding.EncodeToString(b), "=")
	store = sessions.NewCookieStore([]byte(os.Getenv(str)))
	session = sessions.NewSession(store, "session-name")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/auth", authHandler)
	r.Get("/callback", callbackHandler)

	r.Get("/ws", clientHandler)

	go domain.ReadMessageFromslack(logger)
	go domain.BroadcastMessagesToClients(logger)

	logger.Printf("Started to listen: " + listenPort)
	http.ListenAndServe(":"+listenPort, r)
}

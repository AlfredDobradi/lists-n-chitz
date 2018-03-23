package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alfreddobradi/lists-n-chitz/internal/auth"
	"github.com/alfreddobradi/lists-n-chitz/internal/link"
	"github.com/alfreddobradi/lists-n-chitz/internal/logger"
	"github.com/zenazn/goji"
)

var log = logger.New()

func main() {
	goji.Get("/", root)

	goji.Post("/signup", signup)
	goji.Post("/login", login)
	goji.Post("/save", save)

	goji.Serve()
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Routes:\n------\nRegister:\tPOST /signup\nLogin:\t\tPOST /login\nSave:\t\tGET /save")
}

func signup(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var data auth.User

	err := decoder.Decode(&data)
	if err != nil {
		log.Warningf("server: json: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 - Internal Server Error")

		return
	}

	u, err := auth.Register(data)
	if err != nil {
		log.Warningf("server: register: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "500 - Internal Server Error")

		return
	}

	userResponse := auth.UserResponse{
		ID:        u.ID,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		Email:     u.Email,
	}

	response, _ := json.Marshal(userResponse)

	fmt.Fprintf(w, "%s", response)
}

func login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data auth.User
	err := decoder.Decode(&data)
	if err != nil {
		log.Warningf("server: json: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 - Internal Server Error")

		return
	}

	address := getAddress(r)

	t, err := auth.Authenticate(data, address)
	if err != nil {
		log.Warningf("server: authenticate: %+v", err)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "401 - Unauthorized")

		return
	}

	response, err := json.Marshal(t)
	if err != nil {
		log.Warningf("server: json: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 - Internal Server Error")

		return
	}

	fmt.Fprintf(w, "%s", response)
}

func save(w http.ResponseWriter, r *http.Request) {
	address := getAddress(r)
	token := r.Header.Get("Authorization")

	user, err := auth.Authorize(address, token)
	if err != nil {
		log.Warningf("server: authorize: %+v", err)

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "401 - Unathorized")

		return
	}

	decoder := json.NewDecoder(r.Body)
	var data link.Link
	err = decoder.Decode(&data)
	if err != nil {
		log.Warningf("server: json: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 - Internal Server Error")

		return
	}

	data.UserID = user.ID

	data, err = link.Save(data)
	if err != nil {
		log.Warningf("server: json: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 - Internal Server Error")

		return
	}

	response, err := json.Marshal(data)
	if err != nil {
		log.Warningf("server: json: %+v", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "500 - Internal Server Error")

		return
	}
	log.Debugf("%+v", response)
	fmt.Fprintf(w, "%s", response)
}

func getAddress(r *http.Request) (address string) {
	index := strings.LastIndex(r.RemoteAddr, ":")
	if index != -1 {
		address = r.RemoteAddr[1 : index-1]
	} else {
		address = "127.0.0.1"
	}

	return
}

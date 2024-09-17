package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type myHandler string

func (mh myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Powered-By", "Matcha")

	http.SetCookie(w, &http.Cookie{
		Name:    "session-id",
		Value:   "12345",
		Expires: time.Now().Add(24 * time.Hour * 365),
	})

	w.WriteHeader(http.StatusAccepted)

	_, err := fmt.Fprintln(w, string(mh))
	if err != nil {
		http.Error(w, "Internal server error: ", http.StatusInternalServerError)
		log.Println("Err -> ", err)
	}

	_, err = fmt.Fprintln(w, r.Header)
	if err != nil {
		log.Println("Error on header -> ", err)
	}
}

func main() {

	http.Handle("/", myHandler("Customer service"))

	http.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Customer service")
		if err != nil {
			log.Println("Error -> ", err)
		}
	})

	http.HandleFunc("/url/", customHandlerFunc())

	// This is the way to create a default server
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// This is the way to create a default server with TLS (Transport Layer Security)
	// log.Fatal(http.ListenAndServeTLS(":8080", "./cert.pem", "./key.pem", nil))

	customServer(":3030")
}

func customHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
	var handlerFunc http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, r.URL.String())
		if err != nil {
			log.Fatal("Error -> ", err)
		}
	}
	return handlerFunc
}

func customServer(Addr string) {
	// This is the way to create our own custom server
	// By creating our servers, we can have as many as we want
	// But if we use a default server such as ListenAndServer, we can only have 1.
	s := http.Server{
		Addr: Addr,
	}

	go func() {
		log.Fatal(s.ListenAndServeTLS("./cert.pem", "./key.pem"))
	}()

	fmt.Println("Server has started, press <Enter> to shutdown.")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}

	err = s.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Error -> ", err)
	}
}

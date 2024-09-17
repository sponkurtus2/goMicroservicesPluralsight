package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

// This is a custom type of an http Handler
type myHandler string

// This is a function for myHandler custom type
func (mh myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We add a header by providing a key and it's value
	w.Header().Add("X-Powered-By", "Matcha")

	// And we also add some Cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session-id",
		Value:   "12345",
		Expires: time.Now().Add(24 * time.Hour * 365),
	})

	w.WriteHeader(http.StatusAccepted)

	// With this we print the current path of the url
	_, err := fmt.Fprintln(w, string(mh))
	if err != nil {
		http.Error(w, "Internal server error: ", http.StatusInternalServerError)
		log.Println("Err -> ", err)
	}

	// Print the current headers
	_, err = fmt.Fprintln(w, r.Header)
	if err != nil {
		log.Println("Error on header -> ", err)
	}
}

func main() {

	// We handle a function with a handler (myHandler)
	http.Handle("/", myHandler("Customer service"))

	// We handle an anonymous func
	http.HandleFunc("/service", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintln(w, "Customer service")
		if err != nil {
			log.Println("Error -> ", err)
		}
	})

	// We handle a normal function
	http.HandleFunc("/url/", customHandlerFunc())

	// This is the way to create a default server
	// log.Fatal(http.ListenAndServe(":8080", nil))

	// This is the way to create a default server with TLS (Transport Layer Security)
	// log.Fatal(http.ListenAndServeTLS(":8080", "./cert.pem", "./key.pem", nil))

	customServer(":3030")
}

// Instead of using a default function, we create our own.
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

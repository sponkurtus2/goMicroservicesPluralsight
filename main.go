package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	// Handle our printing function
	http.HandleFunc("/file", serveFprint("./pets.csv"))
	http.HandleFunc("/fileV2", serveServeFile("./pets.csv"))
	http.HandleFunc("/filev3", serveServeContent("./pets.csv"))

	// https://localhost:3030/files/customer.csv
	http.Handle("/files/", http.StripPrefix("/files/",
		http.FileServer(http.Dir("."))))

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
	fmt.Println("Current endpoints:  \n / \n /file")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}

	err = s.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Error -> ", err)
	}
}

// Ways to serve static content

//func serveData(fileName string) (written *os.File) {
//	dataToReturn, err := os.Open(fileName)
//	if err != nil {
//		log.Fatal("Error opening file in serveData() -> ", err)
//	}
//
//	defer func() {
//		if err := dataToReturn.Close(); err != nil {
//			panic(err)
//		}
//	}()
//
//	var dataPointer *os.File = dataToReturn
//
//	return dataPointer
//
//}

func serveFprint(fileName string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		customerFile, err := os.Open(fileName)
		if err != nil {
			log.Fatal("File not fount -> err", err)
		}

		defer func() {
			if err := customerFile.Close(); err != nil {
				panic(err)
			}
		}()

		// We only need to do this if we want to print our data using fmt.Fprint
		// This happens since io.Copy streams directly from the file and doesn't need to
		// read all the data before printing it out.
		//data, err := io.ReadAll(customerFile)
		//if err != nil {
		//	log.Fatal("Error reading data -> ", err)
		//}
		_, err = io.Copy(w, customerFile)
		if err != nil {
			log.Fatal("Error serving the data -> ", err)
		}
	}
}

func serveServeFile(filename string) func(w http.ResponseWriter, r *http.Request) {
	var contentToReturn http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
	return contentToReturn
}

// Serve Content is used to directly download a file to the device, once you enter the url associated to the Serve Content
// Function, the file will start to download
func serveServeContent(filename string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		customerFile, err := os.Open(filename)
		if err != nil {
			log.Fatal("Could'nt find file -> ", err)
		}

		defer func() {
			if err := customerFile.Close(); err != nil {
				log.Fatal("Error closing the file -> ", err)
			}
		}()

		// In order to get the modification time from our file, we first need to get the data from the file
		fileInfo, err := os.Stat(filename)
		if err != nil {
			log.Fatal("Error reading file stats -> ", err)
		}
		// Once we read all the file stats, we can get the modification time
		lastModificationTime := fileInfo.ModTime()

		http.ServeContent(w, r, "customer_data.csv", lastModificationTime, customerFile)
	}
}

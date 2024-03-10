package main

import (
	"io"
	"log"
	"net/http"
)

var customTransport = http.DefaultTransport

func init() {
	// Here, you can customize the transport, e.g., set timeouts or enable/disable keep-alive
}

func main() {

	server := http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(HandleRequest),
	}

	log.Println("Starting proxy server on :8080")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error on proxy server : ", err)
	}

}

func HandleRequest(w http.ResponseWriter, r *http.Request) {

	// Create a copy of the client request to send it to the target server

	tu := r.URL

	// Copying method url and body
	proxyReq, err := http.NewRequest(r.Method, tu.String(), r.Body)

	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	//Copying headers

	// Header example:
	// {
	// 	"header1" : ["Value 1","Value 2","Value 3","Value 4","Value 5",],
	// 	"header2" : ["Value 1","Value 2","Value 3","Value 4","Value 5",],
	// 	"header3" : ["Value 1","Value 2","Value 3","Value 4","Value 5",],
	// }

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Send the new request

	resp, err := customTransport.RoundTrip(proxyReq)

	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// the inverse : Copy the server response to send it to the client

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	// copy the body
	io.Copy(w, resp.Body)

}

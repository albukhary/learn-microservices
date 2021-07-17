package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Salam")

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Oops", http.StatusBadRequest)
		}
		log.Printf("Data %s\n", d)

	})

	http.ListenAndServe(":9000", nil)
}

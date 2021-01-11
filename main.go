package main

import (
	"awesomeProject1/src/apis/upload_api"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)
func main()  {
	router := mux.NewRouter()

	router.HandleFunc("/api/uploads", upload_api.UploadFile).Methods("POST")

	err := http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Println(err)
	}
}

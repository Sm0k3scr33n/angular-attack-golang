package main

import (
  "net/http"
  //"log"
  //"fmt"
  "github.com/gorilla/mux"
  //"github.com/gorilla/securecookie"
  //"gopkg.in/mgo.v2"
  //"gopkg.in/mgo.v2/bson"
)


func indexPageHandler(response http.ResponseWriter, request *http.Request) {

}

func pageHandler404(response http.ResponseWriter, request *http.Request) {

}

func challengePost(response http.ResponseWriter, request *http.Request) {

}

func challengeGet(response http.ResponseWriter, request *http.Request) {

}

func uploadPhotoHandler(response http.ResponseWriter, request *http.Request) {

}
func voteHandler(response http.ResponseWriter, request *http.Request) {

}





/// create a router with the gorilla mux router and handle the requests
var router = mux.NewRouter()
func main() {
    router.NotFoundHandler = http.HandlerFunc(pageHandler404)
    ///handlers for the gorilla mux router
    router.HandleFunc("/", indexPageHandler)
    router.HandleFunc("/challenge", challengePost).Methods("POST")
    router.HandleFunc("/challenge", challengeGet).Methods("GET")
    // router.HandleFunc("/location", Location).Methods("GET")
    router.HandleFunc("/uploadphoto", uploadPhotoHandler).Methods("POST")
    router.HandleFunc("/vote", voteHandler).Methods("GET")
    http.Handle("/", router)
    http.ListenAndServeTLS(":8000", "cert.pem", "key.pem", nil)
}

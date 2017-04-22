package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	//"github.com/gorilla/securecookie"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

type GeoLoc struct {
     Name string
     Lat string
     Lng string
     Url string
}

type Challenge struct {
     Place string
     Lat string
     Lng string
     Url string
     Verb string
     Noun string
}

var nouns []interface {}
var verbs []interface {}

func check(e error) {
    if e != nil {
        panic(e)
    }
}


func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func pageHandler404(response http.ResponseWriter, request *http.Request) {

}

func indexPageHandler(response http.ResponseWriter, request *http.Request) {

}

func challengePost(response http.ResponseWriter, request *http.Request) {

}

func challengeGet(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	lat := vars["lat"]
	lng := vars["lng"]

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s,%s&radius=1000&key=AIzaSyDv8MqlGr7dxQDCb7STkYGRqdCA5wuLHMM", lat, lng)
	resp, err := http.Get(url)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var f interface{}
	jsonParseErr := json.Unmarshal(body, &f)

	if jsonParseErr != nil {
		//handle error next time
	}

	m := f.(map[string]interface{})
	results := m["results"].([]interface{})
	//	var placeIDs = results[1:len(results)]

	for i := 0; i < len(results); i++ {

		result := results[i].(map[string]interface{})

		photos := result["photos"].([]interface{})
		photo := photos[0].(map[string]interface{})
		photoID := photo["photo_reference"]
		fmt.Println(photoID)
		url := photoID.(string)
		geometry := result["geometry"].(map[string]interface{})
		location := geometry["location"].(map[string]interface{})

		name := result["name"].(string)
		lat := location["lat"].(float64)
		lng := location["lng"].(float64)

		item := GeoLoc{Name: name,Lat: FloatToString(lat),Lng: FloatToString(lng), Url: url}

		results[i] = item
	}

	//fmt.Println(results)
	//fmt.Print(rand.Intn(len(results)))
	result := results[rand.Intn(len(results))].(GeoLoc)
	verbsMap := verbs[rand.Intn(len(verbs))].(map[string]interface{})

	//fmt.Println(verbsMap["present"])
	//fmt.Println(nouns[rand.Intn(len(nouns))])
	challenge := Challenge{result.Name, result.Lat, result.Lng, result.Url, verbsMap["present"].(string), nouns[rand.Intn(len(nouns))].(string)}
        //fmt.Println(result)

	b, err := json.Marshal(challenge)

	//fmt.Println(b)

	response.Write(b)
}

func uploadPhotoHandler(response http.ResponseWriter, request *http.Request) {

}

func voteHandler(response http.ResponseWriter, request *http.Request) {

}

func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + ":8000" + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		// see @andreiavrammsd comment: often 307 > 301
		http.StatusTemporaryRedirect)
}

func appVersion(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, I am version 0.0.1 %q", html.EscapeString(req.URL.Path))
}

/// create a router with the gorilla mux router and handle the requests
var router = mux.NewRouter()

func main() {

	dat, err := ioutil.ReadFile("./nouns.json")
	check(err)
	fmt.Println(string(dat))

        var n interface{}
        nounJsonParseErr := json.Unmarshal(dat, &n)
	check(nounJsonParseErr)

        nounsMap := n.(map[string]interface{})
	nouns = nounsMap["nouns"].([]interface {})
        fmt.Println(nouns)

	var v interface{}
	verbDat, verbIOerr := ioutil.ReadFile("./verbs.json")
	check(verbIOerr)
	verbJsonParseErr := json.Unmarshal(verbDat, &v)
        check(verbJsonParseErr)

        verbsMap := v.(map[string]interface{})
        verbs = verbsMap["verbs"].([]interface {})
        fmt.Println(verbs)

	router.NotFoundHandler = http.HandlerFunc(pageHandler404)
	///handlers for the gorilla mux router
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/challenge", challengePost).Methods("POST")
	router.HandleFunc("/challenge/{lat}/{lng}", challengeGet).Methods("GET")
	// router.HandleFunc("/location", Location).Methods("GET")
	router.HandleFunc("/uploadphoto", uploadPhotoHandler).Methods("POST")
	router.HandleFunc("/vote", voteHandler).Methods("POST")
	http.Handle("/", router)
	port := ":80"
	go http.ListenAndServe(port, http.HandlerFunc(appVersion))
	// port2 := ":443"
	// go http.ListenAndServeTLS(port2,"cert.pem", "key.pem", http.HandlerFunc(appVersion))
	http.ListenAndServeTLS(":443", "../fullchain.pem", "../privkey.pem", nil)
}

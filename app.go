package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"html"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	//"github.com/gorilla/securecookie"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

const TLS_CERT_FILE = "/root/fullchain.pem"
const TLS_KEY_FILE = "/root/privkey.pem"

type GeoLoc struct {
	Name string
	Lat  string
	Lng  string
	Url  string
}

type Challenge struct {
	Place string
	Lat   string
	Lng   string
	Url   string
	Verb  string
	Noun  string
}

// Structure for class variables and members for use
// in handlers.
type AppObject struct {
	dal DataLayerInterface
}

type AppInterface interface {
	Close()
	appVersion(w http.ResponseWriter, req *http.Request)
	voteHandler(response http.ResponseWriter, request *http.Request)
	pageHandler404(response http.ResponseWriter, request *http.Request)
	indexPageHandler(response http.ResponseWriter, request *http.Request)
	challengePost(response http.ResponseWriter, request *http.Request)
	challengeGet(response http.ResponseWriter, request *http.Request)
	uploadPhotoHandler(response http.ResponseWriter, request *http.Request)
}

func NewApp() AppInterface {
	dal := NewDataLayer()
	obj := AppObject{dal}
	return &obj
}

var nouns []interface{}
var verbs []interface{}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (app *AppObject) pageHandler404(response http.ResponseWriter, request *http.Request) {

}

func (app *AppObject) indexPageHandler(response http.ResponseWriter, request *http.Request) {

}

func (app *AppObject) challengePost(response http.ResponseWriter, request *http.Request) {

}

func (app *AppObject) uploadPhotoHandler(response http.ResponseWriter, request *http.Request) {

}

func (app *AppObject) voteHandler(response http.ResponseWriter, request *http.Request) {

}

// func redirect(w http.ResponseWriter, req *http.Request) {
// 	// remove/add not default ports from req.Host
// 	target := "https://" + req.Host + ":8000" + req.URL.Path
// 	if len(req.URL.RawQuery) > 0 {
// 		target += "?" + req.URL.RawQuery
// 	}
// 	log.Printf("redirect to: %s", target)
// 	http.Redirect(w, req, target,
// 		// see @andreiavrammsd comment: often 307 > 301
// 		http.StatusTemporaryRedirect)
// }

func appVersion(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, I am version 0.0.1 %q", html.EscapeString(req.URL.Path))
}

func (app *AppObject) Close() {
	// close the connection to the database when the app instance destructs
	app.dal.Close()
}

func (app *AppObject) appVersion(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, I am version 0.0.1 %q", html.EscapeString(req.URL.Path))
}

func (app *AppObject) challengeGet(response http.ResponseWriter, request *http.Request) {
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
	fmt.Println(string(body))
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
		photoID := ""
		if result["photos"] != nil {
			photos := result["photos"].([]interface{})
			photo := photos[0].(map[string]interface{})
			photoID = photo["photo_reference"].(string)
			fmt.Println(photoID)
		}

		url := photoID
		geometry := result["geometry"].(map[string]interface{})
		location := geometry["location"].(map[string]interface{})

		name := result["name"].(string)
		lat := location["lat"].(float64)
		lng := location["lng"].(float64)

		item := GeoLoc{Name: name, Lat: FloatToString(lat), Lng: FloatToString(lng), Url: url}

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

	app.dal.SaveChallenge(challenge)

	b, err := json.Marshal(challenge)

	//fmt.Println(b)

	response.Write(b)
}

func start80(wg *sync.WaitGroup, handler http.HandlerFunc) {
	log.Println("starting server on 80")
	defer wg.Done()
	err := http.ListenAndServe(":80", handler)
	log.Println(err)
	log.Println("done with 80")
}

func start443(wg *sync.WaitGroup, handler http.Handler) {
	log.Println("starting server on 443")
	defer wg.Done()
	err := http.ListenAndServeTLS(":443", TLS_CERT_FILE, TLS_KEY_FILE, handler)
	log.Println(err)
	log.Println("done with 443")
}

func main() {
	// create application state and connect to database
	var app = NewApp()
	defer app.Close()

	dat, err := ioutil.ReadFile("./nouns.json")
	check(err)
	// fmt.Println(string(dat))

	var n interface{}
	nounJsonParseErr := json.Unmarshal(dat, &n)
	check(nounJsonParseErr)

	nounsMap := n.(map[string]interface{})
	nouns = nounsMap["nouns"].([]interface{})
	// fmt.Println(nouns)

	var v interface{}
	verbDat, verbIOerr := ioutil.ReadFile("./verbs.json")
	check(verbIOerr)
	verbJsonParseErr := json.Unmarshal(verbDat, &v)
	check(verbJsonParseErr)

	verbsMap := v.(map[string]interface{})
	verbs = verbsMap["verbs"].([]interface{})
	// fmt.Println(verbs)

	/// create a router with the gorilla mux router and handle the requests
	var router = mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(app.pageHandler404)
	///handlers for the gorilla mux router
	router.HandleFunc("/", app.indexPageHandler)
	router.HandleFunc("/challenge", app.challengePost).Methods("POST")
	router.HandleFunc("/challenge/{lat}/{lng}", app.challengeGet).Methods("GET")
	// router.HandleFunc("/location", Location).Methods("GET")
	router.HandleFunc("/uploadphoto", app.uploadPhotoHandler).Methods("POST")
	router.HandleFunc("/vote", app.voteHandler).Methods("POST")
	http.Handle("/", router)

	// corsOpts := cors.Default()
	corsOpts := cors.New(cors.Options{
		Debug: false,
	})
	handler := corsOpts.Handler(router)

	// start http and https servers
	var wg sync.WaitGroup
	wg.Add(2)
	go start80(&wg, http.HandlerFunc(app.appVersion))
	go start443(&wg, handler)
	wg.Wait()
}

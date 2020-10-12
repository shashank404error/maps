package main

import (
	//"fmt"
	"net/http"
	"encoding/json"
	//"math"
	//"io/ioutil"
	"log"
	//"reflect"
	"github.com/gorilla/mux"
	//"github.com/shashank404error/account"
	"github.com/shashank404error/shashankMongo"
	"html/template"
	"os"
)

var connectDBInfo *shashankMongo.ConnectToDataBase = &shashankMongo.ConnectToDataBase{
	CustomApplyURI:"mongodb://shashank404error:Y9ivXgMQ5ZrjL4N@parkpoint-shard-00-00.0bxqn.mongodb.net:27017,parkpoint-shard-00-01.0bxqn.mongodb.net:27017,parkpoint-shard-00-02.0bxqn.mongodb.net:27017/parkpoint?ssl=true&replicaSet=atlas-21pobg-shard-0&authSource=admin&retryWrites=true&w=majority", 
	DatabaseName:"parkpoint", 
	CollectionName:"businessAccounts", 
}

var templates *template.Template

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == ""{
	  port = "8000"
	}
	return ":" + port, nil
  }

func main(){

	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	templates = template.Must(template.ParseGlob("templates/*.gohtml"))
	r := mux.NewRouter()
	r.HandleFunc("/", showIndexPage).Methods("GET")
	r.HandleFunc("/create/account/{userName}/{businessName}/{password}/{city}", createAccount).Methods("POST")
	r.HandleFunc("/create/profile/{userID}/{plan}", createProfile).Methods("POST")
	//r.HandleFunc("/login/{userName}/{password}", loginAccount).Methods("POST")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",fs))
	http.Handle("/",r)
	if err := http.ListenAndServe(addr, nil);
	err != nil {
		panic(err)
	  }
}

func showIndexPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.gohtml", nil)
}

func createAccount(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	load:=`{
		"username":"`+vars["userName"]+`",
		"businessname": "`+vars["businessName"]+`",
		"password": "`+vars["password"]+`",
		"city": "`+vars["city"]+`"
		}`

	loadToJson:=byteToJsonInterface(load)
	
	id:=shashankMongo.InsertOne(connectDBInfo,loadToJson)
	templates.ExecuteTemplate(w, "selectPlan.gohtml", id)
}

func createProfile(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key:="businessplan"
	value:=vars["plan"]

	res:=shashankMongo.UpdateOneByID(connectDBInfo,vars["userID"],key,value)
	if(res==1){
		connectDBInfo.CollectionName = "profileConfig"
		config:=shashankMongo.FetchProfileConfiguration(connectDBInfo,value)
		connectDBInfo.CollectionName = "businessAccounts"
		profileRes := shashankMongo.UpdateProfileConfiguration(connectDBInfo,vars["userID"],config)
		if(profileRes==1) {
			userConfig:=shashankMongo.FetchProfile(connectDBInfo,vars["userID"])
			templates.ExecuteTemplate(w, "profile.gohtml", userConfig)
		}
	}
}

func byteToJsonInterface(load string) map[string]interface{} {
	var loadArr = []byte(load)
    var loadToJson map[string]interface{}
    err := json.Unmarshal(loadArr, &loadToJson)
    if (err != nil) {
		log.Fatal(err)
	}
	return loadToJson
}
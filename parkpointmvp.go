package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	//"math"
	//"io/ioutil"
	"log"
	//"reflect"
	"github.com/gorilla/mux"
	//"github.com/shashank404error/account"
	"github.com/shashank404error/shashankMongo"
	"github.com/shashank404error/middlework"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"os"
)

var connectDBInfo *shashankMongo.ConnectToDataBase = &shashankMongo.ConnectToDataBase{
	CustomApplyURI:"mongodb://shashank404error:Y9ivXgMQ5ZrjL4N@parkpoint-shard-00-00.0bxqn.mongodb.net:27017,parkpoint-shard-00-01.0bxqn.mongodb.net:27017,parkpoint-shard-00-02.0bxqn.mongodb.net:27017/parkpoint?ssl=true&replicaSet=atlas-21pobg-shard-0&authSource=admin&retryWrites=true&w=majority", 
	DatabaseName:"parkpoint",  
}

var templates *template.Template

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == ""{
	  port = "80"
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
	r.HandleFunc("/login/{userName}/{password}", loginAccount).Methods("POST")
	r.HandleFunc("/overview/{userID}", loadOverview).Methods("POST")
	r.HandleFunc("/zones/{userID}", loadZone).Methods("POST")
	r.HandleFunc("/zone/assign/{UserID}", assignToZone).Methods("POST")
	r.HandleFunc("/tracking/{UserID}", liveTracking).Methods("POST")
	r.HandleFunc("/liveTracking/{UserID}", getLiveLocation).Methods("POST")
	
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
		"city": "`+vars["city"]+`",
		"deliveryPending":"0",
		"deliveryDelivered":"0"
		}`

	loadToJson:=byteToJsonInterface(load)
	
	id:=shashankMongo.InsertOne(connectDBInfo,"businessAccounts",loadToJson)
	templates.ExecuteTemplate(w, "selectPlan.gohtml", id)
}

func createProfile(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key:="businessplan"
	value:=vars["plan"]

	res:=shashankMongo.UpdateOneByID(connectDBInfo,"businessAccounts",vars["userID"],key,value)
	if(res==1){
		config:=shashankMongo.FetchProfileConfiguration(connectDBInfo,"profileConfig",value)
		profileRes := shashankMongo.UpdateProfileConfiguration(connectDBInfo,"businessAccounts",vars["userID"],config)
		middlework.CreateZones(connectDBInfo,"parking",vars["userID"],config)
		if(profileRes==1) {
			userConfig:=shashankMongo.FetchProfile(connectDBInfo,"businessAccounts",vars["userID"])		
			templates.ExecuteTemplate(w, "profile.gohtml", userConfig)
		}
	}
}

func loginAccount(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	username:=vars["userName"]
	password:=vars["password"]
	userConfig,err := shashankMongo.FetchLogin(connectDBInfo,"businessAccounts",username,password)
	if err!=nil {
		error:="Username and password donot match! Please close this tab and try again."
		templates.ExecuteTemplate(w, "error.gohtml", error)	
	}else{
	templates.ExecuteTemplate(w, "profile.gohtml", userConfig)
	}
}

func loadZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	account:=shashankMongo.GetZone(connectDBInfo,"parking",vars["userID"])
	templates.ExecuteTemplate(w, "zone.gohtml", account)
}

func loadOverview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userConfig:=shashankMongo.FetchProfile(connectDBInfo,"businessAccounts",vars["userID"])
	templates.ExecuteTemplate(w, "profile.gohtml", userConfig)
}

func assignToZone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	
	file, _, err := r.FormFile("xlsxFile")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
    }
	defer file.Close()
	
	deliveryArrUpdate,noOfDeliveryUpdate,userId:=middlework.UploadToExcel(file,connectDBInfo,"parking",vars["UserID"])
	if(deliveryArrUpdate == 1 && noOfDeliveryUpdate == 1 && userId!= ""){
		userConfig:=shashankMongo.FetchProfile(connectDBInfo,"businessAccounts",userId)
		templates.ExecuteTemplate(w, "profile.gohtml", userConfig)
	}

}

func liveTracking(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	var deliveredDetails []shashankMongo.DeliveredDetail
	var deliveredDetail shashankMongo.DeliveredDetail
	documents:=shashankMongo.GetFieldByFilter (connectDBInfo, "delivered", "businessUid", vars["UserID"])

	for _,val:= range documents{
		bsonBytes, _ := bson.Marshal(val)
		bson.Unmarshal(bsonBytes, &deliveredDetail)
		deliveredDetails = append(deliveredDetails,deliveredDetail)	
	} 
	account:=shashankMongo.FetchProfile(connectDBInfo,"businessAccounts", vars["UserID"])
	load:= shashankMongo.DeliveredAndAccount{
		DeliveredDetails: deliveredDetails,
		BusinessAccount: account,
	}
	templates.ExecuteTemplate(w, "pastorders.gohtml", load)
}

func getLiveLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profile:=shashankMongo.FetchProfile(connectDBInfo, "businessAccounts", vars["UserID"])
	templates.ExecuteTemplate(w, "track.gohtml", profile)
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
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)



func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://localhost")
if err != nil {
	panic(err)
}
return session
}
var collection = getSession().DB("peliculas").C("movies")

func responseMovie( rw http.ResponseWriter, status int,results Movie ){
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(results)
}

func responseMovies( rw http.ResponseWriter, status int,results []Movie ){
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(results)
}

func Index (rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(rw,"Hola mundo")
	}

func MovieList (rw http.ResponseWriter, r *http.Request) {
		var results []Movie
		err := collection.Find(nil).Sort("_id").All(&results)
		if err !=nil {
			log.Fatal(err)
			
		}else{
			fmt.Println("resultados: ", results)
		}
		responseMovies (rw,200,results)
		}

func MovieShow (rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieId := params["id"]

		if !bson.IsObjectIdHex(movieId){
			rw.WriteHeader(404)
			return
		}	
		oid := bson.ObjectIdHex(movieId)
		results := Movie{}
		err := collection.FindId(oid).One(&results)
		if err!=nil {
			rw.WriteHeader(404)
			return			
		}
		responseMovie(rw,200,results)
		
		}

func MovieAdd(rw http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  var movieData Movie
  err := decoder.Decode(&movieData)

  if (err != nil) {
	  panic(err)	  
  }
  defer r.Body.Close()

err = collection.Insert(movieData)

if err !=nil {
	rw.WriteHeader(500)
	
}
responseMovie(rw,200,movieData)
}
func MovieUpdate (rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieId := params["id"]

		if !bson.IsObjectIdHex(movieId){
			rw.WriteHeader(404)
			return
		}
		oid := bson.ObjectIdHex(movieId)
		decoder := json.NewDecoder(r.Body)	
		var movieData Movie
		err := decoder.Decode(&movieData)
		 if err != nil{
			 panic(err)
			 rw.WriteHeader(500)
			 return
		 }

		 defer r.Body.Close()

		document := bson.M{"_id":oid}
		change := bson.M{"$set":movieData}
		err = collection.Update(document,change)
		if err != nil{
			panic(err)
			rw.WriteHeader(404)
			return
		}
		responseMovie(rw,200,movieData)
		
		}

type Message struct{
			Status string `json:"status"`
			Message string `json:"mensaje"`
		}
func MovieRemove (rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	movieId := params["id"]
		
		if !bson.IsObjectIdHex(movieId){
			rw.WriteHeader(404)
			return
		}	
		oid := bson.ObjectIdHex(movieId)
		err := collection.RemoveId(oid)
		if err!=nil {
			rw.WriteHeader(404)
			return			
		}
		results := Message {"success", "la pelicula con ID:"+ movieId+"ha sido eliminada correctamente."}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(results)
}

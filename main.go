package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	// "strings"
	"time"

	"github.com/OLUWAMUYIWA/Adel/api"
	"github.com/codegangsta/negroni"

	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var sec = os.Getenv("secret")

// var Signature = []byte(sec)
// var sec = "984fv873rfnvfo9u34rb34340b5geor08343otf89wfbw4893"
var mySigningKey = []byte(sec)



func main() {
	// MongooDB connection

	var connStr = os.Getenv("dbconn")
	var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
	// clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	// var connStr = "mongodb+srv://adel:l7hbnRiL7BEr1bck@drugs.u6k4q.mongodb.net/drugs?retryWrites=true&w=majority"
	clientOptions := options.Client().ApplyURI(connStr)
	
	client, err := mongo.Connect(ctx, clientOptions)
    if err != nil {
        log.Fatalf("this error fucking appened: %v",err)
    }
    err = client.Ping(ctx, readpref.PrimaryPreferred())
    if err != nil { 
        log.Fatalf("while doing shit: %v", err)
	}
	defer client.Disconnect(ctx)
	dBase := client.Database("drugstore")

	drugsColl := dBase.Collection("drugs")
	
	mod := mongo.IndexModel {
			Keys: bson.M{
				"name": 1,
			}, Options: nil,
		}
	if _, err := drugsColl.Indexes().CreateOne(ctx, mod); err != nil {
		log.Println("Could not create index:", err)
	}

	// Expiration index for drugs

	var ttl = int32(0)
	ttl = 60
    keys := bsonx.Doc{{"expiry_date", bsonx.Int32(int32(1))}}
    idx := mongo.IndexModel{Keys: keys, Options: &options.IndexOptions{ExpireAfterSeconds:  &ttl}}
    _, err = drugsColl.Indexes().CreateOne(ctx, idx)
    if err != nil {
        log.Println("Error occurred while creating index", err)
    } else {
        log.Println("Index creation success")
    }

	//Indexes for User

	collSeniors := dBase.Collection("seniors")
	optsSen := options.CreateIndexes().SetMaxTime(3 * time.Second)
	indexesSen := []mongo.IndexModel{}
	indexStringsSen := []string{"email", "password"}
	  for _, val := range indexStringsSen {
		temp := mongo.IndexModel{}
		temp.Keys = bsonx.Doc{{Key: val, Value: bsonx.Int32(int32(1))}}
		indexesSen = append(indexesSen, temp)
	  }
	resp, err := collSeniors.Indexes().CreateMany(context.Background(), indexesSen, optsSen)
	if err != nil {
		log.Print("senior index was not created")
	}
	log.Print(resp)

	collJuniors := dBase.Collection("juniors")
	opts := options.CreateIndexes().SetMaxTime(3 * time.Second)
	indexes := []mongo.IndexModel{}
	indexStrings := []string{"email", "password"}
	  for _, val := range indexStrings {
		temp := mongo.IndexModel{}
		temp.Keys = bsonx.Doc{{Key: val, Value: bsonx.Int32(int32(1))}}
		indexes = append(indexes, temp)
	  }
	res, err := collJuniors.Indexes().CreateMany(context.Background(), indexes, opts)
	log.Print(res)
	r := mux.NewRouter().StrictSlash(false)
	
	s := r.PathPrefix("/api").Subrouter()

	noAuth := s.PathPrefix("/no_auth").Subrouter()
	noAuth.Path("/regJunior").HandlerFunc(api.CreateJunior( dBase)).Methods("POST", "OPTIONS")
	noAuth.Path("/regSenior").HandlerFunc(api.CreateSenior(dBase)).Methods("POST", "OPTIONS")
	noAuth.Path("/login").HandlerFunc(api.LoginHandler(dBase)).Methods("POST", "OPTIONS")
	noAuth.Path("/regBoss").HandlerFunc(api.CreateBosss(dBase)).Methods("POST", "OPTIONS")

	s.Handle("/search/{name}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareAll),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.Search(dBase))),
	)).Methods("GET", "OPTIONS")
	
	s.Handle("/junUpdate/{uid}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareJunior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.UpdateJunior(dBase))),
	)).Methods("PUT", "OPTIONS")

	
	s.Handle("/uploadManyDrugs/{uid}/{cname}/{cphone}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.UploadMany(dBase))),
	)).Methods("POST", "OPTIONS")

	s.Handle("/uploadDrug/{uid}/{cname}/{cphone}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.Upload( dBase))),
	)).Methods("POST", "OPTIONS")

	s.Handle(("/sendMyDrugs/{uid}/{cname}"),
		negroni.New(
			negroni.NewRecovery(),
			negroni.HandlerFunc(api.AuthorizeWareSenior),
			negroni.NewLogger(),
			negroni.Wrap(http.HandlerFunc(api.SendMyDrugs(dBase))) ,
		)).Methods("GET", "OPTIONS")
		
	s.Handle("/updateMyDrugs/{uid}/{cname}", 
		negroni.New(
			negroni.NewRecovery(),
			negroni.HandlerFunc(api.AuthorizeWareSenior),
			negroni.NewLogger(),
			negroni.Wrap(http.HandlerFunc(api.UpdateMyDrugs(dBase))),
		)).Methods("PUT", "OPTIONS")

	s.Handle("/updateSenior/{uid}/{cname}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.UpdateSenior(dBase))),
	)).Methods("PUT", "OPTIONS")
	
	s.Handle("/deleteDrug/{uid}/{id}/{cname}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.DeleteDrug(dBase))),
	)).Methods("DELETE", "OPTIONS")

	

	s.Handle("/sendUnverifiedJuniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.SendUnverifiedJuniors(dBase))),
	)).Methods("GET", "OPTIONS")
	
	s.Handle("/sendUnverifiedSeniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.SendUnverifiedSeniors(dBase))),
	)).Methods("GET", "OPTIONS")
	
	s.Handle("/verifyManyJuniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.VerifyManyJuniors(dBase))),
	)).Methods("PUT", "OPTIONS")

	s.Handle("/verifyManySeniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.VerifyManySeniors(dBase))),
	)).Methods("PUT", "OPTIONS")
	
	s.HandleFunc("/servestatic/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("we came here to conquer"))
	})

	clientApp := http.Dir("./clientdist/index.html")
	catchAll := "/"
	fs := http.FileServer(clientApp)
	r.Handle(catchAll, fs)

	// http.ListenAndServe(":3000", r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		panic(err)
	}
	
}


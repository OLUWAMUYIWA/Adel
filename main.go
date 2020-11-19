package main
import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/OLUWAMUYIWA/Adel/api"
	//jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	//"github.com/gorilla/handlers"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)




var mySigningKey = []byte("984fv873rfnvfo9u34rb34340b5geor08343otf89wfb")

// type Connection struct {
// 	Seniors		*mongo.Collection
// 	Juniors		*mongo.Collection
// 	Drugs		*mongo.Collection
// }


func main() {
	// MongooDB connection
        JavascriptISOString := "2006-01-02T15:04:05.999Z07:00"
	time.Now().UTC().Format(JavascriptISOString)
	var ctx, _  = context.WithTimeout(context.Background(), 200 * time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	
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
	// var Conn Connection = Connection {
	// 	Seniors: dBase.Collection("seniors"),
	// 	Juniors: dBase.Collection("juniors"),
	// 	Drugs: dBase.Collection("drugs"),

	// }

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
	//r.Host("www.example.com")

	noAuth := s.PathPrefix("/no_auth").Subrouter()
	noAuth.Path("/regJunior").HandlerFunc(api.CreateJunior(ctx, dBase)).Methods("POST")
	noAuth.Path("/regSenior").HandlerFunc(api.CreateSenior(ctx, dBase)).Methods("POST")
	noAuth.Path("/login").HandlerFunc(api.LoginHandler(ctx, dBase)).Methods("POST")
	noAuth.Path("/regBoss").HandlerFunc(api.CreateBosss(ctx, dBase)).Methods("POST")

	s.Handle("/search/{name}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareAll),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.Search(ctx, dBase))),
	)).Methods("GET")
	
	s.Handle("/junUpdate/{uid}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareJunior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.UpdateJunior(ctx, dBase))),
	)).Methods("PUT")

	
	s.Handle("/uploadManyDrugs/{uid}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.UploadMany(ctx, dBase))),
	)).Methods("POST")

	s.Handle("/uploadDrug/{uid}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.Upload(ctx, dBase))),
	)).Methods("POST")

	s.Handle(("/sendMyDrugs/{uid}"),
		negroni.New(
			negroni.NewRecovery(),
			negroni.HandlerFunc(api.AuthorizeWareSenior),
			negroni.NewLogger(),
			negroni.Wrap(http.HandlerFunc(api.SendMyDrugs(ctx, dBase))) ,
		)).Methods("GET")
		
	s.Handle("/updateMyDrugs/{uid}", 
		negroni.New(
			negroni.NewRecovery(),
			negroni.HandlerFunc(api.AuthorizeWareSenior),
			negroni.NewLogger(),
			negroni.Wrap(http.HandlerFunc(api.UpdateMyDrugs(ctx, dBase))),
		)).Methods("PUT")

	s.Handle("/updateSenior/{uid}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.UpdateSenior(ctx, dBase))),
	)).Methods("PUT")
	
	s.Handle("/deleteDrug/{uid}/{id}", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareSenior),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.DeleteDrug(ctx, dBase))),
	)).Methods("DELETE")

	s.Handle("/sendUnverifiedJuniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.SendUnverifiedJuniors(ctx, dBase))),
	)).Methods("GET")
	
	s.Handle("/sendUnverifiedSeniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.SendUnverifiedSeniors(ctx, dBase))),
	)).Methods("GET")
	
	s.Handle("/verifyManyJuniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.VerifyManyJuniors(ctx, dBase))),
	)).Methods("PUT")

	s.Handle("/verifyManySeniors", negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(api.AuthorizeWareBoss),
		negroni.NewLogger(),
		negroni.Wrap(http.HandlerFunc(api.VerifyManySeniors(ctx, dBase))),
	)).Methods("PUT")

	http.ListenAndServe(":3000", r)
	//http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r))

}

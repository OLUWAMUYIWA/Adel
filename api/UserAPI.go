package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/OLUWAMUYIWA/Adel/data"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type VerifiedPlusUnverified struct {
	Juniors 	[]data.Junior
	Seniors 	[]data.Senior
}

type Message struct {
	Message string `json:"message"`
}


//CreateJunior creates a junior
func CreateJunior(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		//w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		var junior data.Junior
		err := json.NewDecoder(r.Body).Decode(&junior)
		if err != nil {
			ErrorWithJSON(w, Message{"Error while decoding json"}, http.StatusBadRequest)
		}
		coll := base.Collection("juniors")
		junior.Email = strings.ToLower(junior.Email)
		filter := bson.M{"email": junior.Email}
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		no, err := coll.CountDocuments(ctx, filter)
		if no > 0 {
			w.WriteHeader(http.StatusUnauthorized)
			ErrorWithJSON(w, Message{"You have registered before"}, http.StatusUnauthorized)
			return
		}
		junior.Role = "junior"
		junior.Verified = false
		junior.TimeRegd = time.Now()
		//junior.Email = "junior" + junior.Email
		junior.Id = primitive.NewObjectID()
		fmt.Printf("jun: %v", junior)
		_, err = coll.InsertOne(ctx, junior)
		if err != nil {
			ErrorWithJSON(w, Message{err.Error()}, http.StatusInternalServerError)
			log.Println("Failed create junior ", err)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		m := Message{"Created Successfully"}
		response, err := json.Marshal(m)
		w.Write([]byte(response))
		
	}
}

//CreateSenior creates a senior
func CreateSenior(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		var senior data.Senior
		err := json.NewDecoder(r.Body).Decode(&senior)

		if err != nil {
			ErrorWithJSON(w, Message{"Error while decoding json"}, http.StatusBadRequest)
			return
		}
		coll := base.Collection("seniors")
		senior.Email = strings.ToLower(senior.Email)
		senior.Role = "senior"
		filter := bson.M{"email": senior.Email}
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		no, err := coll.CountDocuments(ctx, filter)
		if no > 0 {
			w.WriteHeader(http.StatusUnauthorized)
			ErrorWithJSON(w, Message{"You have registered before"}, http.StatusUnauthorized)
			return
		}
		senior.Verified = false
		senior.TimeRegd = time.Now()
		senior.Id= primitive.NewObjectID()
		
		_ , err = coll.InsertOne(ctx, senior)
		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Println("Failed create senior ", err)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusOK)
		m := Message{"Created Successfully"}
		response, err := json.Marshal(m)
		w.Write([]byte(response))
	}
}
func CreateBosss(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		var boss data.Boss
		err := json.NewDecoder(r.Body).Decode(&boss)
		if err != nil {
			ErrorWithJSON(w, Message{"Error while decoding json"}, http.StatusBadRequest)
			return
		}
		boss.Email = strings.ToLower(boss.Email)
		boss.Role = "boss"
		boss.CompanyName = "Adel"
		boss.Verified = true
		coll := base.Collection("bosses")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		_ , err = coll.InsertOne(ctx, boss)
		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Println("Failed create boss: ", err)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusOK)
		m := Message{"Created Successfully"}
		response, err := json.Marshal(m)
		w.Write([]byte(response))
	}
}

//UpdateJunior updates a junior
func UpdateJunior(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		idStr := vars["uid"]
		id, _ := primitive.ObjectIDFromHex(idStr)
		var junior data.Junior

		err := json.NewDecoder(r.Body).Decode((&junior))
		if err != nil {
			ErrorWithJSON(w, Message{"Bad request body"}, http.StatusBadRequest)
			return
		}
		junior.Email = strings.ToLower(junior.Email)
		coll := base.Collection("juniors")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		res, err := coll.UpdateOne(
			ctx, 
			bson.M{"_id": id},
			bson.D{
				{"$set", junior},
			},
		)

		if err != nil {
				ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
				log.Println("Failed update junior: ", err)
				return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusOK)
		log.Print(res)
		jsonM, _ := json.Marshal("Upate done")
		w.Write(jsonM)
		}

	}

//UpdateSenior updates a senior
func UpdateSenior(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
		return func (w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		idStr := vars["uid"]
		id, _ := primitive.ObjectIDFromHex(idStr)
		var senior data.Senior
		err := json.NewDecoder(r.Body).Decode((&senior))
		if err != nil {
			ErrorWithJSON(w, Message{"Bad request body"}, http.StatusBadRequest)
			return
		}
		senior.Email = strings.ToLower(senior.Email)
		coll := base.Collection("seniors")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		res, err := coll.UpdateOne(
			ctx, 
			bson.M{"_id": id},
			bson.D{
				{"$set", senior},
			},
		)

		if err != nil {
				ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
				log.Println("Failed update senior: ", err)
				return
		}
		//w.WriteHeader(http.StatusOK)
		
		jsonM, _ := json.Marshal(res.ModifiedCount)
		w.Write(jsonM)
		}
}

func SendUnverifiedJuniors (base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		coll := base.Collection("juniors")
		opt := options.Find()
		opt.SetSort(bson.D{{"time_regd", -1}})
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		filter := bson.D{{"verified", bson.D{{"$eq", false}}}}
		jFilterCursor, err := coll.Find(ctx, filter, opt)
		juniors := []data.Junior {}
		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Println("Failed to finnd unverified users")
		}
		defer jFilterCursor.Close(ctx)
		for jFilterCursor.Next(ctx) {
			var junior data.Junior
			err := jFilterCursor.Decode(&junior)
			if err != nil {
				log.Print(err)
			}
			//junior.Email = strings.TrimPrefix(junior.Email, "junior")
			juniors = append(juniors, junior)


		}
		

		resp, er := json.Marshal(juniors)
		if er != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, resp, http.StatusOK)
	}	
	
	}	

func SendUnverifiedSeniors(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		coll2 := base.Collection("seniors")
		filter := bson.D{{"verified", bson.D{{"$eq", false}}}}
		opt := options.Find()
		opt.SetSort(bson.D{{"time_regd", -1}})
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		sFilterCursor, err := coll2.Find(ctx, filter, opt)
		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Println("Failed to finnd unverified users")
		}
		seniors := []data.Senior{}
		defer sFilterCursor.Close(ctx)
		for sFilterCursor.Next(ctx) {
			var senior data.Senior
			err := sFilterCursor.Decode(&senior)
			if err != nil {
				log.Print(err)
			}
			//senior.Email = strings.TrimPrefix(senior.Email, "senior")
			seniors = append(seniors, senior)

		}
		resp, er := json.Marshal(seniors)
		if er != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, resp, http.StatusOK)
	}
}

func VerifyManyJuniors (base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		var all []string
		if err := json.NewDecoder(r.Body).Decode(&all); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		jColl := base.Collection("juniors")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		
		for _,j := range all {
			id, _ := primitive.ObjectIDFromHex(j)
			
			_, err := jColl.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", bson.M{"verified": true}}})
			if err != nil {
				fmt.Fprintf(w, "error while checking db", err)
				return
			}
		}

		//w.WriteHeader(http.StatusOK)
		ret := Message{"All Verified"}
		returned, _ := json.Marshal(ret)
		w.Write(returned)
	}
}

func VerifyManySeniors (base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		var all []string
		if err := json.NewDecoder(r.Body).Decode(&all); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		sColl := base.Collection("seniors")
		log.Printf("we got here")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		for _,j := range all {
			id, _ := primitive.ObjectIDFromHex(j)
			_, err := sColl.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", bson.M{"verified": true}}})
			if err != nil {
				fmt.Fprintf(w, "error while checking db", err)
				return
			}
		}
		//w.WriteHeader(http.StatusOK)
		ret := Message{"All Verified"}
		returned, _ := json.Marshal(ret)
		w.Write(returned)
	}
}
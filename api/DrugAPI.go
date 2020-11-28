package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/OLUWAMUYIWA/Adel/data"
	"github.com/gorilla/mux"

	//."github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db string = "drugstore"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
//ResponseWithJSON is used in most handlers
func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.WriteHeader(code)
	w.Write(json)
}

//ErrorWithJSON is used in many handlers to print errors
func ErrorWithJSON(w http.ResponseWriter, message Message, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}


//Up handles the drug picture upload
func Up(w http.ResponseWriter, r *http.Request) {

	drugId := r.FormValue("drugId")

	file, header, err := r.FormFile("	")
	defer file.Close()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	filename := path.Join("uploads", drugId + path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
	io.WriteString(w, err.Error())
	return
	}
	io.WriteString(w, "Successful")
	
} 


//Upload caters to the upload of drugs
func Upload(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request)  {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		uploaderID := vars["uid"]
		cname := vars["cname"]
		cphone := vars["cphone"]
		var drug data.Drug
		err := json.NewDecoder(r.Body).Decode(&drug)
		if err != nil {
			ErrorWithJSON(w, Message{"Incorrect body"}, http.StatusBadRequest)
			return
		}
		uId, err := primitive.ObjectIDFromHex(uploaderID)
		drug.UploaderID = uId
		drug.CompPhone = cphone
		drug.Name = strings.ToLower(drug.Name)
		drug.TimeUploaded = time.Now()
		drug.ExpiryDate = drug.TimeUploaded.AddDate(0, drug.ExpiryMonth, 0)

		// uFilter := bson.M{"_id": uId}
		// uColl := base.Collection("seniors")
		// curs := uColl.FindOne(ctx, uFilter)
		// var sen data.Senior
		// err = curs.Decode(&sen)
		// if err != nil {
		// 	log.Print(err)
		// }

		drug.CompanyName = cname

		drugsColl := base.Collection("drugs")

		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		_, err = drugsColl.InsertOne(ctx, drug)
		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Print(err)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusOK)
		jsonRes, err := json.Marshal(Message{"Drug uploaded"})
		w.Write([]byte(jsonRes))

	}

}
func UploadMany(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request)  {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		uploaderID := vars["uid"]
		cname := vars["cname"]
		cphone := vars["cphone"]
		var drugs []data.Drug
		err := json.NewDecoder(r.Body).Decode(&drugs)
		if err != nil {
			log.Printf("error decode: %s", err)
			ErrorWithJSON(w, Message{"Incorrect body"}, http.StatusBadRequest)
			return
		}
		upId, err := primitive.ObjectIDFromHex(uploaderID)
		// uFilter := bson.M{"_id": upId}
		// uColl := base.Collection("seniors")
		// curs := uColl.FindOne(ctx, uFilter)
		// var sen data.Senior
		// err = curs.Decode(&sen)


		//j := data.Drug{}
		drugsMany := []interface{}{}
		for _, j := range drugs {
			j.Name = strings.ToLower(j.Name)
			j.UploaderID= upId
			j.TimeUploaded = time.Now()
			j.CompanyName = cname
			j.CompPhone = cphone
			j.Id = primitive.NewObjectID()
			j.ExpiryDate = j.TimeUploaded.AddDate(0, j.ExpiryMonth, 0)
			drugsMany = append(drugsMany, j)
			

		}

		if err != nil {
			log.Print(err)
		}

		
		drugsColl := base.Collection("drugs")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		result, err := drugsColl.InsertMany(ctx, drugsMany)
		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Print(err)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusOK)
		jsonRes, err := json.Marshal(result.InsertedIDs)
		w.Write([]byte(jsonRes))

	}
}

//Update updates a drug
func Update (base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		idHex := vars["id"]
		id, _ := primitive.ObjectIDFromHex(idHex)
		var drug data.Drug
		err := json.NewDecoder(r.Body).Decode((&drug))
			if err != nil {
			ErrorWithJSON(w, Message{"Bad request body"}, http.StatusBadRequest)
			return
		}
		drug.Name = strings.ToLower(drug.Name)
		coll := base.Collection("drugs")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		result, err := coll.ReplaceOne(
			ctx,
			bson.M{"_id": id},
			bson.D{
				{"$set", drug},
			},
		 )
		if err != nil {
			switch err {
			default:
				ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
				log.Println("Failed to delete drug: ", err)
				return
		}

		//w.WriteHeader(http.StatusOK)
		jsonRes, _ := json.Marshal(result.UpsertedCount)
		ResponseWithJSON(w, jsonRes, http.StatusOK)
	}
	
	}
}
func SendMyDrugs(base *mongo.Database) func (w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		uId := vars["uid"]
		id, err := primitive.ObjectIDFromHex(uId)
		opts := options.Find()
		opts.SetSort(bson.D{{"time_uploaded", -1}})
		if err != nil {
			ErrorWithJSON(w, Message{"You have no valid Id"}, http.StatusBadRequest)
		}
		filter := bson.D{{
			"uploader_id", bson.D{{"$eq", id}},
		}}
		drugColl := base.Collection("drugs")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		cursor, err := drugColl.Find(ctx, filter, opts)

		if err != nil {
			ErrorWithJSON(w, Message{"You have no drug in the database"}, http.StatusBadRequest)
		}
		defer cursor.Close(ctx)
		drugs := []data.Drug{}
		for cursor.Next(ctx) {
			var drug data.Drug
			err := cursor.Decode(&drug)
			if err != nil {
				log.Print(err)
			}
			drug.Name = strings.ToTitle(drug.Name)
			drugs = append(drugs, drug)
		}
		resp, err := json.Marshal(drugs)
		ResponseWithJSON(w, resp, http.StatusOK)
	}
}

func UpdateMyDrugs(base *mongo.Database) func (w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		idStr := mux.Vars(r)["uid"]
		uId, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			ErrorWithJSON(w, Message{"Wrong id"}, http.StatusBadRequest)
		}
		var drugs  []data.Drug
		
		 if err = json.NewDecoder(r.Body).Decode(&drugs); err != nil {
			 log.Print(err)
			 ErrorWithJSON(w, Message{"Bad drug fields"}, http.StatusBadRequest)
		 }
		 coll := base.Collection("drugs")
		 var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		 for _, drug := range drugs {
			 //dIdStr := drug.Id
			 //dId, _ := primitive.ObjectIDFromHex(dIdStr)
			 drug.TimeUpdated  =time.Now()
			 drug.UploaderID = uId
			 drug.Name = strings.ToLower(drug.Name)
			 
			_, err := coll.UpdateOne(
				ctx,
				bson.M{"_id": drug.Id},
				bson.D{
					{"$set", drug }},

			)
			 if err != nil {
				 log.Print("error updating this one")
			 }
			
		 }
		 resp, _ := json.Marshal("done")
		 ResponseWithJSON(w, resp, http.StatusOK)

	}
}

//Search returns a list of drugs tet have a particular name
func Search(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}

		w.Header().Add("Content-Type", "application/json")
		opt := options.Find()
		opt.SetSort(bson.D{{"time_uploaded", -1}})

		vars := mux.Vars(r)
		name := vars["name"]
		name = strings.ToLower(name)
		coll := base.Collection("drugs")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		f := bson.D{{"$and", []bson.D{
			bson.D{{"name", bson.D{{"$eq", name}}}},
			bson.D{{"quantity_available", bson.D{{"$gt", 0}}}},
		}}}
		// filterCursor, err := coll.Find(ctx, bson.D{{"name",bson.D{{"$eq", name}}}}, opt)
		filterCursor, err := coll.Find(ctx, f, opt)
		drugs := []data.Drug {}
		if err != nil {
			ErrorWithJSON(w, Message{"no result"}, http.StatusNotFound)
			return
		}
		defer filterCursor.Close(ctx)
		for filterCursor.Next(ctx) {
			var drug data.Drug
			err := filterCursor.Decode(&drug)
			if err != nil {
				log.Print(err)
			}
			drug.Name = strings.ToTitle(drug.Name)
			drugs = append(drugs, drug)
		}
		
		// For paged data
		// var limit int64 = 10
		// var page int64 = 1
		// projection := bson.D {
		// 	{"name", 1},
		// 	{"description", 1},
		// 	{"brand", 1},
		// 	{"expiry_date", 1},
		// 	{"batcn_no", 1}, 
		// 	{"quantity_available", 1},
		// 	{"product_image", 1},
		// 	{"price", 1},
		// 	{"location", 1},
		// 	{"time_uploaded", 1},
		// }
		// //projection = bson.D{{}}
		// match := bson.D{{"name", bson.D{{"$eq", name}}}}
		// pagedData, err := New(coll).Limit(limit).Page(page).Sort("name", 1).Select(projection).Filter(match).Find()
		// if err != nil {
		// 	ErrorWithJSON(w, Message{"No item in collection"}, http.StatusBadRequest)
		// }
		// var pagedDrugs []data.Drug
		// for _, raw := range pagedData.Data {
		// 	var drug data.Drug
		// 	if marshallErr := bson.Unmarshal(raw, drug); marshallErr == nil {
		// 		pagedDrugs = append(pagedDrugs, drug)
		// 	}
		// }
		

		if err != nil {
			ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
			log.Println("Failed to find drugs with this name")
			return
		}

		resp, er := json.Marshal(drugs)
		if er != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, resp, http.StatusOK)
		
	}	
	
}

//ReturnThisDrug returns a drug with given id
func ReturnThisDrug(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, _ := primitive.ObjectIDFromHex(idStr)

		var drug data.Drug
		coll := base.Collection("drugs")
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&drug)
		if err != nil {
			ErrorWithJSON(w, Message{"error while looking for id"}, http.StatusInternalServerError)
			log.Println(err)
			return
		}
		respBody, err := json.Marshal(&drug)
		ResponseWithJSON(w, respBody, http.StatusOK)
	}

}

//DeleteDrug deletes a drug from the database
func DeleteDrug(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		w.Header().Add("Content-Type", "application/json")
		vars := mux.Vars(r)
		idStr := vars["id"]
		dIdHex := vars["uid"]
		collSen := base.Collection("seniors")
		dId, _ := primitive.ObjectIDFromHex(dIdHex)
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		var user data.Senior
		err := collSen.FindOne(ctx, bson.M{"_id": dId}).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		id, _ := primitive.ObjectIDFromHex(idStr)
		coll := base.Collection("drugs")
		res, err := coll.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
				ErrorWithJSON(w, Message{"Database error"}, http.StatusInternalServerError)
				log.Println("Failed to delete drug: ", err)
				return
			}

		resp := res.DeletedCount
		//w.WriteHeader(http.StatusOK)
		jResp, _ := json.Marshal(resp)
		w.Write(jResp)
	}

}
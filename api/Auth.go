package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"os"
	"github.com/OLUWAMUYIWA/Adel/data"
	jwt "github.com/dgrijalva/jwt-go"
	//"github.com/dgrijalva/jwt-go/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)




type Token struct {
	Token string `json:"accessToken"`
	UserID string	`json:"id"`
	Role	string	`json:"role"`
	CompanyName	string `json:"company_name"`
	Name	string		`json:"name"`
	ConpanyPhone	string	`json:"cphone"`
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var sec = os.Getenv("secret")

var Signature = []byte(sec)

// var Signature = []byte("984fv873rfnvfo9u34rb34340b5geor08343otf89wfbw4893")

// LoginHandler handles the login of usersm returnin the token
func LoginHandler(base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		if (*r).Method == "OPTIONS" {
			return
		}
		var ctx, _  = context.WithTimeout(context.Background(), 10 * time.Second)
		log.Printf("Route: %v", r.URL)
		w.Header().Add("Content-Type", "application/json")
		var user data.UserLogin
		err := json.NewDecoder(r.Body).Decode(&user)
		email := strings.ToLower(user.Email)
		password := strings.ToLower(user.Password)
		collSen := base.Collection("seniors")
		collJun := base.Collection("juniors")
		collBoss := base.Collection("bosses")
		log.Printf("This is what entered: %s", user)


		filter :=  bson.D{{ "$and", []bson.D {
					bson.D{{"email", email}},
					bson.D{{"password", password}},
					bson.D{{"verified", bson.D{{"$eq", true}}}},
			}}}

		log.Printf("filter: %v", filter)
		
		cursor := collJun.FindOne(ctx,filter)
		log.Print(cursor.Err())
		if cursor.Err() == mongo.ErrNoDocuments {
			cursor = collSen.FindOne(ctx, filter)
			if cursor.Err() == mongo.ErrNoDocuments {
				cursor = collBoss.FindOne(ctx, filter)
			}
		}
		if cursor.Err() != nil {
			log.Printf("error happened: %s", cursor.Err())
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "user not found in DB")
			return
		}
		var thisUser data.Senior
		err = cursor.Decode(&thisUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Error while signing the token")
			
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			
			"username": thisUser.Email,
			"id": thisUser.Id,
			"role": thisUser.Role,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
			"iat": time.Now().Unix(),
		})
		tokenString, err := token.SignedString(Signature)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "Error while signing the token")
			fatal(err)
		}

		response := Token{tokenString, thisUser.Id.Hex(), thisUser.Role, thisUser.CompanyName, thisUser.Name, thisUser.PhoneNo}
		jsonResponse, err := json.Marshal(response)
		log.Print(jsonResponse)
		ResponseWithJSON( w, jsonResponse, http.StatusOK)
		}

}

// AuthorizeWareAll serves a the middleware to both senior and junior
func AuthorizeWareAll(w http.ResponseWriter, r *http.Request,  next http.HandlerFunc) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
			return
		}
	w.Header().Add("Content-Type", "application/json")
	reqT := r.Header.Get("Authorization")
	if reqT == "" {
		w.WriteHeader(http.StatusUnauthorized)
			//fmt.Fprint(w, "Token invalid for all")
			return
	}
	splitToken := strings.Split(reqT, "access_token ")
	reqT =strings.TrimSpace(splitToken[1]) 
	//tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(r)
	token, err := jwt.Parse(reqT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v",
		token.Header["alg"])
		}
		return Signature, nil
	})
	
	if err == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && (claims["role"] == "junior" || claims["role"] == "senior" || claims["role"] == "boss") {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token invalid for all")
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized")
	}
}

//AuthorizeWareJunior authorizes seniors
func AuthorizeWareJunior(w http.ResponseWriter, r *http.Request,  next http.HandlerFunc) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
			return
		}
	w.Header().Add("Content-Type", "application/json")
	reqT := r.Header.Get("Authorization")
	if reqT == "" {
		w.WriteHeader(http.StatusUnauthorized)
			//fmt.Fprint(w, "Token invalid for all")
			return
	}
	splitToken := strings.Split(reqT, "access_token ")
	reqT =strings.TrimSpace(splitToken[1])
	//tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(r)
		token, err := jwt.Parse(reqT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v",
			token.Header["alg"])
			}
			return Signature, nil
		})
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && (claims["role"] == "junior" || claims["role"] == "senior" || claims["role"] == "boss") {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Token invalid for jun")
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
}

//AuthorizeWareSenior authorizes seniors
func AuthorizeWareSenior(w http.ResponseWriter, r *http.Request,  next http.HandlerFunc) {
	w.Header().Add("Content-Type", "application/json")
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
			return
		}
	reqT := r.Header.Get("Authorization")
	splitToken := strings.Split(reqT, "access_token ")
	reqT =strings.TrimSpace(splitToken[1])
	//tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(r)
		token, err := jwt.Parse(reqT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v",
			token.Header["alg"])
			}
			return Signature, nil
		})
		// checkClaim, _ := token.Claims.(jwt.MapClaims)
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && (claims["role"] == "senior" || claims["role"] == "boss") {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Token invalid for sen")
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
}

//AuthorizeWareBoss deals with the authorization of the boss
func AuthorizeWareBoss(w http.ResponseWriter, r *http.Request,  next http.HandlerFunc) {
	setupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
			return
		}
	w.Header().Add("Content-Type", "application/json")
	reqT := r.Header.Get("Authorization")
	if reqT == "" {
		w.WriteHeader(http.StatusUnauthorized)
			return
	}
	splitToken := strings.Split(reqT, "access_token ")
	reqT =strings.TrimSpace(splitToken[1])
	//tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(r)
		token, err := jwt.Parse(reqT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v",
			token.Header["alg"])
			}
			return Signature, nil
		})
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["role"] == "boss" {
				next(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Token invalid for boss")
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Unauthorized")
			return
		}
}
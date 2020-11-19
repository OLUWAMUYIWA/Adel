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
	jwt "github.com/dgrijalva/jwt-go"
	//"github.com/dgrijalva/jwt-go/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)




type Token struct {
	Token string `json:"accessToken"`
	UserID string	`json:"id"`
	Role	string	`json:"role"`
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var Signature = []byte("984fv873rfnvfo9u34rb34340b5geor08343otf89wfbw4893")

// LoginHandler handles the login of usersm returnin the token
func LoginHandler(ctx context.Context, base *mongo.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
			log.Print("got here")
			cursor = collSen.FindOne(ctx, filter)
			if cursor.Err() == mongo.ErrNoDocuments {
				log.Print("here after")
				cursor = collBoss.FindOne(ctx, filter)
			}
		}
		if cursor.Err() != nil {
			log.Printf("error happened: %s", cursor.Err())
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "user not found in DB")
			log.Print("here")
			return
		}
		var thisUser data.Senior
		err = cursor.Decode(&thisUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			log.Printf("this is the error", err)
			fmt.Fprintln(w, "Error while signing the token")
			
		}
		//log.Print(thisUser.Name)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			
			"username": thisUser.Email,
			"id": thisUser.Id,
			"role": thisUser.Role,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
			"iat": time.Now().Unix(),
		})
		log.Printf("token: %s", token.Claims)
		tokenString, err := token.SignedString(Signature)
		log.Println(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//log.Printf("this is the error", err)
			fmt.Fprintln(w, "Error while signing the token")
			fatal(err)
		}

		response := Token{tokenString, thisUser.Id.Hex(), thisUser.Role}
		jsonResponse, err := json.Marshal(response)
		ResponseWithJSON( w, jsonResponse, http.StatusOK)
		}

}

// AuthorizeWareAll serves a the middleware to both senior and junior
func AuthorizeWareAll(w http.ResponseWriter, r *http.Request,  next http.HandlerFunc) {
	w.Header().Add("Content-Type", "application/json")
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
	w.Header().Add("Content-Type", "application/json")
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
		log.Printf("token Valid? %v", token.Valid )
		checkClaim, _ := token.Claims.(jwt.MapClaims)
		log.Printf("role: %s", checkClaim["role"])
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && (claims["role"] == "senior" || claims["role"] == "boss") {
				log.Print("still in")
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
	w.Header().Add("Content-Type", "application/json")
	reqT := r.Header.Get("Authorization")
	splitToken := strings.Split(reqT, "access_token ")
	reqT =strings.TrimSpace(splitToken[1])
	//tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(r)
		token, err := jwt.Parse(reqT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v",
			token.Header["alg"])
			log.Print("Boss 1")
			}
			return Signature, nil
		})
		if err == nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["role"] == "boss" {
				log.Print("boss 2")
				next(w, r)
			} else {
				log.Print("boss no")
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
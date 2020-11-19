package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User contains properties common to all users
type User struct {
	Id				primitive.ObjectID	`json:"id" bson: "_id,omitemmpty"`
	CompanyName		string			`json: "company_name" bson: "company_name,omitempty"`
	Email 			string			`json: "email" bson: "email,omitempty"`
	Address			string			`json: "address" bson: "address,omitempty"`
	//UserName		string			`json:"UserName" bson: "user_name,omitempty"`
	Password		string			`json: "password" bson: "password,omitempty"`
	Verified 		bool			`json: "verified" bson: "verified,omitempty"`
}

//Junior refers to the user who can only search but cannot post
type Junior struct {
	Id				primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	CompanyName		string			`json:"company_name" bson:"company_name,omitempty"`
	Email 			string			`json:"email" bson:"email,omitempty"`
	Address			string			`json:"address" bson:"address,omitempty"`
	Phone			string			`json:"phone" bson:"phone"`
	Name			string			`json:"name" bson:"name,omitempty"`
	Password		string			`json:"password" bson:"password,omitempty"`
	Verified 		bool			`json:"verified" bson:"verified"`
	PracticeArea string 		`json:"practice_area" bson:"practice_area,omitempty"`
	Role		string			`json:"role" bson: "role,omitempty"`
	TimeRegd	time.Time		`json:"time_regd" bson:"time_regd,omitempty"`
}

//Senior is the second most-privileged user. He can upload/post drug items
type Senior struct {
	Id				primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	CompanyName		string			`json:"company_name" bson:"company_name,omitempty"`
	Email 			string			`json:"email" bson:"email,omitempty"`
	Address			string			`json:"address" bson: "address,omitempty"`
	Phone			string			`json:"phone" bson:"phone"`
	Name		string				`json:"name" bson:"name,omitempty"`
	Password		string			`json:"password" bson:"password,omitempty"`
	Verified 		bool			`json:"verified" bson:"verified"`
	PracticeArea string 			`json:"practice_area" bson:"practice_area,omitempty"`
	Role		string				`json "role" bson: "role,omitempty"`
	SuperintendentPharmName 		string	`json:"superintendentPharmName" bson:"pharm_name,omitempty"`
	SuperintendentPharmRegNo		string	`json:"superintendentPharmRegNo" bson:"pharm_reg,omitempty"`
	SuperintendentPharmLicenceNo	string	`json:"superintendentPharmLicenceNo" bson:"pharm_licence,omitempty"`
	TimeRegd	time.Time			`json:"time_regd" bson:"time_regd,omitempty"`

}

type Boss struct {
	Id 		primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	Name	string		`json:"name" bson:"name"`
	Email	string	`json:"email" bson:"email"`
	Password string  `json:"password" bson:"password"`
	Role		string		`json: "role" bson: "role,omitempty"`
	Verified	bool 	`json:"verified" bson:"verified,omitempty"`
}

type UserLogin struct {
	Email		string				`json:"email"`
	Password		string			`json: "password"`
}
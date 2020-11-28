package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Drug struct {
	Id					primitive.ObjectID	`json:"id" bson:"_id,omitempty"`
	Name 				string			`json:"name" bson:"name,omitempty"`
	Description			string			`json:"description" bson:"description,omitempty"`
	Brand				string			`json:"brand" bson:"brand,omitempty"`
	ExpiryMonth			int				`json:"exp,omitempty" bson:"exp,omitempty"`
	BatcnNo				string			`json:"batch_no" bson:"batch_no,omitempty"`
	QuantityAvailable	int				`json:"quantity_available" bson:"quantity_available,omitempty"`
	Price 				float64			`json:"price" bson:"price,omitempty"`
	Location 			string			`json:"location" bson:"location,omitempty"`
	CompPhone			string			`json:"cphone" bson:"cphone,omitempty"`
	ProductImage		string			`json:"product_image" bson:"product_image,omitempty"`
	UploaderID			primitive.ObjectID	`json:"uploader_id" bson:"uploader_id"`
	TimeUploaded		time.Time		`json:"time_uploaded" bson:"time_uploaded,omitempty"`
	TimeUpdated			time.Time		`json:"time_updated" bson:"time_updated,omitempty"`
	CompanyName			string			`json:"company_name" bson:"company_name,omitempty"`
	ExpiryDate			time.Time		`json:"expiry_date" bson:"expiry_date"`
}
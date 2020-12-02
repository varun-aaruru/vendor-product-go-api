package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Vendor Struct
type Vendor struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	VendorName  string             `json:"vendorname" bson:"vendorname"`
	Designation string             `json:"designation" bson:"designation"`
	Age         int                `json:"age,omitempty" bson:"age,omitempty"`
	Products    []string           `json:"product_names,omitempty" bson:"product_names,omitempty"`
}

//Product Struct
type Product struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Price  float64            `json:"price" bson:"price"`
	Vendor string             `json:"vendor_name,omitempty" bson:"vendor_name,omitempty"`
}

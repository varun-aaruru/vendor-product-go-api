package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gitjub.com/user/vendorapi/helper"
	"gitjub.com/user/vendorapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with helper class
var vendorcollection = helper.ConnectDB().Collection("vendors")
var productcollection = helper.ConnectDB().Collection("products")

func main() {
	//Init Router
	r := mux.NewRouter()

	// routes
	r.HandleFunc("/api/vendors", getVendors).Methods("GET")
	r.HandleFunc("/api/vendors/{id}", getVendor).Methods("GET")
	r.HandleFunc("/api/vendors", createVendor).Methods("POST")
	r.HandleFunc("/api/vendors/{id}", updateVendor).Methods("PUT")
	r.HandleFunc("/api/vendors/{id}", deleteVendor).Methods("DELETE")
	r.HandleFunc("/api/products", getProducts).Methods("GET")
	r.HandleFunc("/api/products/{id}", getProduct).Methods("GET")
	r.HandleFunc("/api/products", createProduct).Methods("POST")
	r.HandleFunc("/api/products/{id}", updateProduct).Methods("PUT")
	r.HandleFunc("/api/products/{id}", deleteProduct).Methods("DELETE")
	// set our port address
	r.Use(mux.CORSMethodMiddleware(r))
	log.Fatal(http.ListenAndServe(":8080", r))

}
func getVendors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var vendors []models.Vendor
	cur, err := vendorcollection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var vendor models.Vendor
		// & character returns the memory address of the following variable.
		err := cur.Decode(&vendor) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		vendors = append(vendors, vendor)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(vendors)
}
func getVendor(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var vendor models.Vendor
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	err := vendorcollection.FindOne(context.TODO(), filter).Decode(&vendor)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(vendor)
}
func createVendor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var vendor models.Vendor

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&vendor)

	result, err := vendorcollection.InsertOne(context.TODO(), vendor)
	fmt.Println(vendor.Products)

	if err != nil {
		helper.GetError(err, w)
		return
	}
	newResult, newErr := productcollection.UpdateMany(
		context.TODO(),
		bson.M{"name": bson.M{"$in": vendor.Products}},
		bson.D{
			{"$set", bson.D{{"vendor_name", vendor.VendorName}}},
		},
	)
	if newErr != nil {
		log.Fatal(newErr)
	}
	fmt.Printf("Updated %v Documents!\n", newResult.ModifiedCount)
	json.NewEncoder(w).Encode(result)
}
func updateVendor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var vendor models.Vendor

	// filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&vendor)

	//update model.
	fmt.Println(vendor.Products)
	update := bson.D{
		{"$set", bson.D{
			{"vendorname", vendor.VendorName},
			{"designation", vendor.Designation},
			{"age", vendor.Age},
			{"product_names", vendor.Products},
		}},
	}
	fmt.Println(vendor.VendorName)
	venName := vendor.VendorName
	err := vendorcollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&vendor)

	if err != nil {
		helper.GetError(err, w)
		return
	}
	fmt.Println(vendor.VendorName)
	newResult, newErr := productcollection.UpdateMany(
		context.TODO(),
		bson.M{"name": bson.M{"$in": vendor.Products}},
		bson.D{
			{"$set", bson.D{{"vendor_name", venName}}},
		},
	)
	if newErr != nil {
		log.Fatal(newErr)
	}
	fmt.Printf("Updated %v Documents!\n", newResult.ModifiedCount)
	vendor.ID = id
	fmt.Println(vendor.VendorName)
	json.NewEncoder(w).Encode(vendor)
}
func deleteVendor(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := vendorcollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// we created Book array
	var products []models.Product

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := productcollection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var product models.Product
		// & character returns the memory address of the following variable.
		err := cur.Decode(&product) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		products = append(products, product)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(products) // encode similar to serialize process.
}
func getProduct(w http.ResponseWriter, r *http.Request) {
	// set header.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var product models.Product
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := productcollection.FindOne(context.TODO(), filter).Decode(&product)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(product)
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var product models.Product

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&product)

	// insert our book model.
	result, err := productcollection.InsertOne(context.TODO(), product)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	var productNames []string
	fmt.Println(product.Name)
	productNames = append(productNames, product.Name)
	fmt.Println(productNames)
	newResult, newErr := vendorcollection.UpdateMany(
		context.TODO(),
		bson.M{"vendorname": product.Vendor},
		bson.M{
			"$set": bson.D{{"product_names", productNames}},
		},
	)
	fmt.Println(newResult)
	fmt.Println(newErr)
	json.NewEncoder(w).Encode(result)
}
func updateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var product models.Product

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&product)
	fmt.Println(product.Vendor)
	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", product.Name},
			{"price", product.Price},
			{"vendor_name", product.Vendor},
		}},
	}
	fmt.Println(update)

	err := productcollection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)

	if err != nil {
		helper.GetError(err, w)
		return
	}
	fmt.Println(product.Vendor)

	var productNames []string
	productNames = append(productNames, product.Name)
	fmt.Println(product.Vendor)
	newResult, newErr := vendorcollection.UpdateMany(
		context.TODO(),
		bson.M{"vendorname": product.Vendor},
		bson.M{
			"$set": bson.D{{"product_names", productNames}},
		},
	)
	fmt.Println(newResult)
	fmt.Println(newErr)

	product.ID = id

	json.NewEncoder(w).Encode(product)
}
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := productcollection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

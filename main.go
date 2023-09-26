package main

import (
	"encoding/json"
	"fmt"
)

// Golang has own data structure called struct
// Golang does not natively understand json data structure
// string, .Number are primitive data types
// 'Address' is a custom data type

// Address struct
type Address struct {
	City string
	State string
	Country string
	Zipcode json.Number
}

// user struct
type User struct {
	Name string
	Age json.Number
	Contact string
	Company string
	Address Address
}


func main() {
	dir := "./" //golang to create a directory with collections data

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}
	// a slice of 'employees' of type User
	employees := []User{
		{"John", "23", "071234356", "Freelance Ltd", Address{"bangalore", "karnataka", "india", "098765"}},
		{"Jane", "28", "075468012", "My Code Ltd", Address{"kolkata", "karnataka", "india", "097001"}},
		{"Sara", "25", "071357901", "WebDev Co", Address{"chennai", "tamil nadu", "india", "096165"}},
		{"Amrit", "29", "072468012", "CDCI Ltd", Address{"panaji", "goa", "india", "095765"}},
		{"Divya", "27", "0732343569", "Backend Co", Address{"jaipur", "rajasthan", "india", "094265"}},
		{"Laksmi", "26", "074534351", "Frontend Ltd", Address{"hyderabad", "telangana", "india", "093165"}},
		{"Akhil", "31", "075214356", "Dominate", Address{"bhopal", "madhya pradesh", "india", "092265"}},
	
	}
	// value.Name gives the name of the user file
	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name: value.Name,
			Age: value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}
}
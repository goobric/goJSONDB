package main

import "encoding/json"

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

}
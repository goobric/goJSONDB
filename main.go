package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber" //logger package
)

const Version = "1.0.0"

type (
	Logger interface{
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debud(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct{
		mutex sync.Mutex
		mutexes map[string]*sync.Mutex
		dir string
		log Logger

	}
)

type Options struct {
	Logger
}

// creates db and returns a Driver
func New(dir string, options *Options)(*Driver, error){
	dir = filepath.Clean(dir)

	opts := Options{}
	if options != nil {
		opts = *options
	}
	if opts.Logger != nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir: dir,
		mutexes: make(map[string]*sync.Mutex),
		log: opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (db already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating the db at '%s'... ", dir)
	return &driver, os.MkdirAll(dir, 0755)
}
//struct methods
func (d *Driver) Write(collection, resource string, v interface{}) error{
	if collection == ""{
		return fmt.Errorf("Missing collection - no place to store the record!")
	}
	if resource == ""{
		return fmt.Errorf("Missing resource - unable to save record - no name!")
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resource+".json")
	tmpPath := fnlPath + ".tmp"

	if err := os.MkdirAll(dir, 0755); err != nil{
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil{
		return err
	}
	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, fnlPath)
}

func (d *Driver) Read(collection, resource string, v interface{}) error{
	if collection == ""{
		return fmt.Errorf("missing collection - no place to save record")
	}
	if resource == ""{
		return fmt.Errorf("missing resource - no place to save record - no name")
	}
	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}
	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &v)
}

func (d *Driver) ReadAll(collection string)([]string, error){
	if collection == ""{
		return nil, fmt.Errorf("missing collection unable to read")
	}
	
	if _, err := stat(dir); err != nil{
		return nil, err
	}
	files, _ := ioutil.ReadDir(dir)

	var records []string

	for _, file := range files{
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil{
			return nil, err
		}
	}

}

func (d *Driver) Delete() error{

}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex{
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[collection]

	if !ok{
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
}

func stat(path string)(fi os.FileInfo, err error){
	if fi, err = os.Stat(path); os.IsNotExist(err){
		fi, err = os.Stat(path + ".json")
	}
	return
}

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
	// value.Name is the resource
	// the elements are the values
	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name: value.Name,
			Age: value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error user", err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _, record := range records {
		empolyeeFound := User{}
		if err := json.Unmarsh([]byte(record), &empolyeeFound); err != nil {
			fmt.Println("Error employee", err)
		}
		allusers = append(allusers, empolyeeFound)
	}
	fmt.Println((allusers))

	// if err := db.Delete("users", "john"); err != nil {
	// 	fmt.Println("Error", err)
	// }
	// if err := db.Delete("user", ""); err != nil {
	// 	fmt.Println("Error", err)
	// }
}
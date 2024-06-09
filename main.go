package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type ( //this is so that we dont have to type "type" everytime
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{}) // Fixed method name to start with a capital letter
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex // Fixed field name from "mutextes" to "mutexes"
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("using '%s' (Database already Exists)\n", dir)
		return &driver, nil
	}
	opts.Logger.Debug("Creating the database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

func (d *Driver) Write(string, string, User) error {
	// Implementation needed
	return nil
}

func (d *Driver) Read(string, string) (string, error) {
	// Implementation needed
	return "", nil
}

func (d *Driver) ReadAll(string) (map[string]string, error) {
	// Implementation needed
	return nil, nil
}

func (d *Driver) Delete(string, string) error {
	// Implementation needed
	return nil
}

func (d *Driver) getOrCreateMutex(name string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, ok := d.mutexes[name]
	if !ok {
		m = &sync.Mutex{}
		d.mutexes[name] = m
	}
	return m
}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
		{"Alice", "30", "1002003001", "Software Engineer", Address{"New York", "NY", "USA", "10001"}},
		{"Bob", "25", "1002003002", "Data Scientist", Address{"San Francisco", "CA", "USA", "94105"}},
		{"Charlie", "28", "1002003003", "Product Manager", Address{"Austin", "TX", "USA", "73301"}},
		{"David", "35", "1002003004", "DevOps Engineer", Address{"Seattle", "WA", "USA", "98101"}},
		{"Eve", "32", "1002003005", "UX Designer", Address{"Chicago", "IL", "USA", "60601"}},
	}

	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}

	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(records) //these are in json data type we have to further process it to use it

	allusers := []User{}

	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil { //t's a common practice to handle errors immediately after they occur rather than letting them propagate through the program unchecked
			fmt.Println("Error", err)
		}
		allusers = append(allusers, employeeFound)
	}
	fmt.Println((allusers))

	// if err := db.Delete("user","Alice");err!=nil{
	// fmt.Println("Error", err)
	// }

	// if err := db.Delete("user","");err!=nil{  //this is for delete all operation where we pass an empty string
	// 	fmt.Println("Error", err)
	// 	}
}

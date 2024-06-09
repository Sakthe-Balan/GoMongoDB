package main

import (
	"encoding/json"
	"fmt"
)

const Version = "1.0.1"

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
		{"Frank", "27", "1002003006", "QA Engineer", Address{"Denver", "CO", "USA", "80201"}},
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
		fmt.Println("Erros", err)
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

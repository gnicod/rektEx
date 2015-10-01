package main

import (
	"encoding/json"
	"fmt"

	r "github.com/dancannon/gorethink"
)

var session *r.Session

// Struct tags are used to map struct fields to fields in the database
type ErrorException struct {
	Id      string `gorethink:"id,omitempty"`
	Name    string `gorethink:"name"`
	Message string `gorethink:"message"`
}

func init() {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "logs",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	id := insertRecord()
	printStr(id)

	fetchAllRecords()

}

func insertRecord() string {
	var data = map[string]interface{}{
		"Name":    "Dashboard",
		"Message": "Some usefull message",
	}

	result, err := r.Table("exceptions").Insert(data).RunWrite(session)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	printStr("*** Insert result: ***")
	printObj(result)
	printStr("\n")

	return result.GeneratedKeys[0]
}

func fetchAllRecords() {
	rows, err := r.Table("exceptions").Run(session)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read records into persons slice
	var errorsExceptions []ErrorException
	err2 := rows.All(&errorsExceptions)
	if err2 != nil {
		fmt.Println(err2)
		return
	}

	printStr("*** Fetch all rows: ***")
	for _, p := range errorsExceptions {
		printObj(p)
	}
	printStr("\n")
}

func printStr(v string) {
	fmt.Println(v)
}

func printObj(v interface{}) {
	vBytes, _ := json.Marshal(v)
	fmt.Println(string(vBytes))
}

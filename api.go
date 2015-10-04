package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	rethink "github.com/dancannon/gorethink"
	"github.com/go-martini/martini"
	"net/http"
	"strings"
)

var session *rethink.Session

// Struct tags are used to map struct fields to fields in the database
type Log struct {
	Id      string `gorethink:"id,omitempty"`
	AppName string `gorethink:"AppName"`
	Message string `gorethink:"Message"`
	Ip      string
}

// TODO sucks, need to use the same session than the one declared in main.go
func init() {
	var err error
	session, err = rethink.Connect(rethink.ConnectOpts{
		Address:  "localhost:28015",
		Database: "logs",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func onChange() {
	res, err := rethink.Table("exceptions").Changes().Run(session)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Use goroutine to wait for changes. Prints the first 10 results
	go func() {
		var response interface{}
		for res.Next(&response) {
			fmt.Println("okkkkk")
		}

		if res.Err() != nil {
			fmt.Println(res.Err())
		}

	}()
}

// TODO separate view from *model*
func NewLog(log Log, req *http.Request, args martini.Params, r render.Render) {
	ip := strings.Split(req.RemoteAddr, ":")[0]
	log.Ip = ip
	result, err := rethink.Table("exceptions").Insert(log).RunWrite(session)
	if err != nil {
		fmt.Println(err)
	}
	key := result.GeneratedKeys[0]
	b, err := json.Marshal(log)
	broadcastMessage(1, b)
	r.JSON(200, map[string]interface{}{"key": key})
}

// Get all log with appname passed in params
func GetLogForApp(appname string) (logs []Log, err error) {
	rows, err := rethink.Table("exceptions").Filter(map[string]interface{}{
		"AppName": appname,
	}).Limit(10).Run(session)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err2 := rows.All(&logs)
	if err2 != nil {
		return nil, err
	}
	if logs == nil {
		return nil, err
	}
	return logs, nil
}

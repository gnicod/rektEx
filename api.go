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
func GetLogForApp(args martini.Params, r render.Render) {
	rows, err := rethink.Table("exceptions").Filter(map[string]interface{}{
		"AppName": args["appname"],
	}).Run(session)
	if err != nil {
		fmt.Println(err)
		r.JSON(200, map[string]interface{}{})
	}

	var logs []Log
	err2 := rows.All(&logs)
	if err2 != nil {
		r.JSON(200, map[string]interface{}{})
		return
	}
	if logs == nil {
		r.JSON(200, map[string]interface{}{})
		return
	}
	r.JSON(200, logs)
}

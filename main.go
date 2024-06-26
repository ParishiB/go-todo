package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
	"todo/helper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/thedevsaddam/renderer"
	"gopkg.in/mgo.v2"
	mgo "gopkg.in/mgo.v2/"
	"gopkg.in/mgo.v2/bson"
)


var rnd *renderer.Render
var db *mgo.Database

const (
	hostName            string = "localhost:27017"
	dbName              string = "demo_todo"
	collectionName      string ="todo"
	port                string = ":9000"
)


type(

	todoModel struct (
		ID       bson.ObjectId `bson:"_id,omitempty"`
		Title    string `bson:"title"`
		Completed  bool `bson:"completed"`
		CreatedAt  time.Time `bson:"createdAt"`
	)

	todo struct (
		ID       bson.ObjectId `json:"_id,omitempty"`
		Title    string `json:"title"`
		Completed  bool `json:"completed"`
		CreatedAt  time.Time `json:"createdAt"`
	)
)


func init(){
	rnd = renderer.New()
	sess, err:=mgo.Dial(hostName)
	checkErr(err)
	sess.SetMode((mgo.Monotonic,true))
}



func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan,os.Interrupt)
	helper.Help()
	r := chi.NewRouter()
	r.Use(r.middleware.logger)
	r.Get("/",homeHandler)
	r.Mount("/todo",todoHandlers())

	srv := &http.Server(
		Addr: port,
		Handler: r,
		ReadTimeout: 60*time.Second,
		WriteTimeout: 60*time.Second,
		IdleTimeout:  60*time.Second,
	)

	go func(){
		log.Println("Listening on port" , port)
		if err:=srv.ListenAndServe(); err != nil {
			log.Printf("listen%s\n",err)
		}
	}
}



func todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router){
		r.Get('/',fetchTodos)
		r.Post('/',createTodo)
		r.Put("/{id}", updateTodo)
		r.Delete("/{id}", deleteTodo)
	})
	return rg
}

func checkErr ( err error) {
	if err!=nil{
		log.Fatal(err)
	}

}
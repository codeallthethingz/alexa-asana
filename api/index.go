package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codeallthethingz/alexa-asana/api/asana"
	"github.com/codeallthethingz/alexa-asana/api/common"
	"github.com/gorilla/mux"
)

var router *mux.Router

type model struct {
	Request request `json:"request"`
	Session session `json:"session"`
}

type session struct {
	User user `json:"user"`
}

type user struct {
	UserID string `json:"userId"`
}

type request struct {
	Intent intent `json:"intent"`
}

type intent struct {
	Slots slots `json:"slots"`
}

type slots struct {
	Utterance utterance `json:"utterance"`
}

type utterance struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	ConfirmationStatus string `json:"confirmationStatus"`
	Source             string `json:"source"`
}

// Handler - check routing and call correct methods
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL.Path)
	if router == nil {
		r, err := CreateRouter()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		router = r
	}
	router.ServeHTTP(w, r)
}

// CreateRouter public so we can test it.
func CreateRouter() (*mux.Router, error) {
	router := mux.NewRouter()
	router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		defer body.Close()
		b, err := ioutil.ReadAll(body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(b))
		model := &model{}
		if err := json.Unmarshal(b, model); err != nil {
			fmt.Println(err.Error())
			return
		}
		c, _ := json.Marshal(model)

		fmt.Println(string(c))
		task := model.Request.Intent.Slots.Utterance.Value
		user := model.Session.User.UserID
		dao := asana.NewDAO(common.NewClient(http.DefaultClient))
		if err := dao.CreateTask(task, user); err != nil {
			task = err.Error()
		}
		w.Write([]byte(`{
			"response": {
			  "outputSpeech": {
				"type": "PlainText",
				"text": "added task to ` + task + `"  
			  },
			  "shouldEndSession": true
			}
		  }`))
	})
	return router, nil
}

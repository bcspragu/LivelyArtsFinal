package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var templates = template.Must(template.ParseGlob("templates/*"))

var wordFreq = make(map[string]int)

type Word struct {
	Word string `json:"text"`
	Size int    `json:"size"`
}

func main() {
	go h.run()

	r := mux.NewRouter()

	r.HandleFunc("/", mainHandler)
	r.HandleFunc("/input", inputHandler)
	r.HandleFunc("/cloud", cloudHandler)
	r.HandleFunc("/ws", serveWs)

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	http.Handle("/", r)

	log.Println("Starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "input.html", struct{}{})
	if err != nil {
		log.Print("Error executing template:", err)
	}
}

func inputHandler(w http.ResponseWriter, r *http.Request) {

	words := strings.Split(r.FormValue("words"), " ")
	for _, word := range words {
		if word != "" {
			if _, ok := wordFreq[word]; ok {
				wordFreq[word]++
			} else {
				wordFreq[word] = 1
			}
		}
	}

	h.broadcast <- []byte(wordsJSON())
}

func cloudHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		WordList string
	}{wordsJSON()}
	err := templates.ExecuteTemplate(w, "cloud.html", data)
	if err != nil {
		log.Print("Error executing template:", err)
	}
}

func wordsJSON() string {
	words := make([]Word, len(wordFreq))
	index := 0
	for word, count := range wordFreq {
		words[index] = Word{word, validCount(count)}
		index++
	}
	enc, err := json.Marshal(words)

	if err != nil {
		log.Print("Error making JSON:", err)
	}

	return string(enc)
}

func validCount(count int) int {
	count = count*5 + 25
	if count > 175 {
		return 175
	} else {
		return count
	}
}

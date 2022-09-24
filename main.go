package main

import (

	"image/png"
	"log"
	"net/http"
	"encoding/json"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"os"
	"io/ioutil"
	"fmt"
	"text/template"

    //"github.com/joho/godotenv"
    "github.com/spf13/viper"
    "github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QrText struct {
	Text string `json:"text"`
}

type Page struct {
	Title string
}



func main() {

	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	//viper.SetConfigType("env")

	//viper.AutomaticEnv()

	port := fmt.Sprint(viper.Get("PORT"))
	

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/qr", qrgenerator).Methods("POST")
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/generator/", viewCodeHandler).Methods("POST")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	

	log.Println("listening on http://localhost:"+port)
	log.Println(http.ListenAndServe(":"+port, loggedRouter))	
}

func qrgenerator(w http.ResponseWriter, r *http.Request) {
	

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Problem reading the body")
	}

	var text QrText

	json.Unmarshal(reqBody, &text)

	
	log.Println(text)

	qrCode, err := qr.Encode(text.Text, qr.L, qr.Auto)
	if err != nil {
		log.Println("Problem encoding the text")
	}

	qrCode, _ = barcode.Scale(qrCode, 128, 128)

	png.Encode(w, qrCode)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "QR Code Generator"}

	t, err := template.ParseFiles("generator.html")
	if err != nil {
		log.Println("Problem parsing html file")
	}

	t.Execute(w, p)
}

func viewCodeHandler(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 128, 128)

	png.Encode(w, qrCode)
}  
package main

import (

	"image/png"
	"net/http"
	"text/template"
	
	"log"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/joho/godotenv"
)

type Page struct {
	Title string
}

func main() {

	err := godotenv.Load(".env")

	  if err != nil {
	    log.Fatalf("Error loading .env file")
	  }

	port := os.Getenv("PORT")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generator/", viewCodeHandler)
	
	log.Println("Listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func homeHandler( w http.ResponseWriter, r * http.Request) {
	p := Page{Title: "QR Code Generator"}

	t, _ := template.ParseFiles("generator.html")

	t.Execute(w, p)
}

func viewCodeHandler(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	png.Encode(w, qrCode)
}

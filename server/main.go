package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"m28909_Hausarbeit_WebprogGo/internal/database"
	"m28909_Hausarbeit_WebprogGo/internal/handlers"
	"m28909_Hausarbeit_WebprogGo/internal/util"
)

var tmpDir = flag.String("tmp", "internal/templates", "Template -Dir.") // Dateipfad für die HTML-Template-Dateien

// main startet den Webserver, stellt die Verbindung zur Datenbank her, lädt die initiale Daten und initialisiert HTTP-Handler.
func main() {
	// Verbindung zu Datenbank herstellen
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Fehler bei der Verbindung zur Datenbank: %v", err)
	}
	defer db.Close() // sicherstellen, dass die Verbindung ordnungsgemäß geschlossen wird

	// Rohdaten aus der zip Datei in die Datenbank laden
	util.UnzipAndLoadDataToDatabase(db)

	flag.Parse()
	// Überprüfen, ob der Pfad von tmpDir existiert
	if _, err := os.Stat(*tmpDir); os.IsNotExist(err) {
		log.Fatalf("Template-Verzeichnis '%s' existiert nicht", *tmpDir)
	}

	// Handler für die Statischen Dateien (CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handler für die Medien (Bilder, Audios)
	http.HandleFunc("/media/", handlers.MakeMediaHandler(db))

	// Handler für die Indexseite
	http.HandleFunc("/", handlers.MakeIndexHandler(db, *tmpDir))
	// Handler für alle weiteren Seiten
	http.HandleFunc("/page/", handlers.MakePageHandler(db, *tmpDir))

	// Server starten
	log.Print(">>> Der Server hört auf :8080 ...")
	log.Print(">>> Der Server wurde gestartet")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

package util

import (
	"archive/zip"
	"io"
	"log"
	"m28909_Hausarbeit_WebprogGo/internal/database"
	"path/filepath"
)

// UnzipAndLoadDataToDatabase entpackt eine ZIP-Datei und speichert die enthaltenen Dateien in der MongoDB-Datenbank.
func UnzipAndLoadDataToDatabase(db *database.MongoDB) {
	// Öffnen der ZIP-Datei
	zipFile, err := zip.OpenReader("data.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	log.Println("\n>>> Rohdateien aus der ZIP Datei lesen ...")
	// Schleife über alle Dateien in der ZIP-Datei
	for _, file := range zipFile.File {
		fileName := filepath.Base(file.Name)
		log.Println(">>> Eingelesene Datei:", fileName) // eingelesenen Dateinamen ausgeben

		// Öffnen der Datei
		rc, err := file.Open()
		if err != nil {
			log.Println("UnzipAndLoadDataToDatabase.Open: ", err)
		}
		defer rc.Close()

		// Lesen der Datei
		data, err := io.ReadAll(rc)
		if err != nil {
			log.Println("UnzipAndLoadDataToDatabase.ReadAll: ", err)
		}

		// Speichern der Datei in der MongoDB
		err = db.SaveFile(fileName, data)
		if err != nil {
			log.Println("UnzipAndLoadDataToDatabase.SaveFile: ", err)
		}
	}
	log.Println(">>> Dateien wurden erfolgreich in Datenbank gespeichert")
}

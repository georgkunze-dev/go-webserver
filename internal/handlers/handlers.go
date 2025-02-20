package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"m28909_Hausarbeit_WebprogGo/internal/database"
	"m28909_Hausarbeit_WebprogGo/internal/structure"
)

// MakeIndexHandler erstellt einen HTTP-Handler, der die Indexseite rendert, indem die
// benötigten Daten der index.md aus der Datenbank geladen und diese in einem Template dargestellt werden.
func MakeIndexHandler(db *database.MongoDB, tmpDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Seite laden
		p, err := structure.LoadPage(db, "index.md")
		if err != nil {
			http.Error(w, "Interner Serverfehler MakeIndexHandler.LoadPage", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		// Seite rendern
		err = renderPage(w, p, tmpDir)
		if err != nil {
			http.Error(w, "Interner Serverfehler MakeIndexHandler.renderPage", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		//Caching verhindern
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
	}
}

// MakePageHandler erstellt einen HTTP-Handler, der dynamische Seiten basierend auf der URL rendert,
// indem die entsprechenden Markdowninhalte aus der Datenbank geladen und diese in einem Template dargestellt werden.
func MakePageHandler(db *database.MongoDB, tmpDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Seite laden
		f := r.URL.Path[len("/page/"):]
		p, err := structure.LoadPage(db, f)
		if err != nil {
			http.Error(w, "Interner Serverfehler MakePageHandler.LoadPage", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		// Seite rendern
		err = renderPage(w, p, tmpDir)
		if err != nil {
			http.Error(w, "Interner Serverfehler MakePageHandler.renderPage", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		//Caching verhindern
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
	}
}

// MakeMediaHandler erstellt einen HTTP-Handler, der Mediendateien aus der Datenbank lädt
// und dem Client mit dem korrekten MIME-Typ zurückgibt.
func MakeMediaHandler(db *database.MongoDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Datei aus Datenbank laden
		filename := r.URL.Path[len("/media/"):]
		data, err := db.GetFile(filename)
		if err != nil {
			http.Error(w, "Media nicht gefunden", http.StatusNotFound)
			return
		}
		// Content-Typ setzen, je nach Dateityp
		mediaType := http.DetectContentType(data)
		w.Header().Set("Content-Type", mediaType)
		//Caching verhindern
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Write(data)
	}
}

// renderPage rendert die übergebenen Daten mithilfe von HTML-Templates und gibt die generierte Seite aus.
// Es erwartet den Pfad zum Template-Verzeichnis und lädt mehrere Templates, die zusammen eine vollständige Seite bilden.
func renderPage(w io.Writer, data interface{}, tmpDir string) error {
	// Templates laden
	temp, err := template.ParseFiles(
		filepath.Join(tmpDir, "base.templ.html"),
		filepath.Join(tmpDir, "head.templ.html"),
		filepath.Join(tmpDir, "header.templ.html"),
		filepath.Join(tmpDir, "footer.templ.html"),
	)
	if err != nil {
		return fmt.Errorf("renderPage.Parsefiles: %w", err)
	}
	// Templates ausführen
	err = temp.ExecuteTemplate(w, "base", data)
	if err != nil {
		return fmt.Errorf("renderPage.ExecuteTemplate: %w", err)
	}
	return nil
}

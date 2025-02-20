package structure

import (
	"html/template"
	"m28909_Hausarbeit_WebprogGo/internal/database"
	"strings"

	"github.com/russross/blackfriday/v2"
)

// Page repräsentiert eine Seite mit Titel, Inhalt und Informationen, ob diese die Indexseite oder Corporate-Design-Seite ist.
type Page struct {
	Title       string
	IsIndexPage bool
	IsCDPage    bool
	Content     template.HTML
}

// LoadPage lädt eine Seite aus der Datenbank anhand des Dateinamens, setzt gegebenfalls den entsprechenden boolean
// wert auf true, wenn es sich um eine Index- oder Corporate-Design-Seite handelt.
func LoadPage(db *database.MongoDB, fileName string) (Page, error) {
	var p Page

	p.Title = fileName

	// Überprufung, ob p die Indexseite oder Corporate Design Seite ist, um im Template Scripte hinzuzufügen
	if strings.Contains(fileName, "index") {
		p.IsIndexPage = true
	} else if strings.Contains(fileName, "corporate") {
		p.IsCDPage = true
	}

	// Markdowndatei aus der Datenbank laden
	markdown, err := db.GetFile(fileName)
	if err != nil {
		return p, err
	}
	// Markdown in html umwandeln
	p.Content = template.HTML(blackfriday.Run(markdown))
	return p, nil
}

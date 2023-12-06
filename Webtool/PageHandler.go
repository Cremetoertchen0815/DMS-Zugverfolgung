package webtool

import (
	"html/template"
	"net/http"
)

// Manager for handling the web tool
type WebtoolHandler struct {
	mainPage *template.Template
}

// Creates an instance of the Webtool handler
func CreateWebtoolHandler() *WebtoolHandler {
	mainPage := template.Must(template.ParseFiles("Webtool/pages/main.gohtml"))
	return &WebtoolHandler{mainPage: mainPage}
}

func (page *WebtoolHandler) HandleMainPage(w http.ResponseWriter, r *http.Request) {
	err := page.mainPage.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

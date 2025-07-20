package fakenav

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

type Handler struct{}

type TemplateData struct {
	Base64Image string
}

func NewFakeNavHandler() *Handler {
	return &Handler{}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get current working directory: %v", err), http.StatusInternalServerError)
		return
	}
	//fmt.Println("Current working directory:", cwd)

	templatePath := filepath.Join(cwd, "htdocs", "fake-nav.html")

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		templatePath = filepath.Join(cwd, "app", "htdocs", "fake-nav.html")
		tmpl, err = template.ParseFiles(templatePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse template: %v", err), http.StatusInternalServerError)
			return
		}
	}

	imagePath := filepath.Join(cwd, "htdocs", "img", "nav1098x84.png")
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		imagePath = filepath.Join(cwd, "app", "htdocs", "img", "nav1098x84.png")
		imageBytes, err = os.ReadFile(imagePath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read image: %v", err), http.StatusInternalServerError)
			return
		}
	}

	base64Image := base64.StdEncoding.EncodeToString(imageBytes)

	data := TemplateData{
		Base64Image: "data:image/png;base64," + base64Image,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

}

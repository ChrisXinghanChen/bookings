package render

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ChrisXinghanChen/bookings/pkg/config"
	"github.com/ChrisXinghanChen/bookings/pkg/models"
)

var functions = template.FuncMap{}
var app *config.AppConfig

//NewTemplate sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a
}

//AddDefaultData add default data to evey template
func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

//renderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//get the template cache from the app config

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache.")
	}

	//buffer := new(bytes.Buffer)
	//_ = t.Execute(buffer, nil)
	td = AddDefaultData(td)
	err := t.Execute(w, td)
	if err != nil {
		fmt.Println("error writing template to browser;", err)
	}
}

//CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("page is currently:", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}

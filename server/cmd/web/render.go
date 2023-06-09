package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
)

type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	IsAuthenticated      int
	ErrorData            []string
	StripeSecretKey      string
	StripePublishableKey string
	CSRFToken            string
	Flash                string
	Warning              string
	Error                string
	API                  string
	CSSVersion           string
}

//go:embed templates
var templateFS embed.FS

var functions = template.FuncMap{
	"formatCurrency": formatCurrency,
}

func formatCurrency(n int) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("$%.2f", f)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = app.cfg.api
	td.StripeSecretKey = app.cfg.stripe.secret
	td.StripePublishableKey = app.cfg.stripe.key
	return td
}

func (app *application) renderTemplate(c *gin.Context, page string, td *templateData) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	_, templateInMap := app.templateCache[templateToRender]

	if templateInMap {
		t = app.templateCache[templateToRender]
	} else {
		t, err = app.parseTemplate(page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = app.addDefaultData(td, c.Request)

	err = t.Execute(c.Writer, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)

	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}

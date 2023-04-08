package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	driver "new-e-commerce/drivers"
	"new-e-commerce/models"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	api  string
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
}

var session *scs.SessionManager

type application struct {
	cfg           config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	database      models.DatabaseModel
	Session       *scs.SessionManager
	data          map[string]map[string]interface{}
}

func main() {
	gob.Register(models.User{})
	var cfg config

	cfg.api = "http://localhost:4001"
	cfg.port = 4000
	cfg.env = "development"

	cfg.db.dsn = "user=postgres dbname=e-commerce password=Matwyenko1_ host=localhost sslmode=disable binary_parameters=yes"
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDb(cfg.db.dsn)

	if err != nil {
		fmt.Println(err)
	}

	tc := make(map[string]*template.Template)

	session = scs.New()
	session.Lifetime = 24 * time.Hour

	data := make(map[string]map[string]interface{})

	app := &application{
		cfg:           cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		database:      models.DatabaseModel{DB: conn},
		Session:       session,
		data:          data,
	}

	router := gin.Default()

	app.SetupRoutes(router)

	if err := router.Run(":4000"); err != nil {
		app.errorLog.Printf("error running server on port: %v\n", err)
	}
}

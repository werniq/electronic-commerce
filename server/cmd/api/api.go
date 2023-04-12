package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	driver "new-e-commerce/drivers"
	"new-e-commerce/models"
	"os"
	"strconv"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	stripe struct {
		secret string
		key    string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
	}
	secretKey string
	client    string
}

var session *scs.SessionManager

type application struct {
	cfg           config
	infoLog       *log.Logger
	errorLog      *log.Logger
	templateCache map[string]*template.Template
	database      models.DatabaseModel
	Session       *scs.SessionManager
}

func main() {
	gob.Register(models.User{})
	var cfg config
	var err error

	if err = godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	cfg.port = 4000
	cfg.env = "development"

	cfg.db.dsn = "user=postgres dbname=e-commerce password=Matwyenko1_ host=localhost sslmode=disable binary_parameters=yes"
	cfg.stripe.key = os.Getenv("STRIPE_KEY")
	cfg.stripe.secret = os.Getenv("STRIPE_SECRET")

	cfg.smtp.password = os.Getenv("SMTPPass")
	cfg.smtp.username = os.Getenv("SMTPUser")
	cfg.smtp.host = os.Getenv("SMTPHost")
	cfg.smtp.port, err = strconv.Atoi(os.Getenv("SMTPPort"))
	cfg.secretKey = os.Getenv("SECRET_SIGNED_KEY")
	cfg.client = "http://localhost:4000"

	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	conn, err := driver.OpenDb(cfg.db.dsn)

	if err != nil {
		fmt.Println(err)
	}

	tc := make(map[string]*template.Template)

	app := &application{
		cfg:           cfg,
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		database:      models.DatabaseModel{DB: conn},
		Session:       session,
	}

	router := gin.Default()

	app.SetupApiRoutes(router)

	if err := router.Run(":4001"); err != nil {
		app.errorLog.Printf("error running server on port: %v\n", err)
	}
}

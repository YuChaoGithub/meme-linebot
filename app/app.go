package app

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/YuChaoGithub/meme-linebot/app/models"
	"github.com/YuChaoGithub/meme-linebot/config"
	"github.com/line/line-bot-sdk-go/linebot"

	_ "github.com/lib/pq" // PostgreSQL driver.
)

const (
	reconnectionInterval = 5
)

// App contains all the required models for the application.
type App struct {
	adminSecret   string
	memeModel     *models.MemeModel
	bot           *linebot.Client
	pageTemplates templateCache
}

// InitializeAndRun initializes the app with predefined configuration and run the app.
// It tries to reconnect to the database indefinitely when it fails.
func (a *App) InitializeAndRun(config *config.Config) {
	// PostgreSQL connection.
	var db *sql.DB

	dbTicker := time.NewTicker(reconnectionInterval * time.Second)

	// Database connection with retries.
Loop:
	for {
		select {
		case <-dbTicker.C:
			log.Println("Trying to establish connection with the database...")
			var err error
			db, err = sql.Open(config.DB.Dialect, config.DB.ConnectionURL)
			if err != nil {
				log.Println(err)
				continue Loop
			}

			if err = db.Ping(); err != nil {
				log.Println(err)
				continue Loop
			}

			// Successfully connected to the database.
			log.Println("Successfully connected to the database. The app is running.")
			dbTicker.Stop()
			defer db.Close()
			break Loop
		}
	}

	// Inject the DBs into the models.
	a.memeModel = &models.MemeModel{DB: db}

	// Start a new linebot client.
	bot, err := linebot.New(config.LineBot.ChannelSecret, config.LineBot.ChannelAccessToken)
	if err != nil {
		log.Println("Error creating a linebot client. Shutting down.")
		return
	}
	a.bot = bot

	// Compile page templates.
	a.pageTemplates, err = newTemplateCache([]string{"./ui/html/home.html"})
	if err != nil {
		log.Println("Error compiling templates.")
		return
	}

	// Admin secret.
	a.adminSecret = config.AdminSecret

	// Web server.
	server := &http.Server{
		Addr:         config.Server.Port,
		Handler:      a.routes(),
		IdleTimeout:  config.Server.IdleTimeout,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	log.Printf("Starting meme linebot server on %s\n", config.Server.Port)
	log.Fatal(server.ListenAndServe())
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.homepageHandler)
	mux.HandleFunc("/callback", a.callbackHandler)
	mux.HandleFunc("/add", a.addMeme)
	mux.HandleFunc("/delete", a.deleteMeme)

	// For static files on the home page.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

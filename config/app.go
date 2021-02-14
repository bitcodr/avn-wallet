package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type AppInterface interface {
	Environment()
	DB() (*sql.DB, error)
}

type App struct {
	DBName, DBUserName, DBPass string
	CurrentTime                time.Time
}

var (
	AppConfig  *viper.Viper
	LangConfig *viper.Viper
)

func (app *App) Init() {
	app.appConfig()
	app.langConfig()
	app.setAppConfig()
}

func (app *App) appConfig() {
	AppConfig = viper.New()
	AppConfig.SetConfigType("yaml")
	AppConfig.SetConfigName("config")
	AppConfig.AddConfigPath(os.Getenv("WALLET_SERVICE_ROOT_DIR"))
	err := AppConfig.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *App) langConfig() {
	LangConfig = viper.New()
	LangConfig.SetConfigType("yaml")
	LangConfig.SetConfigName("lang")
	LangConfig.AddConfigPath(os.Getenv("WALLET_SERVICE_ROOT_DIR") + "/resource/lang")
	err := LangConfig.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
}

func (app *App) setAppConfig() {
	app.DBName = AppConfig.GetString("DATABASES.MYSQL.DATABASE")
	app.DBUserName = AppConfig.GetString("DATABASES.MYSQL.USERNAME")
	app.DBPass = AppConfig.GetString("DATABASES.MYSQL.PASSWORD")
	app.CurrentTime = time.Now().UTC()
}

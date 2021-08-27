package config

import (
	"database/sql"

	di "github.com/isirfanm/online-store/inventory"
	pi "github.com/isirfanm/online-store/persistent/postgres/inventory"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type config struct {
	DriverName     string    `envconfig:"db_driver_name" default:"postgres"`
	DataSourceName string    `envconfig:"db_data_source_name" default:"postgresql://postgres:postgres@localhost:15432/postgres?sslmode=disable"`
	LogLevel       log.Level `envconfig:"log_level" default:"debug"`
	ReportCaller   bool      `envconfig:"report_caller" default:"false"`
}

var Cfg config

func SetupAll() {
	// read config
	if err := envconfig.Process("", &Cfg); err != nil {
		log.Fatal(err)
	}

	// Init DB connection
	db, err := sql.Open(Cfg.DriverName, Cfg.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	// make repo
	repo := pi.NewRepo(db)

	// setup inventory package
	di.Setup(db, repo)
}

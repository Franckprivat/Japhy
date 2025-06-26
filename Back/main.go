package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"github.com/japhy-tech/backend-test/database_actions"
	"github.com/japhy-tech/backend-test/internal"
)

const (
	MysqlDSN = "myuser:mypass@(mysql-test:3306)/myapp?parseTime=true"
	ApiPort  = "5000"
)

func main() {
	logger := charmLog.NewWithOptions(os.Stderr, charmLog.Options{
		Formatter:       charmLog.TextFormatter,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "🧑‍💻 backend-test",
		Level:           charmLog.DebugLevel,
	})

	err := database_actions.InitMigrator(MysqlDSN)
	if err != nil {
		logger.Fatal(err.Error())
	}

	msg, err := database_actions.RunMigrate("up", 0)
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info(msg)
	}

	db, err := sql.Open("mysql", MysqlDSN)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxIdleConns(0)

	err = db.Ping()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	logger.Info("Database connected")

	// Pour passer la base de données à l'application
	app := internal.NewApp(logger, db)

	r := mux.NewRouter()
	app.RegisterRoutes(r)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	logger.Info(fmt.Sprintf("Service started and listen on port %s", ApiPort))

	err = http.ListenAndServe(
		net.JoinHostPort("", ApiPort),
		r,
	)

	if err != nil {
		logger.Fatal("Erreur lors du démarrage du serveur", "error", err)
	}
}
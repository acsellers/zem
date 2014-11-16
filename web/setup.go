package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/acsellers/dr/migrate"
	"github.com/acsellers/zem/store"
)

func PostgresConn() string {
	if Config.DBHost == "localhost" {
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s",
			Config.DBUser,
			Config.DBPass,
			Config.DBName,
		)
	} else {
		if Config.DBPort == "" {
			Config.DBPort = "5432"
		}
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s port=%s",
			Config.DBUser,
			Config.DBPass,
			Config.DBName,
			Config.DBHost,
			Config.DBPort,
		)
	}
}

func SetupDB() {
	migrator := migrate.Database{
		// Later: DB, DBMS, Log
		Schema:     store.Schema,
		Translator: &store.AppConfig{},
		Log:        log.New(os.Stdout, "Migrate: ", 0),
	}

	switch Config.DBType {
	case "postgres":
		migrator.DBMS = migrate.Postgres
	default:
		log.Fatal("Unrecognized DBType, was:", Config.DBType)
	}

	ConnectDB()
	migrator.DB = Conn.DB
	err := migrator.Migrate()
	if err != nil {
		log.Fatal("Encountered error getting schema up to date:", err)
	}
	SetupManagement()
}

func ConnectDB() {
	var connString string
	switch Config.DBType {
	case "postgres":
		connString = PostgresConn()
	default:
		log.Fatal("Unrecognized DBType, was:", Config.DBType)
	}

	var err error
	Conn, err = store.Open(Config.DBType, connString)
	if err != nil {
		log.Fatal("Couldn't open Database, got error:", err)
	}

}

func SetupManagement() {
	if Conn.User.Management().Eq(true).Name().Eq("The Management").Count() > 1 {
		return
	}

	user := store.User{
		Name:       "The Management",
		Email:      "management@example.com",
		CreatedAt:  time.Now(),
		Management: true,
	}
	user.SetPassword("management_change_me")
	err := user.Save(Conn)
	if err != nil {
		log.Fatal("Couldn't create management user, error was:", err)
	}
}

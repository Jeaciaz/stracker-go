package pkg

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDb() {
	conn, err := sql.Open("sqlite", "file:/db/dev.db?mode=rwc&_busy_timeout=1000")
	if err != nil {
		log.Fatal("Failed to open DB: ", err)
	}
	db = conn

	// Handle shutdown
	cleanup := make(chan os.Signal, 1)
	signal.Notify(cleanup, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cleanup
		log.Print("Shutting down... Closing DB")
		closeDb()
		os.Exit(0)
	}()
}

func GetDb() *sql.DB {
	if db == nil {
		initDb()
	}
	return db
}

func closeDb() {
	if db != nil {
		_, _ = db.Exec("PRAGMA wal_checkpoint(FULL)")
		println("Closing DB")
		db.Close()
	}
}

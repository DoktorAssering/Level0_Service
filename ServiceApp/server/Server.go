package main

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"mainservice/ServiceApp/handler"
	"mainservice/ServiceApp/service"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var (
	db       *sql.DB
	nc       *nats.Conn
	isOnline bool
)

const (
	connStr = "user=postgres password=15109600 dbname=NewServicesDataBase sslmode=disable"
	natsURL = "localhost"
	addr    = "127.0.0.1:8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logsFolderPath := filepath.Join(currentDir, "..", "..", "..", "ResearchService", "Logs")
	backupsFolderPath := filepath.Join(currentDir, "..", "..", "..", "ResearchService", "Backups")
	logFile := filepath.Join(currentDir, "..", "..", "Logs", "Service.log")
	startBackupFile := filepath.Join(currentDir, "..", "..", "Backups", "startBackup.json")
	endBackupFile := filepath.Join(currentDir, "..", "..", "Backups", "endBackup.json")

	if err := createLogsFolder(logsFolderPath); err != nil {
		log.Fatal("Error creating Logs folder:", err)
		fmt.Print("Error creating Logs folder:", err)
	}

	if err := createBackupsFolder(backupsFolderPath); err != nil {
		log.Fatal("Error creating Backups folder:", err)
		fmt.Print("Error creating Backups folder:", err)
	}

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
		fmt.Print("Error opening log file:", err)
	}
	defer f.Close()

	log := log.New(f, "Server: ", log.LstdFlags|log.Lshortfile)

	if err := initDB(log); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := connectToNATS(log); err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	databaseService := service.NewDatabaseService(db, nc, startBackupFile, endBackupFile)
	databaseHandler := handler.NewDatabaseHandler(databaseService)

	handleError := func(err error) {
		log.Printf("Error occurred: %v", err)

		backupErr := databaseService.Backup(false)
		if backupErr != nil {
			log.Printf("Error during backup after error: %v", backupErr)
		}
	}
	handleError(err)

	if err := databaseService.RestoreCacheFromDB(); err != nil {
		log.Println("Error restoring cache from the database:", err)
	}
	fmt.Print("Successful data recovery from database")
	log.Print("Successful data recovery from database")

	if err := databaseService.Backup(true); err != nil {
		log.Println("Error creating backup:", err)
		fmt.Println("Error creating backup:", err)
	}
	fmt.Print("Successful creation of a backup")
	log.Print("Successful creation of a backup")

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/get-json", databaseHandler.GetInfo)
	http.HandleFunc("/add-json", databaseHandler.AddData)
	http.HandleFunc("/get-all-ids", databaseHandler.GetAllIDsHandler)

	go func() {
		isOnline = true
		fmt.Println("Server is online...")

		for {
			select {
			case <-time.After(1 * time.Minute):
				if isOnline {
					log.Println("Server online", time.Now())
				} else {
					log.Println("Server is offline", time.Now())
					fmt.Println("Server is offline", time.Now())
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		markServerOffline(log)
		cancel()
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}()

	fmt.Println("Starting the HTTP server...")
	log.Printf("HTTP server is starting with IP: %s", addr)
	server := &http.Server{Addr: addr}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
			fmt.Printf("Error starting server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	if err := databaseService.RestoreCacheFromDB(); err != nil {
		log.Println("Error restoring cache from the database:", err)
		fmt.Println("Error restoring cache from the database:", err)
	}
	if err := databaseService.Backup(false); err != nil {
		log.Println("Error creating backup:", err)
		fmt.Println("Error creating backup:", err)
	}
	fmt.Print("Successful creation of a backup")
	log.Print("Successful creation of a backup")
}

func initDB(logger *log.Logger) error {
	logger.Println("Connected to the database...")
	fmt.Println("Connected to the database...")
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	logger.Println("Connected to the database successfully")
	fmt.Println("Connected to the database successfully")
	return nil
}

func connectToNATS(logger *log.Logger) error {
	logger.Println("Connecting to NATS...")
	fmt.Println("Connecting to NATS...")
	var err error
	nc, err = nats.Connect(natsURL)
	if err != nil {
		logger.Println("Error connecting to NATS")
		fmt.Println("Error connecting to NATS")
		return err
	}
	logger.Println("Connected to NATS successfully")
	fmt.Println("Connected to NATS successfully")
	return nil
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../web/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func markServerOffline(logger *log.Logger) {
	isOnline = false
	db.Close()
	nc.Close()
	logger.Println("Server is shutting down...")
}

func createLogsFolder(pathFolder string) error {
	return os.MkdirAll(pathFolder, os.ModePerm)
}

func createBackupsFolder(pathFolder string) error {
	return os.MkdirAll(pathFolder, os.ModePerm)
}

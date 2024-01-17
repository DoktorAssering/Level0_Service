package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nats-io/nats.go"
)

func main() {
	natsURL := "localhost"

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logsFolderPath := filepath.Join(currentDir, "..", "..", "..", "ResearchService", "Logs")
	logFile := filepath.Join(currentDir, "..", "..", "Logs", "Service.log")

	if err := createLogsFolder(logsFolderPath); err != nil {
		log.Fatal("Error creating Logs folder:", err)
		fmt.Print("Error creating Logs folder:", err)
	}

	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Error opening log file:", err)
		fmt.Print("Error opening log file:", err)
	}
	defer f.Close()

	log := log.New(f, "Nats: ", log.LstdFlags|log.Lmicroseconds)

	log.Println("Connection...")
	fmt.Println("Connection...")
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("Error connecting to NATS: %v\n", err)
		fmt.Printf("Error connecting to NATS: %v\n", err)
		log.Fatal(err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("log", func(m *nats.Msg) {
		log.Printf("Received message: %s\n", string(m.Data))
		fmt.Printf("Received message: %s\n", string(m.Data))
	})
	if err != nil {
		log.Printf("Error subscribing to 'log' subject: %v\n", err)
		fmt.Printf("Error subscribing to 'log' subject: %v\n", err)
		return
	}

	log.Println("Listening for messages...")
	fmt.Println("Listening for messages...")

	select {}
}

func createLogsFolder(pathFolder string) error {
	return os.MkdirAll(pathFolder, os.ModePerm)
}

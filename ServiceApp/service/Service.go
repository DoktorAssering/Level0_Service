package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/nats-io/nats.go"
)

type Cache struct {
	OrderUID string `json:"order_uid"`
}

type DatabaseService struct {
	db              *sql.DB
	nc              *nats.Conn
	cache           map[int]Cache
	startBackupFile string
	endBackupFile   string
}

func NewDatabaseService(db *sql.DB, nc *nats.Conn, startBackupFile string, endBackupFile string) *DatabaseService {
	return &DatabaseService{
		db:              db,
		nc:              nc,
		cache:           make(map[int]Cache),
		startBackupFile: startBackupFile,
		endBackupFile:   endBackupFile,
	}
}

func (s *DatabaseService) RestoreCacheFromDB() error {
	ids, err := s.GetAllIDs()
	if err != nil {
		return fmt.Errorf("failed to get all IDs from the database: %v", err)
	}

	for _, id := range ids {
		number, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("failed to convert ID to integer: %v", err)
		}

		jsonData, err := s.GetInfo(number)
		if err != nil {
			return fmt.Errorf("failed to get data for ID %d: %v", number, err)
		}

		newCacheData := Cache{OrderUID: jsonData}
		s.cache[number] = newCacheData
	}

	return nil
}

func (s *DatabaseService) GetInfo(number int) (string, error) {
	var jsonData string
	err := s.db.QueryRow("select name_json_info from json_info where id_json_info = $1", number).Scan(&jsonData)
	if err != nil {
		return "", fmt.Errorf("Database query error")
	}

	message := fmt.Sprintf("A query to the database with the number has been executed ID: %d", number)
	s.nc.Publish("log", []byte(message))

	return jsonData, nil
}

func (s *DatabaseService) AddData(jsonData string) (int, error) {
	var newID int
	err := s.db.QueryRow("insert into json_info (Name_Json_Info) VALUES ($1) RETURNING ID_Json_Info", jsonData).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при добавлении данных в базу: %v", err)
	}

	newCacheData := Cache{OrderUID: jsonData}
	s.cache[newID] = newCacheData

	message := fmt.Sprintf("New data added: %s", jsonData)
	s.nc.Publish("log", []byte(message))

	return newID, nil
}

func (s *DatabaseService) GetAllIDs() ([]string, error) {
	rows, err := s.db.Query("select id_json_info from json_info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *DatabaseService) Backup(isStartBackup bool) error {
	var backupFile string

	if isStartBackup {
		backupFile = s.startBackupFile
	} else {
		backupFile = s.endBackupFile
	}

	_, err := os.Stat(backupFile)

	if err != nil && os.IsNotExist(err) {
		err := ioutil.WriteFile(backupFile, []byte("{}"), 0644)
		if err != nil {
			return fmt.Errorf("failed to create backup file: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check backup file existence: %v", err)
	}

	existingData, err := ioutil.ReadFile(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %v", err)
	}

	if len(existingData) == 0 {
		s.cache = make(map[int]Cache)
	} else {
		err = json.Unmarshal(existingData, &s.cache)
		if err != nil {
			return fmt.Errorf("failed to unmarshal existing backup data: %v", err)
		}
	}

	data, err := json.MarshalIndent(s.cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %v", err)
	}

	err = ioutil.WriteFile(backupFile, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write backup file: %v", err)
	}

	return nil
}

func (s *DatabaseService) RestoreCache(isStartRestore bool) error {
	var backupFile string

	if isStartRestore {
		backupFile = s.startBackupFile
	} else {
		backupFile = s.endBackupFile
	}

	data, err := ioutil.ReadFile(backupFile)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла бэкапа: %v", err)
	}

	err = json.Unmarshal(data, &s.cache)
	if err != nil {
		return fmt.Errorf("ошибка при декодировании данных кэша: %v", err)
	}

	return nil
}

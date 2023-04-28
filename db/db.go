package db

import (
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseConnection struct {
	*gorm.DB
	Servers []*PugServer
	Queues  []*QueueChannel
}

type PugServer struct {
	ID       int
	ServerID string `gorm:"unique"`
	Channels []QueueChannel
}
type QueueChannel struct {
	ID          int
	PugServerID int
	ChanID      string `gorm:"unique"`
}
type Match struct {
	ID     int
	Result int
}

func getDbPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	filePath := filepath.Join(dir, "dev.db")
	return filePath
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&PugServer{}, &QueueChannel{}, &Match{})
}

func GetDb(m bool) *DatabaseConnection {
	db, err := gorm.Open(sqlite.Open(getDbPath()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if m {
		migrate(db)
	}

	return &DatabaseConnection{DB: db}
}

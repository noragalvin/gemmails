package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // ...
)

// Model ...
type Model struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

var db *gorm.DB

// InitDB ..
func InitDB() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println("================== DB URL =====================")
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		// panic(err)
	}
	db = conn
	db.Debug().AutoMigrate(Shop{}, Subscriber{})
}

// OpenDB created connect database
func OpenDB() *gorm.DB {
	return db
}

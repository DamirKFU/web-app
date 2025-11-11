package tests

import (
	"app/config"
	"app/internal"
	"app/internal/core"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"testing"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type TestServer struct {
	Server *core.Server
	TX     *gorm.DB
}

var (
	testServer *core.Server
	once       sync.Once
)

func CreateTestDB(cfg *config.Config) string {
	testDBName := "testdb"

	mainConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password,
	)

	conn, err := sql.Open("postgres", mainConnStr)
	if err != nil {
		log.Fatalf("❌ cannot connect to PostgreSQL: %v", err)
	}
	defer conn.Close()

	var exists int
	err = conn.QueryRow(`SELECT 1 FROM pg_database WHERE datname = $1`, testDBName).Scan(&exists)
	if err == sql.ErrNoRows {
		log.Printf("🧱 Creating test database %s ...", testDBName)
		_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE %s;", testDBName))
		if err != nil {
			log.Fatalf("❌ cannot create test database %s: %v", testDBName, err)
		}
	} else if err != nil {
		log.Fatalf("❌ error checking test database: %v", err)
	} else {
		log.Printf("✅ Test database %s already exists", testDBName)
	}

	return testDBName
}

func GetTestServerWithTx(t *testing.T) *core.Server {
	server := GetTestServer()

	tx := server.DB.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	t.Cleanup(func() {
		tx.Rollback()
	})

	server.DB = tx

	return server
}

func GetTestServer() *core.Server {
	once.Do(func() {
		cfg := config.Load("../../")
		cfg.DB.Name = CreateTestDB(&cfg)
		testServer = internal.CreateApp(cfg)
	})
	return testServer
}

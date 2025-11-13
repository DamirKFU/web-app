package tests

import (
	"app/config"
	"app/internal"
	"app/internal/core"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

type TestServer struct {
	Server *core.Server
	TX     *gorm.DB
}

var (
	testServer *core.Server
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

func DropTestDB(cfg *config.Config, dbName string) {
	mainConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password,
	)
	db, err := sql.Open("postgres", mainConnStr)
	if err != nil {
		log.Fatalf("cannot connect to postgres to drop db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName))
	if err != nil {
		log.Fatalf("cannot drop test database %s: %v", dbName, err)
	}

	log.Printf("✅ Test database %s dropped", dbName)
}

func CleanDatabase(t *testing.T) {
	t.Helper()

	srv := GetTestServer()
	db := srv.DB
	models := srv.GetModels()

	t.Cleanup(func() {
		var tableNames []string

		for _, m := range models {
			stmt := &gorm.Statement{DB: db}
			if err := stmt.Parse(m); err != nil {
				t.Fatalf("failed to parse model %T: %v", m, err)
			}
			tableNames = append(tableNames, stmt.Table)
		}

		if len(tableNames) == 0 {
			return
		}

		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tableNames, ", "))
		if err := db.Exec(query).Error; err != nil {
			t.Fatalf("failed to truncate tables: %v", err)
		}
	})
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	cfg := config.Load("../../")
	cfg.DB.Name = CreateTestDB(&cfg)
	testServer = internal.CreateApp(cfg)
	testDBName := cfg.DB.Name

	code := m.Run()
	sqlDB, err := testServer.DB.DB()
	if err != nil {
		log.Fatalf("cannot get sql.DB from gorm.DB: %v", err)
	}
	sqlDB.Close()

	DropTestDB(&cfg, testDBName)

	os.Exit(code)
}

func GetTestServer() *core.Server {
	return testServer
}

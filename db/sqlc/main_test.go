package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/guncv/Simple-Bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// ✅ Load configuration
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// ✅ Open DB connection
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// ✅ Ping the database to verify connection
	if err = testDB.Ping(); err != nil {
		log.Fatalf("❌ Database connection error: %v", err)
	}

	fmt.Println("✅ Successfully connected to the test database!")

	// ✅ Initialize queries
	testQueries = New(testDB)

	// ✅ Run tests
	code := m.Run()

	// ✅ Close DB connection after tests
	if err := testDB.Close(); err != nil {
		log.Fatalf("❌ Failed to close database connection: %v", err)
	}

	// ✅ Exit with the test run status
	os.Exit(code)
}

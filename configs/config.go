package configs

import (
	"os"
)

var DatabasePath string

func LoadConfig() {
	DatabasePath = os.Getenv("SQLITE_DB_PATH")
}

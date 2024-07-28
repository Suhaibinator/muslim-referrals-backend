package config

import (
	"fmt"
	"os"
)

var DatabasePath string

func init() {
	DatabasePath = os.Getenv("SQLITE_DB_PATH")
	fmt.Println(DatabasePath)
}

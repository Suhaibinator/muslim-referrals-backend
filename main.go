package main

import (
	"fmt"
	"io"

	//	"muslim-referrals-backend/internal/database"
	models "muslim-referrals-backend/internal/models"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	_ "ariga.io/atlas-provider-gorm/gormschema"
)

func main() {

	stmts, err := gormschema.New("sqlite").Load(&models.User{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}

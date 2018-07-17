package main

import (
	"fmt"
	"os"

	"github.com/himidori/pm/db"
	"github.com/himidori/pm/utils"
)

func main() {
	dbPath := os.Getenv("HOME") + "/.PM/db.sqlite"
	if !utils.PathExists(dbPath) {
		err := db.InitBase()
		if err != nil {
			fmt.Println("failed to initialize db:", err)
		}
		os.Exit(0)
	}
	initArgs()
	parseArgs()
}

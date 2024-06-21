package main

import (
	"flag"
)

const databasePath = "./test.json"

var database *Database

func init() {
	data, err := loadFileContents(databasePath)
	if err != nil {
		panic(err)
	}
	database, err = CreateDatabase(databasePath, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	projectPtr := flag.String("project", "", "the project we would like to see entries for")
	specificDate := flag.String("date", "", "the date we want to show entries for")
	flag.Parse()

	ProcessCommand(Flags{ProjectName: *projectPtr, SpecificDate: *specificDate})
}

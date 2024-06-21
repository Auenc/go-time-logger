package main

const databasePath = "./test.json"

func main() {
	// data, err := loadFileContents(databasePath)
	// if err != nil {
	// 	panic(err)
	// }
	// database, err := CreateDatabase(databasePath, data)
	// if err != nil {
	// 	panic(err)
	// }
	// // add some entries
	// now := time.Now()
	// then := now.Add(time.Hour)
	// database.Add("Saving some files", now, then)

	// entry, err := database.DeleteById(8)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%s - you are the weakest entry, goodbyte\n", entry.Name)

	// fmt.Println("we did it! Gwych!")
	ProcessCommand()
}

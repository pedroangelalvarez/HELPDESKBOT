package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)
import "Conexion"

func main() {
	os.MkdirAll("./data/reg", 0755)
	os.Create("./data/reg/data.db")

	db, err := sql.Open("sqlite3", "./data/reg/data.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = db.Exec("CREATE  TABLE `sessions` (ID TEXT PRIMARY KEY, FECHA DATETIME, ACTIVO INTEGER)")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = db.Exec("INSERT INTO `sessions` values('" + "51966614614s.whatsapp.net" + "','" + "2022-01-01 10:00:00" + "'," + "1" + ")")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, errROW := db.Query("SELECT * FROM sessions ")
	if errROW != nil {
		fmt.Println("errROW")
		log.Println(errROW)

	}

	for rows.Next() {
		fmt.Println("where are the rows")
		var ID string
		var FECHA string
		var ACTIVO int
		rows.Scan(&ID, &FECHA, &ACTIVO)
		fmt.Println(ID, FECHA, ACTIVO)

	}
	rows.Close()

	db.Close()
}

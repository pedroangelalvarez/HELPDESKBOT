package main

import (
  "database/sql"
  "fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Conexion struct {
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func AsignandoSesion(id text) {

  var existe bool
  db, err := sql.Open("sqlite3", "./sesiones.db")
	checkErr(err)
  rows, errROW := db.Query("SELECT * FROM sessions WHERE ID "=id)
	checkErr(errROW)
	for rows.Next() {
		var ID string
		var FECHA string
		var ACTIVO int
		rows.Scan(&ID, &FECHA, &ACTIVO)
		fmt.Println(ID, FECHA, ACTIVO)
    existe = true
	}
  rows.Close()

  if existe{
    _, err = db.Exec("UPDATE `sessions` SET FECHA='"+currentTime.Format("2006-01-02 15:04:05")+"' WHERE ID="+id)
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
  else{
    _, err = db.Exec("INSERT INTO `sessions` values('" + "51966614614s.whatsapp.net" + "','" + currentTime.Format("2006-01-02 15:04:05") + "'," + "1" + ")")
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
	db.Close()
}

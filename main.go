package main

import (
	"encoding/gob"
	"database/sql"
  "fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

type person struct {
	id          string
	fecha time.Time
  activo int
}

type waHandler struct {
	wac       *whatsapp.Conn
	startTime uint64
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}


func AsignandoSesion(id string) {
  currentTime := time.Now()
  var existe bool
  db, err := sql.Open("sqlite3", "./sesiones.db")
	checkErr(err)
  rows, errROW := db.Query("SELECT * FROM sessions WHERE ID ="+id)
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
  } else{
    _, err = db.Exec("INSERT INTO `sessions` values('" + "51966614614s.whatsapp.net" + "','" + currentTime.Format("2006-01-02 15:04:05") + "'," + "1" + ")")
    if err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
	db.Close()
}

func (wh *waHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "error caught in handler: %v\n", err)
}

// HandleTextMessage receives whatsapp text messages and checks if the message was send by the current
// user, if it does not contain the keyword '@echo' or if it is from before the program start and then returns.
// Otherwise the message is echoed back to the original author.
func (wh *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.FromMe || message.Info.Timestamp < wh.startTime {
		return
	}

	if strings.Contains(strings.ToLower(message.Text), "hola") || strings.Contains(strings.ToLower(message.Text), "buenos dias") || strings.Contains(strings.ToLower(message.Text), "buenas tardes") || strings.Contains(strings.ToLower(message.Text), "buenas noches") {
		
    AsignandoSesion(message.Info.RemoteJid)
    msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: message.Info.RemoteJid,
			},
			Text: "Hola, Como puedo Ayudarte?",
		}

		if _, err := wh.wac.Send(msg); err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
		}

    msg = whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: message.Info.RemoteJid,
			},
			Text: "Para comenzar indicame qué tipo de problema tienes?",
		}

		if _, err := wh.wac.Send(msg); err != nil {
			fmt.Fprintf(os.Stderr, "error sending message: %v\n", err)
		}

		fmt.Printf("echoed message '%v' to user %v\n", message.Text, message.Info.RemoteJid)
	}
}

func login(wac *whatsapp.Conn) error {
	session, err := readSession()
  
	if err == nil {
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring session failed: %v", err)
		}
	} else {
		qr := make(chan string)

		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()

		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v", err)
		}
	}
	if err = writeSession(session); err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}

	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}

	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&session); err != nil {
		return session, err
	}

	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(session); err != nil {
		return err
	}

	return nil
}

func main() {
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}

	wac.AddHandler(&waHandler{wac, uint64(time.Now().Unix())})

	if err = login(wac); err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	<-time.After(60 * time.Minute)
}

package db

import (
	"database/sql"
	"log"
	"os"

	// Register Mysql driver
	_ "github.com/go-sql-driver/mysql"

	"github.com/tttest25/goreceptionloader/logger"
)

var (
	// Logger variable for logging
	l   *log.Logger
	db  *sql.DB
	dsn string // connection string
)

// Config - get config
func Config() {
	dsn = os.Getenv("DSN")
}

// MustConnectDB check connection
func MustConnectDB() {
	if err := ConnectDatabase(); err != nil {
		l.Fatal(err)
		// panic(err)
	}
	l.Printf("DB connected %#v ", db.Driver())
}

// ConnectDatabase connect instance
func ConnectDatabase() (err error) {

	if db, err = sql.Open("mysql", dsn); err != nil {
		return
	}
	db.SetMaxIdleConns(10)
	err = db.Ping()
	return
}

// CloseDatabase close instance
func CloseDatabase() {

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

/*CreateCon Create mysql connection*/
/*func CreateCon() *sql.DB {
	db1, err := sql.Open("mysql", "test32:PbGFTmw3uzpNWRnT@tcp(172.19.193.159)/modx_reception")
	if err != nil {
		l.Println(err.Error())
	}
	l.Println("db1 is connected")
	// defer db1.Close()
	// make sure connection is available
	err = db1.Ping()
	l.Printf("Error = %#v ", err)
	if err != nil {
		l.Println("MySQL db is not connected")
		l.Fatal( err.Error())
	}
	db = db1
}*/

// InitDB - проверить подключение
func InitDB() {
	defer func() {
		if e := recover(); e != nil {
			l.Println(e)
		}
	}()

	CreateTable()
	// Insert(&Person{Name: "Ale", Phone: "+55 53 1234 4321"})
	// Insert(&Person{Name: "Cla", Phone: "+66 33 1234 5678"})
}

// CreateTable Create table
func CreateTable() {
	stmt, err := db.Prepare(createTable)
	if err != nil {
		l.Print(err)
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		l.Fatal(err)
		// panic(err)
	}
}

// Drop table
func Drop() {
	stmt, err := db.Prepare(dropTable)
	if err != nil {
		l.Fatal(err)
		// panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		l.Fatal(err)
		// panic(err)
	}

	CreateTable()
}

func init() {
	l = logger.ReturnLogger("db")
	Config()
	MustConnectDB()
	InitDB()
}

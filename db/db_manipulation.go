package db

import (
	"database/sql"
	"log"
)

// Smddb Struct for use with DB
type Smddb struct {
	RequestID        int64
	Number           string
	DtModified       string
	DepartmentID     string
	DepartmentName   string
	Format           string
	FormatName       string
	IsDirect         int64
	CreateDate       string
	Name             string
	Address          string
	Email            string
	ReceiveDate      string
	DispatchDate     string
	UploadDate       string
	ExceptionMessage string
	Questions        string
}

// Insert - insert data
func Insert(smd *Smddb) int64 {
	// smd.RequestID > 0 ||
	// if GetResult(smd.RequestID) != "" {
	// if already exist skip
	if GetResultF(smd.RequestID) {
		l.Printf("insert conflict - skip %d", smd.RequestID)
		return 0
	}

	stmt, err := db.Prepare(peopleInsert)
	if err != nil {

		l.Fatal(err)
		// panic(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		&smd.RequestID,
		&smd.Number,
		&smd.DtModified,
		&smd.DepartmentID,
		&smd.DepartmentName,
		&smd.Format,
		&smd.FormatName,
		&smd.IsDirect,
		&smd.CreateDate,
		&smd.Name,
		&smd.Address,
		&smd.Email,
		&smd.ReceiveDate,
		&smd.DispatchDate,
		&smd.UploadDate,
		&smd.ExceptionMessage,
		&smd.Questions,
	)
	if err != nil {
		l.Printf("Error data: %#v", smd)
		l.Fatal(err)
		// panic(err)
	}
	count, err2 := res.RowsAffected()

	if err2 != nil {
		l.Printf(err2.Error())
	}
	return count

}

// List get list data
func List() []*Smddb {
	rows, err := db.Query(smdQuery)
	if err != nil {
		l.Fatal(err)
		// panic(err)
	}
	defer rows.Close()

	list := []*Smddb{}
	for rows.Next() {
		m := new(Smddb)
		if err := rows.Scan(
			&m.RequestID,
			&m.Number,
			&m.DtModified,
			&m.DepartmentID,
			&m.DepartmentName,
			&m.Format,
			&m.FormatName,
			&m.IsDirect,
			&m.CreateDate,
			&m.Name,
			&m.Address,
			&m.Email,
			&m.ReceiveDate,
			&m.DispatchDate,
			&m.UploadDate,
			&m.ExceptionMessage,
			&m.Questions,
		); err != nil {
			log.Println(err)
			l.Fatal(err)
			// panic(err)
		}

		list = append(list, m)
	}

	return list
}

// GetResultF - check existance of row with requestid
func GetResultF(id int64) bool {
	var sid int64
	err := db.QueryRow("select sid FROM smd_data WHERE RequestID =  ?", id).Scan(&sid)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Fatal(err)
		}

		return false
	}
	return true
}

// GetResult - get string
func GetResult(id int64) string {
	rows, err := db.Query(smdQuery+" WHERE RequestID = ?", id)
	if err != nil {
		l.Fatal(err)
		// panic(err)
	}
	defer rows.Close()

	list := []*Smddb{}
	for rows.Next() {
		m := new(Smddb)
		if err := rows.Scan(
			&m.RequestID,
			&m.Number,
			&m.DtModified,
			&m.DepartmentID,
			&m.DepartmentName,
			&m.Format,
			&m.FormatName,
			&m.IsDirect,
			&m.CreateDate,
			&m.Name,
			&m.Address,
			&m.Email,
			&m.ReceiveDate,
			&m.DispatchDate,
			&m.UploadDate,
			&m.ExceptionMessage,
			&m.Questions,
		); err != nil {
			l.Fatal(err)
			// panic(err)
		}

		list = append(list, m)
	}

	if len(list) > 0 {
		return list[0].Number
	}
	return ""
}

var peopleInsert = `
INSERT INTO smd_data(
	RequestID, Number, DtModified, DepartmentID, DepartmentName, Format, FormatName, IsDirect, CreateDate, Name, Address, Email, ReceiveDate, DispatchDate, UploadDate, ExceptionMessage, Questions) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
var smdQuery = `
SELECT
	RequestID,
	Number,
	DtModified,
	DepartmentID,
	DepartmentName,
	Format,
	FormatName,
	IsDirect,
	CreateDate,
	Name,
	Address,
	Email,
	ReceiveDate,
	DispatchDate,
	UploadDate,
	ExceptionMessage,
	Questions
FROM
    smd_data
`

//     --title 			 VARCHAR(255) NOT NULL,

var createTable = `
CREATE TABLE IF NOT EXISTS smd_data (
    sid INT AUTO_INCREMENT PRIMARY KEY,
    RequestID        int,
    Number           VARCHAR(40),
    DtModified       VARCHAR(30),
    DepartmentID     VARCHAR(36),
    DepartmentName   VARCHAR(250),
    Format           VARCHAR(20),
    FormatName       VARCHAR(20),
    IsDirect         int,
    CreateDate       VARCHAR(30),
    Name             VARCHAR(200),
    Address          VARCHAR(200),
    Email            VARCHAR(100),
    ReceiveDate      VARCHAR(30),
	DispatchDate     VARCHAR(30),
	UploadDate		 VARCHAR(30),
    ExceptionMessage VARCHAR(300),
    Questions        VARCHAR(1024) default '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    CHECK (Questions IS NULL OR JSON_VALID(Questions))
)  ENGINE=INNODB;
`
var dropTable = `
DROP TABLE smd_data;
`

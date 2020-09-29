package db

import (
	"log"
)

// Smddb Struct for use with DB
type Smddb struct {
	RequestID        int64
	DepartmentID     string
	Format           string
	IsDirect         int64
	Number           string
	CreateDate       string
	Name             string
	Address          string
	Email            string
	ReceiveDate      string
	UploadDate       string
	DispatchDate     string
	ExceptionMessage string
}

// Insert - insert data
func Insert(smd *Smddb) int64 {
	// smd.RequestID > 0 ||
	if GetResult(smd.RequestID) != "" {
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
		&smd.DepartmentID,
		&smd.Format,
		&smd.IsDirect,
		&smd.Number,
		&smd.CreateDate,
		&smd.Name,
		&smd.Address,
		&smd.Email,
		&smd.ReceiveDate,
		&smd.UploadDate,
		&smd.DispatchDate,
		&smd.ExceptionMessage,
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
			&m.DepartmentID,
			&m.Format,
			&m.IsDirect,
			&m.Number,
			&m.CreateDate,
			&m.Name,
			&m.Address,
			&m.Email,
			&m.ReceiveDate,
			&m.UploadDate,
			&m.DispatchDate,
			&m.ExceptionMessage,
		); err != nil {
			log.Println(err)
			l.Fatal(err)
			// panic(err)
		}

		list = append(list, m)
	}

	return list
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
			&m.DepartmentID,
			&m.Format,
			&m.IsDirect,
			&m.Number,
			&m.CreateDate,
			&m.Name,
			&m.Address,
			&m.Email,
			&m.ReceiveDate,
			&m.UploadDate,
			&m.DispatchDate,
			&m.ExceptionMessage,
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
INSERT INTO smd_data(RequestID, DepartmentID, Format, IsDirect, Number, CreateDate, Name, Address, Email, ReceiveDate, UploadDate, DispatchDate, ExceptionMessage) 
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`
var smdQuery = `
SELECT
	RequestID,
	DepartmentID,
	Format,
	IsDirect,
	Number,
	CreateDate,
	Name,
	Address,
	Email,
	ReceiveDate,
	UploadDate,
	DispatchDate,
	ExceptionMessage
FROM
    smd_data
`

//     --title 			 VARCHAR(255) NOT NULL,

var createTable = `
CREATE TABLE IF NOT EXISTS smd_data (
    sid INT AUTO_INCREMENT PRIMARY KEY,
	RequestID        int,
	DepartmentID     VARCHAR(36),
	Format           VARCHAR(20),
	IsDirect         int,
	Number           VARCHAR(40),
	CreateDate       VARCHAR(30),
	Name             VARCHAR(200),
	Address          VARCHAR(200),
	Email            VARCHAR(100),
	ReceiveDate      VARCHAR(30),
	UploadDate       VARCHAR(30),
	DispatchDate     VARCHAR(30),
	ExceptionMessage VARCHAR(300),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)  ENGINE=INNODB;
`
var dropTable = `
DROP TABLE smd_data;
`

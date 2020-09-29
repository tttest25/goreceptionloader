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
func Insert(smd *Smddb) {

	if smd.RequestID > 0 || GetResult(smd.RequestID) != "" {
		panic("insert conflict")
	}

	stmt, err := db.Prepare(peopleInsert)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
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
		panic(err)
	}

}

// List get list data
func List() []*Smddb {
	rows, err := db.Query(smdQuery)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	list := []*Smddb{}
	for rows.Next() {
		m := new(Smddb)
		if err := rows.Scan(
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
			panic(err)
		}

		list = append(list, m)
	}

	return list
}

// GetResult - get string
func GetResult(id int64) string {
	rows, err := db.Query(smdQuery+" WHERE RequestID = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	list := []*Smddb{}
	for rows.Next() {
		m := new(Smddb)
		if err := rows.Scan(
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
			panic(err)
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
    *
FROM
    smd_data
`

var createTable = `
CREATE TABLE IF NOT EXISTS smd_data (
    sid INT AUTO_INCREMENT PRIMARY KEY,
    title 			 VARCHAR(255) NOT NULL,
	RequestID        int,
	DepartmentID     VARCHAR(36),
	Format           VARCHAR(20),
	IsDirect         int,
	Number           VARCHAR(40),
	CreateDate       VARCHAR(20),
	Name             VARCHAR(100),
	Address          VARCHAR(200),
	Email            VARCHAR(100),
	ReceiveDate      VARCHAR(20),
	UploadDate       VARCHAR(20),
	DispatchDate     VARCHAR(20),
	ExceptionMessage VARCHAR(20),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)  ENGINE=INNODB;
`
var dropTable = `
DROP TABLE smd_data;
`

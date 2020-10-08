package db

import (
	"database/sql"
	"log"
	"time"
)

// Smddb Struct for use with DB
type Smddb struct {
	RequestID        int64
	Number           string
	DtModified       time.Time
	DepartmentID     string
	DepartmentName   string
	Format           string
	FormatName       string
	IsDirect         int64
	CreateDate       time.Time
	Name             string
	Address          string
	Email            string
	ReceiveDate      time.Time
	DispatchDate     time.Time
	UploadDate       time.Time
	ExceptionMessage string
	Questions        string
}

// Insert - insert data
func Insert(smd *Smddb) int64 {
	// smd.RequestID > 0 ||
	// if GetResult(smd.RequestID) != "" {
	// if already exist skip

	// if GetResultF(smd.RequestID) {
	// 	// l.Printf("insert conflict - skip %d", smd.RequestID)
	// 	return 0
	// }

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

// Update - update data
func Update(sid int64, smd *Smddb) int64 {
	stmt, err := db.Prepare(smdUpdate)
	if err != nil {
		l.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
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

		sid,
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

// CheckModified - check modified of row with requestid
func CheckModified(id int64, dt time.Time) (int64, int64) {
	var (
		status    int64 = -1 // -1 error , 0 not changed , 1 new , 2 update
		sid, diff sql.NullInt64
		dt1, dt2  sql.NullTime //dt1,
	)

	//l.Printf(" Date -> %d / %s", id, dt.Format("2006-01-02 15:04:05"))
	//with q as (select 609 as RequestID1, STR_TO_DATE("2020-10-01 22:15:01","%Y-%m-%d %H:%i:%s") as DtMod1)
	// ,s as (select * from smd_data s,q where s.RequestId=q.RequestID1)
	// select s.sid,q.DtMod1,s.DtModified,q.DtMod1-s.DtModified as diff
	// from q,s;
	// search and check if row needed modified STR_TO_DATE(DtModified,'%Y-%m-%d %H:%i:%s.%f')
	err := db.QueryRow(`with q as (select ? as RequestID1, STR_TO_DATE(?,'%Y-%m-%d %H:%i:%s') as DtMod1)    
	, s as (select * from  q left join smd_data s on s.RequestId=q.RequestID1) 
	select RequestID1,s.sid,s.DtMod1,s.DtModified,TIMESTAMPDIFF(second,s.DtModified,s.DtMod1) as diff 
	from s ;`,
		id, dt.Format("2006-01-02 15:04:05")).Scan(&id, &sid, &dt1, &dt2, &diff)

	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			l.Fatal(err)
		}
		// l.Printf(" %d -> return 0 OK ", id)
		// return status, sid.Int64
	}

	// -1 error , 0 not changed , 1 new , 2 update  status
	if sid.Valid {
		// sid is not bull string exist
		if diff.Valid && diff.Int64 > 1 {
			status = 2
		} else {
			status = 0
		}
	} else {
		// sid null -> new string
		status = 1
	}

	// l.Printf(" #DEBUG id  %d sid %d -> status %d  DtMod1 %#v DtModified %s diff_dt %v error  %#v",
	// 	id, sid.Int64, status, dt1.Time.Format(time.RFC3339), dt2.Time.Format(time.RFC3339), diff, err)

	return status, sid.Int64
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

var smdUpdate = `
update  smd_data
    set
	Number 			= ?
	,DtModified 		= ?
	,DepartmentID 		= ?
	,DepartmentName 	= ?
	,Format 			= ?
	,FormatName 		= ?
	,IsDirect 			= ?
	,CreateDate 		= ?
	,Name 				= ?
	,Address 			= ?
	,Email 				= ?
	,ReceiveDate 		= ?
	,DispatchDate 		= ?
	,UploadDate 		= ?
	,ExceptionMessage	= ?
	,Questions  		= ?
where 
 sid = ?
`

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
    DtModified       DATETIME,
    DepartmentID     VARCHAR(36),
    DepartmentName   VARCHAR(250),
    Format           VARCHAR(20),
    FormatName       VARCHAR(20),
    IsDirect         int,
    CreateDate       DATETIME,
    Name             VARCHAR(200),
    Address          VARCHAR(200),
    Email            VARCHAR(100),
    ReceiveDate      DATETIME,
	DispatchDate     DATETIME,
	UploadDate		 DATETIME,
    ExceptionMessage VARCHAR(300),
    Questions        VARCHAR(1024) default '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    CHECK (Questions IS NULL OR JSON_VALID(Questions))
)  ENGINE=INNODB;
`
var dropTable = `
DROP TABLE smd_data;
`

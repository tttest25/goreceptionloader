package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/tttest25/goreceptionloader/db"
	"github.com/tttest25/goreceptionloader/logger"
)

// Smd struct  [] for JSON smd
type Smd []struct {
	RequestID        string      `json:"requestId"`
	DepartmentID     string      `json:"departmentId"`
	Format           string      `json:"format"`
	IsDirect         string      `json:"isDirect"`
	Number           string      `json:"number"`
	CreateDate       string      `json:"createDate"`
	Name             string      `json:"name"`
	Address          string      `json:"address"`
	Email            string      `json:"email"`
	ReceiveDate      string      `json:"receiveDate"`
	UploadDate       string      `json:"uploadDate"`
	DispatchDate     string      `json:"dispatchDate"`
	ExceptionMessage interface{} `json:"exceptionMessage"`
}

type stat struct {
	all      int64
	inserted int64
	skiped   int64
}

var (
	// Logger variable for logging
	l *log.Logger

// log nage

)

func get() *Smd {
	l.Println("1. Performing Http Get...")
	resp, err := http.Get("https://smd.permkrai.ru/IPCP/HandlingReportPlugin/Api/analytics/ver.0.1/message")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	// bodyString := string(bodyBytes)
	// fmt.Println("API Response as String:\n" + bodyString)

	// Convert response body to Todo struct
	var todoStruct Smd
	json.Unmarshal(bodyBytes, &todoStruct)
	l.Printf("API Response as struct %#v\n", len(todoStruct))
	l.Printf("Record %#v\n", todoStruct[0])
	return &todoStruct

}

func main() {

	// defer db.Close()
	l = logger.ReturnLogger("main")
	l.Printf("--- Start ")
	smd := get() // get json data

	a := new(db.Smddb)
	stat := new(stat)

	for i, s := range *smd {
		// fmt.Println(i, s)
		// l.Printf("Get result `354` ='%s'", db.GetResult(354))
		id, _ := strconv.ParseInt(s.RequestID, 10, 64)
		a = &db.Smddb{
			RequestID:        id,
			DepartmentID:     s.DepartmentID,
			Format:           s.Format,
			IsDirect:         0,
			Number:           s.Number,
			CreateDate:       s.CreateDate,
			Name:             s.Name,
			Address:          s.Address,
			Email:            s.Email,
			ReceiveDate:      s.ReceiveDate,
			UploadDate:       s.UploadDate,
			DispatchDate:     s.DispatchDate,
			ExceptionMessage: fmt.Sprintf("%d : %v", i, s.ExceptionMessage),
		}
		stat.all++
		if db.Insert(a) > 0 {
			stat.inserted++
		} else {
			stat.skiped++
		}
	}
	l.Printf("Get result from smd  ='%#v'", stat)

	l.Printf("=== Successfully stop elapsed %dms\n \n \n", logger.TimeElapsed()/1000)

}

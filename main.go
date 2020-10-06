package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tttest25/goreceptionloader/db"
	"github.com/tttest25/goreceptionloader/logger"
)

// Smd struct  [] for JSON smd
type Smd []struct {
	RequestID        string      `json:"requestId"`
	Number           string      `json:"number"`
	DtModified       string      `json:"dt_modified"`
	DepartmentID     string      `json:"departmentId"`
	DepartmentName   string      `json:"departmentName"`
	Format           string      `json:"format"`
	FormatName       string      `json:"formatName"`
	IsDirect         bool        `json:"isDirect"`
	CreateDate       string      `json:"createDate"`
	Name             string      `json:"name"`
	Address          string      `json:"address"`
	Email            string      `json:"email"`
	ReceiveDate      string      `json:"receiveDate"`
	DispatchDate     string      `json:"dispatchDate"`
	UploadDate       string      `json:"uploadDate"`
	ExceptionMessage interface{} `json:"exceptionMessage"`
	Questions        []struct {
		Code               string `json:"code"`
		Status             string `json:"status"`
		QuestionStatusName string `json:"questionStatusName"`
		IncomingDate       string `json:"incomingDate"`
	} `json:"questions"`
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
	//resp, err := http.Get("https://smd.permkrai.ru/IPCP/HandlingReportPlugin/Api/analytics/ver.0.1/message")
	resp, err := http.Get("https://smd.permkrai.ru/IPCP/HandlingReportPlugin/Api/analytics/ver.0.2/requests/f1ae1eef-16ea-44cb-b77f-6b978ee4075d")
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

	defer db.CloseDatabase()
	l = logger.ReturnLogger("main")
	l.Printf("--- Start ")
	smd := get() // get json data

	a := new(db.Smddb)
	stat := new(stat)

	start := make(chan bool)
	var done sync.WaitGroup
	done.Add(2)

	val := *smd
	l1 := len(val) / 2

	for i := 0; i < 2; i++ {
		go func() {
			<-start
			for i := 0; i < n; i++ {
				rows, err := db.Query("SELECT 1 as `id`")
				if err != nil {
					log.Fatal(err)
				}
				rows.Close()
			}
			done.Done()
		}()
	}

	at1 := time.Now()
	close(start)
	done.Wait()

	fmt.Printf("sqlasyn finished - %v...\r\n", time.Since(at1))

	for _, s := range val[1:l1] {
		// fmt.Println(i, s)
		// l.Printf("Get result `354` ='%s'", db.GetResult(354))

		// parse data for db
		id, _ := strconv.ParseInt(s.RequestID, 10, 64)
		q, err := json.Marshal(s.Questions)
		if err != nil {
			l.Fatalf("Fatal error questions %s", err)
		}
		intIsdirect := int64(0)
		if s.IsDirect {
			intIsdirect = 1
		}
		a = &db.Smddb{
			RequestID:        id,
			Number:           s.Number,
			DtModified:       s.DtModified,
			DepartmentID:     s.DepartmentID,
			DepartmentName:   s.DepartmentName,
			Format:           s.Format,
			FormatName:       s.FormatName,
			IsDirect:         intIsdirect,
			CreateDate:       s.CreateDate,
			Name:             s.Name,
			Address:          s.Address,
			Email:            s.Email,
			ReceiveDate:      s.ReceiveDate,
			DispatchDate:     s.DispatchDate,
			UploadDate:       s.UploadDate,
			ExceptionMessage: fmt.Sprintf("%v", s.ExceptionMessage),
			Questions:        string(q),
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

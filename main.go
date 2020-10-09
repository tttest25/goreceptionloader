package main

import (
	"encoding/json"
	"flag"
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

//  layout for parse date  "2020-08-07 12:30:01.995389"
const (
	layoutDateTimeSmd = "2006-01-02 15:04:05.999999"
	layoutDtSmd       = "2006-01-02 15:04:05.999999"
)

// Smd struct  [] for JSON smd
type Smd struct {
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
	updated  int64
	skiped   int64
}

var (
	// Logger variable for logging
	l *log.Logger
	// command line arguments = flag
	nqueries, nconns int
)

func init() {
	// config configure parralelism
	flag.IntVar(&nqueries, "n", 10, "number of queries") // !! not nedded , use if need to limit queries in parallel
	flag.IntVar(&nconns, "c", 10, "number of connections")
	ptrV := flag.Bool("v", false, "print stdout")
	flag.Parse()
	logger.SetStdin(*ptrV)

}

func main() {

	defer db.CloseDatabase()
	defer logger.LogCloseFile()

	l = logger.ReturnLogger("main")

	fmt.Printf("Stdout flag %#v \n", logger.GetStdin())
	// set maximim connection and max idle connection
	db.SetNumCons(nconns)

	l.Printf("--- Start ")
	smd := get() // get json data

	// prepeare parallel
	start := make(chan bool)
	chStat := make(chan *stat, nconns)
	var done sync.WaitGroup
	done.Add(nconns)

	// split data for parallel working
	// // Split data for 2 peaces
	// val := *smd
	// l1 := len(*smd) / 2
	var divided [][]Smd
	chunkSize := (len(*smd) + nconns - 1) / nconns

	for i := 0; i < len(*smd); i += chunkSize {
		end := i + chunkSize
		if end > len(*smd) {
			end = len(*smd)
		}
		divided = append(divided, (*smd)[i:end])
	}

	if len(divided) < nconns {
		divided = append(divided, []Smd{})
	}

	l.Printf(" Prepeared %d slices for paralell process.", len(divided))

	// Pass slices for process parallel in ncons gorutine
	for i := 0; i < nconns; i++ {

		go func(val []Smd, i int) {
			<-start
			chStat <- procDbData(val, i)
			defer done.Done()
		}(divided[i], i)
	}

	at1 := time.Now()
	close(start)
	done.Wait()
	l.Printf(" Sqlasyn finished - %v...\r\n", time.Since(at1))

	close(chStat)
	sStat := stat{}

	for elem := range chStat {
		l.Printf(" Stat result %#v", elem)
		sStat.all += elem.all
		sStat.skiped += elem.skiped
		sStat.inserted += elem.inserted
		sStat.updated += elem.updated
	}

	l.Printf(" -> Total stat Get result from smd  ='%#v'", sStat)
	l.Printf("=== Successfully stop elapsed %dms\n \n \n", logger.TimeElapsed()/1000)

}

// procDbData process data for database
func procDbData(val []Smd, routineN int) *stat {
	stat := new(stat)

	for _, s := range val {
		d := s.smdToDbData()
		//  status // -1 error , 0 not changed , 1 new , 2 update
		stat.all++
		if status, sid := db.CheckModified(d.RequestID, d.DtModified); status == 0 {
			stat.skiped++
		} else if status == 1 {
			stat.inserted++
			row := db.Insert(d)
			l.Printf(" --> inserted effectes %d", row)
		} else if status == 2 {
			stat.updated++
			row := db.Update(sid, d)
			l.Printf(" --> updated sid %d effectes %d", sid, row)
		} else {
			l.Printf(" ! WARNING- status %#v sid %#v", status, sid)
		}
	}
	l.Printf("  --> Routine %d Get result from smd  ='%#v'", routineN, stat)
	return stat

}

// get - function for get data from web table
func get() *[]Smd {
	l.Println("1. Performing Http Get...")
	//resp, err := http.Get("https://smd.permkrai.ru/IPCP/HandlingReportPlugin/Api/analytics/ver.0.1/message")
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get("https://smd.permkrai.ru/IPCP/HandlingReportPlugin/Api/analytics/ver.0.2/requests/f1ae1eef-16ea-44cb-b77f-6b978ee4075d")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// Convert response body to string
	// bodyString := string(bodyBytes)
	// fmt.Println("API Response as String:\n" + bodyString)

	// Convert response body to Todo struct
	var todoStruct []Smd
	json.Unmarshal(bodyBytes, &todoStruct)
	l.Printf("API Response as struct %#v\n", len(todoStruct))
	// l.Printf("Record %#v\n", todoStruct[0])
	return &todoStruct

}

func (s *Smd) smdToDbData() *db.Smddb {
	var r = new(db.Smddb)
	id, _ := strconv.ParseInt(s.RequestID, 10, 64)

	// join questions
	q, err := json.Marshal(s.Questions)
	if err != nil {
		l.Printf(" !!! error json create s.Questions %s", err)
	}

	e, err := json.Marshal(s.ExceptionMessage)
	if err != nil {
		l.Printf(" !!! error json create s.ExceptionMessage %s", err)
	}

	intIsdirect := int64(0)
	if s.IsDirect {
		intIsdirect = 1
	}

	r.RequestID = id
	r.Number = s.Number
	r.DtModified = strtoTime(layoutDateTimeSmd, s.DtModified)
	r.DepartmentID = s.DepartmentID
	r.DepartmentName = s.DepartmentName
	r.Format = s.Format
	r.FormatName = s.FormatName
	r.IsDirect = intIsdirect
	r.CreateDate = strtoTime(layoutDtSmd, s.CreateDate)
	r.Name = s.Name
	r.Address = s.Address
	r.Email = s.Email
	r.ReceiveDate = strtoTime(layoutDateTimeSmd, s.ReceiveDate)
	r.DispatchDate = strtoTime(layoutDateTimeSmd, s.DispatchDate)
	r.UploadDate = strtoTime(layoutDateTimeSmd, s.UploadDate)
	r.ExceptionMessage = string(e)
	r.Questions = string(q)
	return r
}

func strtoTime(layout string, s string) time.Time {

	dt, err := time.Parse(layout, s)
	if err != nil {
		// l.Printf("strtoTime error %s", err)
		dt = time.Time{}
	}
	return dt
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tttest25/goreceptionloader/db"
	"github.com/tttest25/goreceptionloader/logger"
)

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

var (
	// Logger variable for logging
	l *log.Logger

// log nage

)

func get() {
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
	// return todoStruct[0]

}

func main() {

	// defer db.Close()
	l = logger.ReturnLogger("main")
	l.Printf("--- Start ")
	get() // get json data
	

	l.Printf("Elapsed send to nagios %dms \n", logger.TimeElapsed()/1000)
	l.Printf("=== Successfully stop\n \n \n")

}

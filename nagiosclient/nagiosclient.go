package nagiosclient

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/tttest25/gonagiosmetric/logger"
)

var (
	// Logger variable for logging
	l *log.Logger
	// Ns Nagios Struct
	Ns = NagiosOutputStruct{ // ns == &NagiosOutputStruct{"Status", 8}
		Status: 0,
	}
)

// NagiosOutputStruct - Example of return metric data - exit code see bellow
// DISK OK - free space: / 3326 MB (56%); | /=2643MB;5948;5958;0;5968
// / 15272 MB (77%);
// /home 69357 MB (27%);
// /var/log 819 MB (84%); | /boot=68MB;88;93;0;98
// /home=69357MB;253404;253409;0;253414
// /var/log=818MB;970;975;0;980
type NagiosOutputStruct struct {
	Host              string
	ServiceName       string // DISK - service name
	Serviceoutput     string // DISK(service name) OK(status - exitcode) - output  free space: / 3326 MB (56%);
	Serviceperfdata   string // 	/=2643MB;5948;5958;0;5968 /boot=68MB;88;93;0;98 /home=69357MB;253404;253409;0;253414 /var/log=818MB;970;975;0;980
	Longserviceoutput string // / 15272 MB (77%);\n/boot 68 MB (69%);\n/var/log 819 MB (84%);
	Status            int    // 0	OK / 1	Warning / 2	Critical / 3 Unknown
}

// NagiosAddService - add new service to metric
func (n *NagiosOutputStruct) NagiosAddService(label string, value float64, format string, UOM string, warn float64, crit float64, min float64, max float64) {
	// add typical service
	status := NagiosGetTresh(value, warn, crit)
	if format == "f" {
		n.NagiosAddParam(fmt.Sprintf("%s: %.3f %s", label, value, UOM), status)
		n.NagiosAddPerformance(NagiosPerfData(label, value, UOM, warn, crit, min, max))
	} else {
		n.NagiosAddParam(fmt.Sprintf("%s: %d %s", label, int(math.Round(value)), UOM), status)
		n.NagiosAddPerformance(NagiosPerfDataI(label, math.Round(value), UOM, warn, crit, min, max))
	}

}

// NagiosSetServiceName - SetServiceName
func (n *NagiosOutputStruct) NagiosSetServiceName(sn string) {
	n.ServiceName = sn
}

// SetHost - SetHost
func (n *NagiosOutputStruct) SetHost(sn string) {
	n.Host = sn
}

// NagiosOutput - from struct to nagios output
func (n *NagiosOutputStruct) NagiosOutput() string {
	// "OK - db: 23 424mb | ",
	return fmt.Sprintf("%s %s - %s | %s \n %s",
		n.ServiceName, NagiosStatus(n.Status), n.Serviceoutput, n.Serviceperfdata, n.Longserviceoutput,
	)
}

// NagiosPassive - return string for passive send in nagios
func (n *NagiosOutputStruct) NagiosPassive() string {
	// fed-serv;check_MISOJ_reception;0;
	return fmt.Sprintf("%s;%s;%d;",
		n.Host, n.ServiceName, n.Status,
	)
}

// NagiosAddParam - add new param to structs
func (n *NagiosOutputStruct) NagiosAddParam(str string, status int) {
	//
	pref := ""
	if status > n.Status {
		n.Status = status
	}
	if status > 1 {
		pref = "!!!"
	} else if status > 0 {
		pref = "!"
	}
	n.Serviceoutput = n.Serviceoutput + pref + str + " "
}

// AddUpdate - add update to service
func (n *NagiosOutputStruct) AddUpdate() {
	//
	currentTime := time.Now()
	n.Serviceoutput = n.Serviceoutput + " Update at " + currentTime.Format("2006-01-02 15:4:5")
}

// NagiosAddPerformance - add performance data
func (n *NagiosOutputStruct) NagiosAddPerformance(str string) {
	n.Serviceperfdata = n.Serviceperfdata + str
}

// 0.123123 -> "1.234" null -> ''
func nagiosFtoS(f float64) string {
	if f < 0 {
		return ""
	}
	return fmt.Sprintf("%.3f", f)
}

// 0.123123 -> "1.234" null -> ''
func nagiosFtoIS(f float64) string {
	if f < 0 {
		return ""
	}
	return fmt.Sprintf("%d", int64(math.Round(f)))
}

// NagiosPerfData - function to add perfdata
func NagiosPerfData(label string, value float64, UOM string, warn float64, crit float64, min float64, max float64) string {
	// 'label'=value[UOM];[warn];[crit];[min];[max]
	return fmt.Sprintf("'%s'=%s%s;%s;%s;%s;%s ", label, nagiosFtoS(value), UOM, nagiosFtoS(warn), nagiosFtoS(crit), nagiosFtoS(min), nagiosFtoS(max))
}

// NagiosPerfDataI - function to add perfdata integer
func NagiosPerfDataI(label string, value float64, UOM string, warn float64, crit float64, min float64, max float64) string {
	// 'label'=value[UOM];[warn];[crit];[min];[max]
	return fmt.Sprintf("'%s'=%s%s;%s;%s;%s;%s ", label, nagiosFtoIS(value), UOM, nagiosFtoIS(warn), nagiosFtoIS(crit), nagiosFtoIS(min), nagiosFtoIS(max))
}

// NagiosService - (str) function to add new measure for service example "free space: / 3326 MB (56%);"
func NagiosService(value string) string {
	// DISK OK - free space: / 3326 MB (56%);
	return fmt.Sprintf("%s; ", value)
}

// NagiosGetTresh Get status by treshold
func NagiosGetTresh(v float64, w float64, c float64) int {
	status := 3
	if v > c || v < 0 {
		status = 2
	} else if v > w {
		status = 1
	} else {
		status = 0
	}
	return status
}

// NagiosStatus - f(int)string get str description of status
func NagiosStatus(stat int) string {
	// 0	OK / 1	Warning / 2	Critical / 3 Unknown
	result := ""
	switch stat {
	case 0:
		result = "OK"
	case 1:
		result = "Warning"
	case 2:
		result = "Critical"
	case 3:
		result = "Unknown"
	default:
		result = fmt.Sprintf("Error of status(%d)", stat)
	}
	return result
}

// SendToNagios - send http post to passive nagios server
func SendToNagios(str string) {
	form := url.Values{
		"perfdata": {str},
	}
	form.Add("ln", "ln")
	form.Add("ip", "ip")
	form.Add("ua", "ua")

	req, err := http.NewRequest("POST", "http://10.59.20.16:8000", strings.NewReader(form.Encode()))
	if err != nil {
		l.Fatal(err)
	}
	req.SetBasicAuth("fed_monitor", "l3tm31n")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	cli := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := cli.Do(req)
	if err != nil {
		l.Fatal(err)
	}
	l.Printf("Resp %#v  \n", resp)
}

func init() {
	l = logger.ReturnLogger("nagios")

}

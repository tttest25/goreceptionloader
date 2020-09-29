package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tttest25/gonagiosmetric/logger"
	//"encoding/json"
)

//data-source="database" data-channel="2" data-userip="46.146.166.64" data-proto="https"

// Metric - type for saved data for metrics out
type Metric struct {
	Database    float64
	Queries     float64
	Application float64
	Total       float64
	Metrics     []float64
	Source      string
	Channel     int64
	Userip      string
	Proto       string
}

var (
	// Logger variable for logging
	l  *log.Logger
	pa *Metric // pa == nil

)

func (m *Metric) String() string {
	return fmt.Sprintf("Metric  %#v\n", m)
}

func stringToFloat(str string) float64 {
	f := strings.ReplaceAll(str, " s", "")
	s, err := strconv.ParseFloat(f, 32)
	if err == nil {
		return s // 3.1415927410125732
	}
	l.Printf("Error convert %s", err.Error())
	return -1
}

func stringToInt(str string) int64 {
	f := strings.ReplaceAll(str, " s", "")
	s, err := strconv.ParseInt(f, 10, 64)
	if err == nil {
		return s // 3.1415927410125732
	}
	l.Printf("Error convert %s", err.Error())
	return -1
}

// Scrape return data from http
func Scrape() *Metric {

	pa = new(Metric)

	l.Printf("Start scrapping")
	// Request the HTML page.
	res, err := http.Get("https://reception.gorodperm.ru/index.php?id=280")
	if err != nil {
		l.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		l.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		l.Fatal(err)
	}

	// strMetric := ""
	// Find the review items
	doc.Find("#stat").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//<div id="stat" data-database="0.0011 s" data-queries="7" data-application="0.0339 s" data-total="0.0349 s" data-source="database">stat</div>
		database, _ := s.Attr("data-database")
		queries, _ := s.Attr("data-queries")
		application, _ := s.Attr("data-application")
		total, _ := s.Attr("data-total")
		source, _ := s.Attr("data-source")
		channel, _ := s.Attr("data-channel")
		userip, _ := s.Attr("data-userip")
		proto, _ := s.Attr("data-proto")
		//l.Printf("Float 64 database: %f", stringToFloat(database))
		pa.Database = stringToFloat(database)
		pa.Queries = stringToFloat(queries)
		pa.Application = stringToFloat(application)
		pa.Total = stringToFloat(total)
		pa.Source = source
		pa.Channel = stringToInt(channel)
		pa.Userip = userip
		pa.Proto = proto
		// l.Printf("Float 64 database: %#v", pa)
		// strMetric = fmt.Sprintf("Review %d: MODx DB %s - Queries %s - App %s - Total %s\n", i, database, queries, application, total)

	})
	pa.Metrics = ScrapeMeasureChannels()
	l.Printf("get result")
	// return pa.String() //strMetric
	return pa
}

// ScrapeMeasureChannels - Parallel collect time of request for 2 channels
func ScrapeMeasureChannels() []float64 {
	start := time.Now()
	// var wg sync.WaitGroup
	var urls = []string{
		"https://reception1.gorodperm.ru/index.php?id=280",
		"https://reception2.gorodperm.ru/index.php?id=280",
	}

	s := []float64{}
	c := make(chan int64)
	l.Printf("  - Start measure 2 channels")

	// // Sequence start
	// for _, url := range urls {
	// 	go scrapeMeasureChannel(url, c)
	// 	l.Printf("   result ->  %#v", <-c)
	// }

	// l.Printf("   = sequent running %d values ->  %#v", logger.TimeTrack(start), s)

	// Parallel start
	start = time.Now()
	for _, url := range urls {
		go scrapeMeasureChannel(url, c)
	}
	s = append(s, float64(<-c)/1000000)
	s = append(s, float64(<-c)/1000000)
	// l.Printf("     result ->  %#v %d", <-c, <-c)
	l.Printf("   = parallel running %d  values ->  %#v", logger.TimeTrack(start), s)
	return s
}

func scrapeMeasureChannel(url string, c chan int64) {
	start := time.Now()
	l.Printf("       - Start measure %s", url)
	// Request the HTML page.
	client := http.Client{
		Timeout: 1000 * time.Millisecond,
	}

	res, err := client.Get(url)
	if err != nil {
		l.Println(err)
		c <- -1
		return
	}
	defer res.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		l.Println(err)
		c <- -1
		return
	}
	// flag := false
	sel := doc.Find("#stat")
	// for _ = range sel.Nodes {
	// 	flag = true
	// 	// use `single` as a selection of 1 node
	// }
	// //	not right page
	if sel.Nodes == nil {
		c <- -1
		return
	}

	l.Printf("status code error: %#v ", doc)

	if res.StatusCode != 200 {
		l.Printf("status code error: %d %s", res.StatusCode, res.Status)
		c <- -1
		return
	}
	c <- logger.TimeTrack(start)
	return
}

func init() {
	l = logger.ReturnLogger("scrapper")

}

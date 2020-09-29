package logger

import (
	"io"
	"log"
	"os"
	"time"
)

const (
	logFilename  string = "log.txt"
	sizeLogfifle int64  = 10 * 1024 * 1024           // n*x megabytes
	reduceLogife int64  = sizeLogfifle - 1*1024*1024 // reduce amount of bytes must be LESS OF reduceLogfile
)

var (
	// Logger variable for logging
	Logger *log.Logger
	start  time.Time
	file   *os.File
	err    error

	// log nage

)

// truncateFile copy file with skipping from head
func truncateFile(src string, lim int64) {
	fin, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	fout, err := os.Create(src + "tmp")
	if err != nil {
		panic(err)
	}
	defer fout.Close()

	// Offset is the number of bytes you want to exclude
	_, err = fin.Seek(int64(lim), io.SeekStart)
	if err != nil {
		panic(err)
	}

	n, err := io.Copy(fout, fin)
	Logger.Printf("Copied %d bytes, err: %v", n, err)

	if err := os.Remove(src); err != nil {
		panic(err)
	}

	if err := os.Rename(src+"tmp", src); err != nil {
		panic(err)
	}

}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// TimeElapsed - return ms elapsed time from starttime
func TimeElapsed() int64 {
	elapsed := time.Since(start)
	return elapsed.Nanoseconds() / 1000
}

// TimeTrack - return ms elapsed time from time
func TimeTrack(start time.Time) int64 {
	elapsed := time.Since(start)
	return elapsed.Nanoseconds() / 1000
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {

	if !fileExists(logFilename) {
		f, err := os.Create(logFilename)
		check(err)
		f.Close()
	}

	fi, err := os.Stat(logFilename)
	if err != nil {
		log.Fatal(err)
		check(err)
	}
	// get the size
	size := fi.Size()

	if fileExists(logFilename) && size > sizeLogfifle {
		truncateFile(logFilename, sizeLogfifle-reduceLogife)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}
	start = time.Now()

	// If the file doesn't exist, create it or append to the file
	file, err = os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Fatal(err)
	}
	Logger = log.New(file, "logger: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.SetOutput(file)

	Logger.Printf("Init")
}

func ReturnLogger(nameLogger string) *log.Logger {
	logger := log.New(file, nameLogger+": ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetOutput(file)
	return logger
}

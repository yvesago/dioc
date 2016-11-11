package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	Pannel        string
	ReconnectTime time.Duration = 180
	Debug         bool          = false
	Surveys                     = make(map[string]*regexp.Regexp)
	CRCa          string
)

func main() {
	filePtr := flag.String("f", "", "File to monitor")
	debugPtr := flag.Bool("d", false, "Debug mode")
	serverPtr := flag.String("s", "http://127.0.0.1:8080/client/api/v1", "Pannel url")
	flag.Parse()
	file := *filePtr
	Debug = *debugPtr
	Pannel = *serverPtr

	if Debug {
		ReconnectTime = 5
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		flag.PrintDefaults()
		os.Exit(0)
	}

	// Register
	type param struct {
		FileSurvey string
	}
	var p = param{FileSurvey: file}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(p)
	req, _ := http.NewRequest("POST", Pannel+"/agent", b)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	for { // wait for server
		rsp, err := client.Do(req)
		if err == nil {
			json.NewDecoder(rsp.Body).Decode(&CRCa)
			break
		} else {
			time.Sleep(ReconnectTime * time.Second)
		}
	}
	DebugLog("crca: " + CRCa)

	var wg sync.WaitGroup
	wg.Add(1)
	go Tail(file)

	time.Sleep(5 * time.Second)
	currentList := httpGETCommands(file)
	SyncList(currentList)
	wg.Add(1)
	go SyncCMD(file)
	wg.Wait()
}

func Tail(fname string) {
	// Monitor file
	t, _ := tail.TailFile(fname, tail.Config{Follow: true, ReOpen: true})
	for line := range t.Lines {
		for crcs, re := range Surveys {
			if re != nil {
				if re.MatchString(line.Text) {
					//fmt.Println(crcs + " " + line.Text)
					httpPOSTAlerte(crcs, line.Text)
				}
			}
		}
	}
}

func SyncCMD(filename string) {
	// Loop
	for {
		time.Sleep(ReconnectTime * time.Second)
		currentList := httpGETCommands(filename)
		SyncList(currentList)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func SyncList(liste []string) {
	for _, crcs := range liste {
		_, ok := Surveys[crcs]
		if ok == false {
			search := httpGETSurvey(crcs)
			if search != "" {
				re, err := regexp.Compile(search)
				if err == nil {
					Surveys[crcs] = re
				} else {
					Surveys[crcs] = nil
				}
			} else {
				delete(Surveys, crcs)
			}
			DebugLog(fmt.Sprintf(" Get Survey: %v", Surveys[crcs]))
		}
	}
	for crcs, _ := range Surveys {
		if contains(liste, crcs) == false {
			DebugLog(" Remove Survey: " + crcs)
			delete(Surveys, crcs)
		}
	}
}

func httpGETSurvey(crcs string) string {
	req, err := http.NewRequest("GET", Pannel+"/survey/"+crcs, nil)
	client := &http.Client{}
	//req.Header.Add("X-myToken", "ixxxx")
	rsp, err := client.Do(req)
	if err != nil {
		DebugLog("Bad connection to panel...")
		return ""
	} else {
		defer rsp.Body.Close()
		buf, _ := ioutil.ReadAll(rsp.Body)
		type ShortSurvey struct {
			CRCs   string
			Search string
			Id     int64
		}
		var s ShortSurvey
		json.Unmarshal(buf, &s)
		return s.Search
	}
	return ""
}

func httpGETCommands(filename string) []string {
	req, err := http.NewRequest("GET", Pannel+"/agent/"+CRCa, nil)
	client := &http.Client{}
	//req.Header.Add("X-myToken", "ixxxx")
	rsp, err := client.Do(req)
	if err != nil {
		DebugLog("Bad connection to panel...")
		return nil
	} else {
		defer rsp.Body.Close()
		buf, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			DebugLog("bad body")
			return nil
		} else {
			type CMDresp struct {
				CMD          string
				ListeSurveys []string
			}
			var cmd CMDresp
			json.Unmarshal(buf, &cmd)
			if cmd.CMD != "" {
				DebugLog(cmd.CMD)
			}
			if cmd.CMD == "STOP" {
				fmt.Println("CMD STOP")
				os.Exit(0)
			}
			if cmd.CMD == "FullSearch" {
				lines := FullSearch(filename)
				for _, l := range lines {
					httpPOSTAlerte(l.crcs, l.line)
				}
				cmd.CMD = "SendLines" // Update lines reset CMD
			}
			// if SendLines PUT lines
			if cmd.CMD == "SendLines" {
				lines := GetLines(filename)
				httpPUTlines(lines)
			}
			DebugLog(strings.Join(cmd.ListeSurveys, ", "))
			return cmd.ListeSurveys
		}

	}
	return nil
}

func httpPUTlines(lines []string) {
	b := new(bytes.Buffer)
	type Lines struct {
		Lines string
	}
	l := Lines{Lines: strings.Join(lines, "\n")}
	json.NewEncoder(b).Encode(l)
	DebugLog(fmt.Sprintf("%v", b))
	req, _ := http.NewRequest("PUT", Pannel+"/agent/"+CRCa, b)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	//req.Header.Add("X-myToken", "ixxxx")
	_, err := client.Do(req)
	if err != nil {
		DebugLog("Bad connection to panel...")
		return
	}
}

func httpPOSTAlerte(crcs string, line string) {
	type NewAlerte struct {
		CRCa string
		CRCs string
		Line string
	}
	b := new(bytes.Buffer)

	a := NewAlerte{CRCa: CRCa, CRCs: crcs, Line: line}
	json.NewEncoder(b).Encode(a)
	DebugLog(fmt.Sprintf("%v", b))
	req, _ := http.NewRequest("POST", Pannel+"/alerte", b)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	//req.Header.Add("X-myToken", "ixxxx")
	_, err := client.Do(req)
	if err != nil {
		DebugLog("Bad connection to panel...")
		return
	}
}

func FullSearch(filename string) []struct{ crcs, line string } {
	f, _ := os.Open(filename)
	defer f.Close()

	var linesFound []struct{ crcs, line string }
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tmpline := scanner.Text()
		for crcs, re := range Surveys {
			if re != nil {
				if re.MatchString(tmpline) {
					//httpPOSTAlerte(crcs, tmpline)
					linesFound = append(linesFound, struct{ crcs, line string }{crcs, tmpline})
				}
			}
		}
	}
	return linesFound
}

func GetLines(filename string) []string {
	var lines []string
	f, err := os.Open(filename)
	if err != nil {
		lines = append(lines, fmt.Sprintf("err: %v", err))
		return lines
	}
	defer f.Close()

	fi, _ := f.Stat()
	size := fi.Size()
	lines = append(lines, fmt.Sprintf("file: %s %s %dK", filename, fi.Mode(), int(size/1024)))
	scanner := bufio.NewScanner(f)
	var i int64 = 0
	var tmpline string = ""
	for scanner.Scan() {
		tmpline = scanner.Text()
		if i < 3 {
			lines = append(lines, tmpline)
		}
		if i == 3 {
			lines = append(lines, "...")
			//break
		}
		i++
	}
	lines = append(lines, tmpline) //last line
	if err := scanner.Err(); err != nil {
		lines = append(lines, fmt.Sprintf("err: %v", err))
	}
	return lines
}

func DebugLog(text string) {
	if !Debug {
		return
	}
	currenttime := time.Now().Local()
	fmt.Println("[", currenttime.Format("2006-01-02 15:04:05"), "] "+text)
}
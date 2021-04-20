package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/url"
	"nvim-grammarly/alert"
	"nvim-grammarly/pkg/transport"
	"os"
	"regexp"
	"strings"
)


// This is a test config that you can edit and drop directly into "config.Header ="
var testHeadersFromPython = map[string][]string{
	"Origin":           []string{"moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b"},
	"X-Client-Version": []string{"8.852.2307"},
	"Accept-Language":  []string{"en-GB,en-US;q=0.9,en;q=0.8"},
	"Accept-Encoding":  []string{"gzip, deflate, br"},
	"X-Client-Type":    []string{"extension-firefox"},
	"X-Container-Id":   []string{"aaukbtnoho4o302"},
	"User-Agent":       []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"},
	"Host":             []string{"capi.grammarly.com"},
	"Cookie":           []string{"firefox_freemium=true; grauth=AABJLPGf_MK7rGO8a-rHC3E3Vqhkx9EpFqyCoHSKx7WnGFjb7W1mPgg1_4WBkOlFNAgw8Usr8QODsbua; funnelType=free; redirect-location=eyJ0eXBlIjogIiIsICJsb2NhdGlvbiI6ICJodHRwczovL3d3dy5ncmFtbWFybHkuY29tL2FmdGVyX2luc3RhbGxfcGFnZT9leHRlbnNpb25faW5zdGFsbD10cnVlJnV0bV9tZWRpdW09c3RvcmUmdXRtX3NvdXJjZT1maXJlZm94In0=; csrf-token=AABJLAyQUhfM8XN5RxV8YFcvwSMPKGTu+T8gLw; gnar_containerId=aaukbtnoho4o302; browser_info=FIREFOX:67:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X; "},
	"Pragma":           []string{"no-cache"},
	"Cache-Control":    []string{"no-cache"},
}

var buildOTMsg = map[string]interface{}{
	"ch":     []string{"+0:0:Should catch the misspelled word.:0"},

}

type flags struct {
	configPath string
	filePath string
	logPath string
	logLevel int
}

func createBuildOTMsg(message string) map[string]interface{} {
	chMsg := fmt.Sprintf("+0:0:%s:0", message)
	return map[string]interface{}{
		"rev":    "0",
		"id":     "0",
		"action": "submit_ot",
		"ch" : []string{chMsg},
	}
}

var (
	configPath = flag.String("configPath",  "/Users/ronan/nvim-grammarly/sampleConfig.toml" ,
		"path to the directory containing the configuration file")
	filePath = flag.String("filePath", "/Users/ronan/tmp.txt",  "path to the file being checked")
	logPath = flag.String("logPath", "/Users/ronan/grammarly-go/grammarly-go.log",  "path to the file to print logs")
	logLevel = flag.Uint("logLevel", 4,  "log level")
)

func parseFlags() (*flags, error) {
	flag.Parse()
	if *filePath == ""{
		return nil, errors.New("no file passed, pass file using -filePath flag")
	}
	return &flags{
		configPath: *configPath,
		filePath:   *filePath,
		logPath: *logPath,
	}, nil
}

func setConfigDefaults(){
	viper.SetDefault("ContentDir", "content")
}


func readAlerts(c *websocket.Conn)[]*alert.Response{
	var alerts []*alert.Response
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("error receiving response; %v", err)
		}
		var r alert.Response
		if err := json.Unmarshal(message, &r); err != nil {
			log.Fatalf("error parsing response; %v", err)
		}
		if r.Action == "finished" {
			return alerts
		}
		if r.Action == "alert" {
			r.LineNum = 1
			alerts = append(alerts, &r)
		}
	}
}

func printAlerts(filePath string, alerts []*alert.Response){
	splitFilePath := strings.Split(filePath, "/")
	for _, alert := range alerts{

		exp := alert.GenExplanation()
		if exp == ""{
			continue
		}
		msg := fmt.Sprintf("%s:%d:%d: %s\n", splitFilePath[len(splitFilePath)-1], alert.LineNum, alert.HighlightBegin, exp)
		fmt.Print(msg)
	}
}
// this is a suboptimal way of getting the line number of the alert. Grammarly returns the character
// number of the start of the word so we have to calculate which line that is on. There is very likely
// a better way but this will do for now
func processLineNum(filepath string, alerts []*alert.Response)[]*alert.Response{
	var retAlert []*alert.Response
	file, err := os.Open(filepath)
	if err != nil {
		glog.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		for  i, alert := range alerts{
			if alert == nil{
				continue
			}
			lChars := uint16(len(l))
			if alert.HighlightBegin <= lChars{
				alerts[i] = nil
				retAlert = append(retAlert, alert)
				if len(retAlert) == len(alerts){
					return retAlert
				}
				continue
			}
			alert.HighlightBegin = alert.HighlightBegin - lChars
			alert.HighlightEnd = alert.HighlightEnd - lChars
			alert.LineNum = alert.LineNum + 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return retAlert
}

func getLogLevel(level int)log.Level{
	switch level {
	case 0:
		return log.PanicLevel
	case 1:
		return log.FatalLevel
	case 2:
		return log.ErrorLevel
	case 3:
		return log.WarnLevel
	case 4:
		return log.InfoLevel
	case 5:
		return log.DebugLevel
	default:
		return log.TraceLevel
	}
}

func main() {
	flags, err := parseFlags()
	if err != nil{
		os.Exit(1)
	}

	// read configuration from config file if present
	conf, err := transport.LoadConfig(flags.configPath)
	if err != nil{
		fmt.Printf("unexpected error reading config file, got %s\n", err)
	}

	// open or create log file with default log level being info level (L4)
	logFile, err := os.OpenFile(flags.logPath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0755)
	if err != nil {
		fmt.Println("error opening logfile")
	}
	log.RegisterExitHandler(func() {
		if logFile == nil {
			return
		}
		logFile.Close()
	})
	defer log.Exit(0)
	log.SetOutput(logFile)
	logLevel := getLogLevel(flags.logLevel)
	log.SetLevel(logLevel)

	content, err := ioutil.ReadFile(flags.filePath)
	if err != nil {
		log.Fatalf("error reading file; %v", err)
	}
	contentStr := string(content)
	// Replace new lines with spaces...
	re := regexp.MustCompile(`\r?\n`)
	contentStr = re.ReplaceAllString(contentStr, " ")
	completeBuildOTMsg := createBuildOTMsg(contentStr)

	u := url.URL{
		Scheme: "wss",
		Host:   "capi.grammarly.com",
		Path:   "/freews",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), testHeadersFromPython)
	if err != nil {
		log.Fatalf("error dialing websocket; %v", err)
	}

	initalMsg := conf.GenerateConfigMessage()
	if err = c.WriteJSON(initalMsg); err != nil {
		log.Fatalf("error sending request; %v", err)
	}
	if err = c.WriteJSON(completeBuildOTMsg); err != nil {
		log.Fatalf("error sending request; %v", err)
	}

	alerts := readAlerts(c)
	alerts = processLineNum(*filePath, alerts)
	printAlerts(*filePath, alerts)
}



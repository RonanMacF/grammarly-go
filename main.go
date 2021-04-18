package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"

	"github.com/gorilla/websocket"
)

const(
	configFiletype = "toml"
	configFileName = "config"
)
type config struct{}

var buildInitialMsg = map[string]interface{}{
	"type":            "initial",
	"docid":           "1234",
	"client":          "extension_chrome",
	"protocolVersion": "1.0",
	"clientSupports": []string{
		"free_clarity_alerts",
		"readability_check",
		"filler_words_check",
		"sentence_variety_check",
		"free_occasional_premium_alerts",
	},
	"dialect":       "british",
	"clientVersion": "14.924.2437",
	"extDomain":     "keep.google.com",
	"action":        "start",
	"id":            "0",
	"sid":           "0",
}

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
	"rev":    '0',
	"id":     '0',
	"action": "submit_ot",
}

type flags struct {
	configPath string
	filePath string
}

func createBuildOTMsg(message string) map[string]interface{} {
	chMsg := fmt.Sprintf("+0:0:%s:0", message)
	completeBuildOTMsg := buildOTMsg
	fmt.Printf("%+v\n", buildOTMsg)
	completeBuildOTMsg["ch"] = []string{chMsg}
	return completeBuildOTMsg
}

var configPath = flag.String("configPath",  "$HOME/.config/grammarly-go" ,
"path to the directory containing the configuration file")
var filePath = flag.String("filePath", "",  "path to the file being checked")

func parseFlags() (*flags, error) {
	flag.Parse()
	if *filePath == ""{
		return nil, errors.New("no file passed, pass file using -filePath flag")
	}
	return &flags{
		configPath: *configPath,
		filePath:   *filePath,
	}, nil
}

func setConfigDefaults(){
	viper.SetDefault("ContentDir", "content")
}

func LoadConfig(path string) (*config, error) {
	setConfigDefaults()
	viper.AddConfigPath(path)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			glog.V(1).Infof("no configuration file called %s.%s found in %s", configFileName,
				configFiletype, path)
			return &config{}, nil
		}
		return &config{}, err
	}
	return &config{}, nil
}

func main() {
	flags, err := parseFlags()
	if err != nil{
		glog.Infof("error received parsing flags, expected nop error, got %s", err)
		os.Exit(1)
	}

	config, _ := LoadConfig(flags.configPath)
	viper.New()
	fmt.Printf("%+v\n", config)
	content, err := ioutil.ReadFile(flags.filePath)
	if err != nil {
		log.Fatalf("error reading file; %v", err)
	}
	contentStr := string(content)
	// Replace new lines with spaces...
	re := regexp.MustCompile(`\r?\n`)
	contentStr = re.ReplaceAllString(contentStr, " ")
	completeBuildOTMsg := createBuildOTMsg(string(contentStr))

	u := url.URL{
		Scheme: "wss",
		Host:   "capi.grammarly.com",
		Path:   "/freews",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), testHeadersFromPython)
	if err != nil {
		log.Fatalf("error dialing websocket; %v", err)
	}
	if err = c.WriteJSON(buildInitialMsg); err != nil {
		log.Fatalf("error sending request; %v", err)
	}
	if err = c.WriteJSON(completeBuildOTMsg); err != nil {
		log.Fatalf("error sending request; %v", err)
	}

	var alerts []Response
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("error receiving response; %v", err)
		}
		tmp := fmt.Sprintf("%s", message)
		fmt.Println(tmp)
		var r Response
		if err := json.Unmarshal(message, &r); err != nil {
			log.Fatalf("error parsing response; %v", err)
		}
		if r.Action == "finished" {
			break
		}
		if r.Action == "alert" {
			fmt.Printf("%+v\n", r)
			alerts = append(alerts, r)
		}
	}
	processLineNum("/Users/ronan/nvim-grammarly/main.go", nil)
}

// this is a suboptimal way of getting the line number of the alert. Grammarly returns the character
// number of the start of the word so we have to calculate which line that is on. There is very likely
// a better way but this will do for now
func processLineNum(filepath string, alerts []Response)[]Response{
	var retAlert []Response
	file, err := os.Open(filepath)
	if err != nil {
		glog.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		for  i, alert := range alerts{
			lChars := uint16(len(l))
			if alert.HighlightBegin <= lChars{
				alerts = append(alerts[:i], alerts[:i+1]...)
				retAlert = append(retAlert, alert)
				if len(alerts) == 0{
					return retAlert
				}
			}
			alert.HighlightBegin = alert.HighlightBegin - lChars
			alert.HighlightEnd = alert.HighlightEnd - lChars
			alert.lineNum = alert.lineNum + 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return retAlert
}
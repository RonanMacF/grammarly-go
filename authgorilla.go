package main

import (
	"os"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"io/ioutil"
	"flag"
	"regexp"
)

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
	fileName string
}

func createBuildOTMsg(message string) map[string]interface{} {
	chMsg := fmt.Sprintf("+0:0:%s:0", message)
	completeBuildOTMsg := buildOTMsg
	fmt.Printf("%+v\n", buildOTMsg)
	completeBuildOTMsg["ch"] = []string{chMsg}
	return completeBuildOTMsg
}

func main(){
	flags := parseFlags()
	fmt.Printf("%+v\n", flags)
	content, err := ioutil.ReadFile(flags.fileName)
	if err != nil {
		log.Fatal("error reading file; %v", err)
	}
	contentStr := string(content)
	// Replace new lines with spaces...
	re := regexp.MustCompile(`\r?\n`)
	contentStr = re.ReplaceAllString(contentStr, " ")
	completeBuildOTMsg := createBuildOTMsg(string(contentStr))

	u := url.URL{
		Scheme: "wss",
		Host: "capi.grammarly.com",
		Path: "/freews",
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
	for  {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatalf("error receiving response; %v", err)
		}
		fmt.Printf("%s\n", message)
		output := map[string]interface{}{}
		if err := json.Unmarshal([]byte(message), &output); err != nil {
			log.Fatalf("error parsing response; %v", err)
		}
		if output["action"] == "finished" {
			os.Exit(0)
		}
	}
}

func parseFlags() *flags {
	f := flags{}
	flag.StringVar(&f.fileName, "fileName", "test.txt", "Name of the file being checked")
	flag.Parse()
	return &f
}


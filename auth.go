package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

var reqStdHeader = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0",
	"Accept":          "*/*",
	"Accept-Language": "en-US,en;q=0.5",
	"Accept-Encoding": "gzip, deflate, br",
	"Origin":          "moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b",
	//"Host":             "auth.grammarly.com",
	//"Connection":       "keep-alive, Upgrade",
	"Pragma":           "no-cache",
	"Cache-Control":    "no-cache",
	"X-Client-Version": "8.852.2307",
	"X-Client-Type":    "extension-firefox",
}

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

var buildOTMsg = map[string]interface{}{
	"ch":     []string{"+0:0:Should catch the misspelled word.:0"},
	"rev":    '0',
	"id":     '0',
	"action": "submit_ot",
}

var furtherHeaders = map[string]string{
	"firefox_freemium": "true",
	"funnelType":       "free",
	"browser_info":     "FIREFOX:67:COMPUTER:SUPPORTED:FREEMIUM:MAC_OS_X:MAC_OS_X",
}

const (
	stdHost = "auth.grammarly.com"
	authURL = "https://auth.grammarly.com/v3/user/oranonymous?app=firefoxExt&containerId=aaukbtnoho4o302"
)

type flags struct {
	fileName string
}

func extractAuthCookies(cookies []*http.Cookie) (map[string]string, error) {
	authCookies := make(map[string]string)
	authCookies["gnar_containerId"] = "aaukbtnoho4o302"
	for _, cookie := range cookies {
		authCookies[cookie.Name] = cookie.Value
	}
	redirectLocation := map[string]string{
		"type":     "",
		"location": "https://www.grammarly.com/after_install_page?extension_install=true&utm_medium=store&utm_source=firefox",
	}
	jsonRedirect, err := json.Marshal(redirectLocation)
	if err != nil {
		return nil, fmt.Errorf("unable to convert %+v to json: %v", redirectLocation, err)
	}
	redirectEnc := base64.StdEncoding.EncodeToString([]byte(jsonRedirect))
	authCookies["redirect-location"] = redirectEnc
	// Below is value python produces...
	//authCookies["redirect-location"] = "eyJ0eXBlIjogIiIsICJsb2NhdGlvbiI6ICJodHRwczovL3d3dy5ncmFtbWFybHkuY29tL2FmdGVyX2luc3RhbGxfcGFnZT9leHRlbnNpb25faW5zdGFsbD10cnVlJnV0bV9tZWRpdW09c3RvcmUmdXRtX3NvdXJjZT1maXJlZm94In0="
	return authCookies, nil
}

func genPlagHeaders(cookie, containerID string) map[string][]string {
	reqStdHeader["Host"] = "capi.grammarly.com"
	reqStdHeader["Accept-Language"] = "en-GB,en-US;q=0.9,en;q=0.8"
	reqStdHeader["Cookie"] = cookie
	reqStdHeader["X-Container-Id"] = containerID
	reqStdHeader["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
	delete(reqStdHeader, "Accept")
	plagHeadersConf := make(map[string][]string)
	for name, value := range reqStdHeader {
		plagHeadersConf[name] = []string{value}
	}
	return plagHeadersConf
}

func mapToString(input map[string]string) string {
	output := new(bytes.Buffer)
	for key, value := range input {
		fmt.Fprintf(output, "%s=%s; ", key, value)
	}
	return output.String()
}

type T struct {
	Msg   string
	Count int
}

func main() {
	//flags := parseFlags()
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		log.Fatalf("unable to create new request; %v", err)
	}
	// Update request headers
	for name, value := range reqStdHeader {
		req.Header.Set(name, value)
	}
	req.Host = stdHost
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error sending get request; %v", err)
	}
	fmt.Printf("%+v\n", resp)
	fmt.Println("________________________________")
	cookies := resp.Cookies()
	fmt.Printf("%+v\n", cookies)
	extCookies, err := extractAuthCookies(cookies)
	if err != nil {
		log.Fatalf("error extracting auth cookies; %v", err)
	}
	// Add further headers
	for name, value := range furtherHeaders {
		extCookies[name] = value
	}
	/*fmt.Println("________________________________")
	fmt.Printf("%+v\n", extCookies)
	fmt.Println("________________________________")
	fmt.Println(mapToString(extCookies))
	fmt.Println("________________________________")*/
	plagHeaders := genPlagHeaders(mapToString(extCookies), extCookies["gnar_containerId"])
	config, err := websocket.NewConfig("wss://capi.grammarly.com/freews", "moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b")
	if err != nil {
		log.Fatalf("error creating websocket config; %v", err)
	}
	fmt.Println("________________________________")
	fmt.Printf("%+v\n", config)
	fmt.Println("________________________________")

	config.Header = plagHeaders
	fmt.Println("CONFIG")
	fmt.Printf("%+v\n", config)
	fmt.Println("________________________________")

	conn, err := websocket.DialConfig(config)
	if err != nil {
		log.Fatalf("error dialing websocket; %v", err)
	}
	defer conn.Close()
	if err = websocket.JSON.Send(conn, buildInitialMsg); err != nil {
		log.Fatalf("error sending request; %v", err)
	}
	if err = websocket.JSON.Send(conn, buildOTMsg); err != nil {
		log.Fatalf("error sending request; %v", err)
	}
	var message T
	if err = websocket.JSON.Receive(conn, &message); err != nil {
		log.Fatalf("error receiving response; %v", err)
	}
	fmt.Printf("%+v\n", message)
}

func parseFlags() *flags {
	f := flags{}
	flag.StringVar(&f.fileName, "fileName", "test.py", "Name of the file being checked")
	flag.Parse()
	return &f
}

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"golang.org/x/net/websocket"
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
	return authCookies, nil
}

func genPlagHeaders(cookie, containerID string) map[string]string {
	updatedHrds := reqStdHeader
	updatedHrds["Cookie"] = cookie
	updatedHrds["X-Container-Id"] = containerID
	updatedHrds["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36"
        delete(updatedHrds, "Accept")
	return updatedHrds
}

func mapToString(input map[string]string) string {
	output := new(bytes.Buffer)
	for key, value := range input {
		fmt.Fprintf(output, "%s=%s; ", key, value)
	}
	return output.String()
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
	fmt.Println("________________________________")
	fmt.Printf("%+v\n", extCookies)
	fmt.Println("________________________________")
        fmt.Println(mapToString(extCookies))
	fmt.Println("________________________________")
        plagHeaders := genPlagHeaders(mapToString(extCookies), extCookies["gnar_containerId"])
	conn, err := websocket.Dial("wss://capi.grammarly.com/freews", "", "moz-extension://6adb0179-68f0-aa4f-8666-ae91f500210b")
}

func parseFlags() *flags {
	f := flags{}
	flag.StringVar(&f.fileName, "fileName", "test.py", "Name of the file being checked")
	flag.Parse()
	return &f
}

package main

// Response is
type Response struct {
	Action         string `json:"action"`
	HighlightBegin uint16 `json:"highlightBegin"`
	HighlightEnd   uint16 `json:"highlightEnd"`
	lineNum uint16
	Group          string `json:"group"`
	Replacements   []string `json:"replacements"`
	Text           string `json:"text"`
}

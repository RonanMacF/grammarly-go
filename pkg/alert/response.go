package alert

// Response is
type Response struct {
	Action         string `json:"action"`
	HighlightBegin uint16 `json:"highlightBegin"`
	HighlightEnd   uint16 `json:"highlightEnd"`
	LineNum uint16
	TransformJson map[string]interface{} `json:"transformJson"`
	Explanation string `json:"explanation"`
	Category  string `json:"category"`
	Pname string`json:"pname"`
	PnameQualifier string`json:"pnameQualifier"`
	Point string`json:"point"`
	Group          string `json:"group"`
	Replacements   []string `json:"replacements"`
	Text           string `json:"text"`
	Title string `json:"title"`
}

func (r *Response) GenExplanation()  string{
	switch r.Category{
	case "Misspelled": return genMisspelledExplanation(r.Group, r.Point, r.Replacements, r.Explanation)
	case "Fragment": return genFragmentExplanation(r.PnameQualifier)
	case "AccidentallyConfused": return genConfusedExplanation(r.Title, r.Replacements)
	case "CommonlyConfused": return genConfusedExplanation(r.Title, r.Replacements)
	}
	return ""
}
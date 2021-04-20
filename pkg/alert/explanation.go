package alert

import (
	"fmt"
	"regexp"
	"strings"
)

const missingPronoun = "Missing Pronoun: It appears this sentence is incomplete"
const sentenceFragment = "Sentence Fragment: Consider rewriting as a complete sentence"

var (
unicodeRegexp = regexp.MustCompile("\\[u][0-9]{1-3}[0-9a-zA-Z]{0-3}")
	paraRegex = regexp.MustCompile("</?p>")
	textRegex = regexp.MustCompile("</?[a-z]>")
)

func genMisspelledExplanation(group string, point string, replacements []string, explanation string) string{
	if group == "Enhancement"{
		return ""
	}
	if point == "WordRepeat"{
		msg := normalizeString(explanation)
		return fmt.Sprintf("Repeated word: %s", msg)
	}
	return fmt.Sprintf("Spelling error: potentially replace with %s"	, normalizeSlice(replacements))
}

func genFragmentExplanation(pnameQualifier string) string{
	switch pnameQualifier {
	case "StMissingSubject": return missingPronoun
	default:return sentenceFragment
	}
}

func genConfusedExplanation(title string, replacements []string)string{
	replacementsStr := normalizeSlice(replacements)
	title = normalizeString(title)
	return fmt.Sprintf("%s: replace with %s", title, replacementsStr)
}

func normalizeSlice(strs []string)string{
	var sb strings.Builder
	for _, s := range strs{
		s = normalizeString(s)
		sb.WriteString(fmt.Sprintf("%s, ", s))
	}
	out := sb.String()
	if out == ""{
		return out
	}
	return out[:len(out)-2]
}
func normalizeString(s string)string{
	s = paraRegex.ReplaceAllString(s, "")
	s = textRegex.ReplaceAllString(s, "")
	s = unicodeRegexp.ReplaceAllString(s, "")
	s = strings.TrimSuffix(s, "\n")
	return s
}

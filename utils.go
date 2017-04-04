package mdrest

import (
	"time"
	"fmt"
	"strings"
	"path/filepath"
	"bytes"
	"sync"
	"unicode/utf8"
	"unicode"
)

func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02 15:04:05 -07:00",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05", // iso8601 without timezone
		"2006-01-02 15:04:05Z07:00",
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05Z07:00",
		"02 Jan 06 15:04 MST",
		"02 Jan 2006",

	})
}



func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.Parse(dateType, s); e == nil {
			return
		}
	}
	return d, fmt.Errorf("Unable to parse date: %s", s)
}


func AbsPath(basePath, currentLocation, path string) string {
	if strings.HasPrefix(path, "#")  || strings.HasPrefix(path, "/") {
		return path
	}
	wd := filepath.Dir(currentLocation)
	result := filepath.Join(wd, path)
	return basePath + result
}



// StripSummary truncates words by runes.
func StripSummary(text string,max int) (summary string, truncated bool) {
	words := strings.Fields(text)
	count := 0
	for index, word := range words {
		if count >= max {
			return strings.Join(words[:index], " "), true
		}
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			count++
		} else if count+runeCount < max {
			count += runeCount
		} else {
			for ri := range word {
				if count >= max {
					truncatedWords := append(words[:index], word[:ri])
					return strings.Join(truncatedWords, " "), true
				}
				count++
			}
		}
	}
	return strings.Join(words, " "), false
}


var stripHTMLReplacer = strings.NewReplacer("\n", " ", "</p>", "\n", "<br>", "\n", "<br />", "\n")

// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	s = stripHTMLReplacer.Replace(s)
	// Walk through the string removing all tags
	var bufferPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
	b := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		b.Reset()
		bufferPool.Put(b)
	}()
	var inTag, isSpace, wasSpace bool
	for _, r := range s {
		if !inTag {
			isSpace = false
		}
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case unicode.IsSpace(r):
			isSpace = true
			fallthrough
		default:
			if !inTag && (!isSpace || (isSpace && !wasSpace)) {
				b.WriteRune(r)
			}
		}
		wasSpace = isSpace

	}
	return  b.String()
}



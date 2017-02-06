package mdrest

import (
	"time"
	"fmt"
	"strings"
	"path/filepath"
	"bufio"
	"bytes"
	"sync"
	"unicode/utf8"
	"unicode"
)

const PrefixInternalMarkDown = "/MDREST_INTERNAL_MARKDOWN/"

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


//cleanMdContent covert relative path to abs path
func cleanMdContent(reader *bufio.Reader, basePath, location string) []byte {
	lineBreak := byte('\n')
	codeBoundary := []byte("```")
	var mdBytes []byte
	var inCode bool
	for {
		line,_, lineErr := reader.ReadLine()
		if lineErr != nil {
			break
		}
		if bytes.HasPrefix(line, codeBoundary) {
			inCode = !(inCode)
		}
		if !inCode {
			//[a link ](path/to/your/file "Logo Title Text 1")
			//![a image](path/to/your/img.png "Logo Title Text 1")
			//[logo]: path/to/your/img.png
			if bytes.HasPrefix(line, []byte("![")) || bytes.HasPrefix(line, []byte("[")) {
				line = cleanLine(line, basePath,location )
			}
		}
		line = append(line,lineBreak)
		mdBytes  = append(mdBytes, line...)
	}
	var contentBuffer bytes.Buffer
	reader.WriteTo(&contentBuffer)
	mdBytes  = append(mdBytes, contentBuffer.Bytes()...)
	return mdBytes
}



func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.Parse(dateType, s); e == nil {
			return
		}
	}
	return d, fmt.Errorf("Unable to parse date: %s", s)
}


/**
lines := [][]byte {
		[]byte(`[a link ](../file)`),
		[]byte(`[a link ](path/to/your/file "Logo Title Text 1")`),
		[]byte(`![a image](path/to/your/img.png "Logo Title Text 1")`),
		[]byte(`[logo]: path/to/your/img.png`),

	}

**/
func cleanLine(line []byte, basePath, location string) []byte {
	leftBoundary := byte(']')
	rightBoundary := byte(')')
	quotes := byte('"')
	leftIdx := -1
	rightIdx := -1
	quotesIdx := -1
	nopath := []byte("://")
	leftSign1 := []byte("](")
	leftSign2 := []byte("]:")
	nopathBoundary := nopath[0]
	for i, v := range line {
		switch v {
		case nopathBoundary:{
			if bytes.Equal(line[i:i+3], nopath) {
				return line
			}
		}
		//only //[a link ](.. and ]:...
		case leftBoundary:
			leftSign := line[i:i+2]
			if !(bytes.Equal(leftSign, leftSign1) || bytes.Equal(leftSign, leftSign2)){
				return line
			}
			if leftIdx > 0 {
				continue
			}
			leftIdx = i
		case quotes:
			if quotesIdx > 0 {
				continue
			}
			quotesIdx = i
		case rightBoundary:
			if rightIdx > 0 {
				continue
			}
			rightIdx = i
		}
	}
	if leftIdx < 2 {
		return line
	}
	var t []byte
	if rightIdx > 0 {
		if quotesIdx > 0 {
			t = line[leftIdx+2:quotesIdx-1]
		} else {
			t = line[leftIdx+2:rightIdx]
		}
	} else {
		t = bytes.TrimPrefix(line[leftIdx+2:],[]byte(" "))
	}
	absPath := []byte(AbsPath(basePath,location,string(t)))
	line = bytes.Replace(line,t, absPath,1)
	return line
}


//AbsPath absolute path
//currentLocation is where the file is , exp /your/path/web/index.html

//for some intreanl md file

//exp: your/path/x.md append /key_internal_markdown/, so in html
//you can replace <a href="/key_internal_markdown/ <a class="interanl_md" href="#

func AbsPath(basePath, currentLocation, path string) string {
	if strings.HasPrefix(path, "#")  || strings.HasPrefix(path, "/") {
		return path
	}
	wd := filepath.Dir(currentLocation)
	result := filepath.Join(wd, path)
	if strings.HasSuffix(path, ".md") {
		return PrefixInternalMarkDown + strings.Replace(result[:len(result)-3]," ","%20",-1)
	}
	return basePath + result
}



func StripSummary(text string, isCJKLanguage bool, summaryLength int) (summary string, truncated bool) {
	if isCJKLanguage {
		summary, truncated = TruncateWordsByRune(strings.Fields(text), summaryLength)
	} else {
		summary, truncated = TruncateWordsToWholeSentence(text, summaryLength)
	}
	return

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

// TruncateWordsByRune truncates words by runes.
func TruncateWordsByRune(words []string, max int) (string, bool) {
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

// TruncateWordsToWholeSentence takes content and truncates to whole sentence
// limited by max number of words. It also returns whether it is truncated.
func TruncateWordsToWholeSentence(s string, max int) (string, bool) {

	var (
		wordCount     = 0
		lastWordIndex = -1
	)

	for i, r := range s {
		if unicode.IsSpace(r) {
			wordCount++
			lastWordIndex = i

			if wordCount >= max {
				break
			}

		}
	}

	if lastWordIndex == -1 {
		return s, false
	}

	endIndex := -1

	for j, r := range s[lastWordIndex:] {
		if isEndOfSentence(r) {
			endIndex = j + lastWordIndex + utf8.RuneLen(r)
			break
		}
	}

	if endIndex == -1 {
		return s, false
	}

	return strings.TrimSpace(s[:endIndex]), endIndex < len(s)
}

func isEndOfSentence(r rune) bool {
	return r == '.' || r == '?' || r == '!' || r == '"' || r == '\n'
}




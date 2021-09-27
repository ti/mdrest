package mdrest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jyaml "github.com/ghodss/yaml"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	KeyTitle      = "title"
	KeyDate       = "date"
	KeyHtml       = "html"
	KeyLocation   = "location"
	KeyText       = "text"
	KeySummary    = "summary"
	KeyRawContent = "raw_content"
	KeyPicture    = "picture"
)

type Article map[string]interface{}

// Articles is an convenience type alias for article slice
type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//Less
func (a Articles) Less(i, j int) bool {
	//the Articles must have date key, if can not read from yaml header, the 'date' be replaced by file last modify time
	left := (*a[i])[KeyDate].(time.Time)
	right := (*a[j])[KeyDate].(time.Time)
	return right.Before(left)
}

//Get By location
func (a Articles) Get(location string) *Article {
	for _, v := range a {
		if (*v)[KeyLocation].(string) == location {
			return v
		}
	}
	return nil
}

func (a Articles) Remove(location string) *Article {
	for i, v := range a {
		if strings.HasPrefix((*v)[KeyLocation].(string), location) {
			a[i] = a[len(a)-1]
			a = a[:len(a)-1]
		}
	}
	return nil
}

// ReadArticle returns an article read from a Reader
func ReadArticle(srcDir, fpath, basePath string) (Article, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	//check first line for if it have fontmatter
	firstLine, lineErr := reader.ReadBytes(byte('\n'))
	var haveFrontMatter bool
	var body = make(map[string]interface{})
	if lineErr == nil && bytes.HasPrefix(firstLine, []byte("---")) {
		haveFrontMatter = true
		if front, err := parseFrontMatter(reader); err == nil {
			body = front
		}
	} else {
		body = make(map[string]interface{})
	}
	fileInfo, _ := file.Stat()
	if date, ok := body[KeyDate]; ok {
		if t, err := StringToDate(fmt.Sprint(date)); err == nil {
			body[KeyDate] = t
		} else {
			body[KeyDate] = fileInfo.ModTime()
		}
	} else {
		body[KeyDate] = fileInfo.ModTime()
	}

	if pic, ok := body[KeyPicture]; ok {
		if picLink, ok := pic.(string); ok {
			if !(strings.HasPrefix(picLink, "http://") || strings.HasPrefix(picLink, "https://")) {
				if strings.HasPrefix(picLink, "/") {
					body[KeyPicture] = basePath + picLink[1:]
				} else {
					picturePath := strings.TrimPrefix(AbsPath("", fpath, picLink), srcDir)
					body[KeyPicture] = basePath + picturePath
				}
			}
		}
	}
	if _, ok := body[KeyTitle]; !ok {
		var title string
		if bytes.HasPrefix(firstLine, []byte("# ")) {
			title = string(firstLine)[2:]
			firstLine = nil
		} else {
			bodyTitle, content := parseBodyTitle(reader)
			if bodyTitle != "" {
				title = string(bodyTitle)
			} else {
				firstLine = content
				haveFrontMatter = false
				title = strings.TrimSuffix(path.Base(fpath), path.Ext(fpath))
			}
		}
		title = strings.TrimSuffix(title, "\n")
		body[KeyTitle] = title
	}

	location := strings.TrimSuffix(strings.TrimPrefix(fpath, srcDir), path.Ext(fpath))
	//fix location case
	body[KeyLocation] = location
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		panic("READ CONTENT ERROR" + err.Error())
	}
	if haveFrontMatter {
		body[KeyRawContent] = content
	} else {
		body[KeyRawContent] = append(firstLine, content...)
	}
	return body, nil
}

func parseBodyTitle(reader *bufio.Reader) (title string, content []byte) {
	var maxTryLines = 5
	lineBreak := byte('\n')
	boundary := []byte("# ")
	boundaryLen := len(boundary)
	for i := 0; i < maxTryLines; i++ {
		line, lineErr := reader.ReadBytes(lineBreak)
		if lineErr != nil {
			break
		}
		if (len(line) <= boundaryLen) {
			continue
		}
		if bytes.HasPrefix(line, boundary) {
			title := line[boundaryLen:]
			return string(title), nil
		} else {
			return "", line
		}
	}
	return "", nil
}

// ParseFrontMatter reads the front matter-type article header
func parseFrontMatter(reader *bufio.Reader) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	lineBreak := byte('\n')
	boundary := []byte("---")
	var frontMatterBytes []byte
	for {
		line, lineErr := reader.ReadBytes(lineBreak)
		if lineErr != nil {
			break
		}
		if bytes.HasPrefix(line, boundary) {
			break
		}
		frontMatterBytes = append(frontMatterBytes, line...)
	}
	//convert yaml to jsonable map struct, for some advaced users
	jsb, err := jyaml.YAMLToJSON(frontMatterBytes)
	if err != nil {
		return nil, errors.New("can not covert front matter header to json")
	}
	err = json.Unmarshal(jsb, &data)
	if err != nil {
		return nil, errors.New("can not unmarshal front matter header to map object")
	}
	if len(data) == 0 {
		return nil, errors.New("empty matter")
	}
	return data, nil
}

// ReadArticle returns an article read from a Reader
func ReadArticles(srcDir, basePath string, showPageTitle bool) (articles Articles, err error) {
	//read files
	sourceFiles, err := ReadFiles(srcDir)
	if err != nil {
		return
	}
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	var rendererParameters blackfriday.HtmlRendererParameters

	htmlPrefix := "/" //可能后面会用到
	htmlPrefix = strings.TrimSuffix(htmlPrefix, "/")
	rendererParameters.AbsolutePrefix = htmlPrefix

	extensions := 0
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS
	extensions |= blackfriday.EXTENSION_HEADER_IDS
	extensions |= blackfriday.EXTENSION_AUTO_HEADER_IDS
	extensions |= blackfriday.EXTENSION_FOOTNOTES
	extensions |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_DEFINITION_LISTS

	//for trim relative path
	if !strings.HasSuffix(srcDir, "/") {
		srcDir += "/"
	}
	for _, sourceFile := range sourceFiles {
		article, readErr := ReadArticle(srcDir, sourceFile, basePath)
		if readErr != nil {
			log.Printf("Skipping file %v due to parse error: %v", sourceFile, readErr)
			continue
		}
		location := strings.TrimSuffix(strings.TrimPrefix(sourceFile, srcDir), path.Ext(sourceFile))
		renderer := &HTMLRenderer{
			basePath: basePath,
			location: location,
			Renderer: blackfriday.HtmlRendererWithParameters(htmlFlags, "", "", rendererParameters),
		}

		htmlContent := blackfriday.Markdown(article[KeyRawContent].([]byte), renderer, extensions)
		content := renderCodeTabs(string(htmlContent), 1)
		if showPageTitle {
			if t, ok := article[KeyTitle]; ok {
				title, ok := t.(string)
				if ok && title != ""{
					content = fmt.Sprintf("<h1>%s</h1>%s", title, content)
				}
			}
		}
		article[KeyHtml] = content
		articles = append(articles, &article)
	}
	sort.Sort(articles)
	return
}

/***

###### > GO
```go
  fmt.print("test cloud")
```
###### > bash
```bash
  fmt.print("haha")
```
###### > CPP
```cpp
  this is cpp
```

<div class="tabs">
	<input id="515736620" type="radio" name="tab1" track-name="go" checked="checked"/>
	<label for="515736620">Golang</label>
	<section>1</section>
	<input id="479272831" type="radio" track-name="bash" name="tab1"/>
	<label for="479272831">Bash</label>
	<section>2</section>
	<input id="479272834" type="radio" track-name="cpp" name="tab1"/>
	<label for="479272834">Test</label>
	<section>3</section>
</div>
*/

func renderCodeTabs(src string, startTabsID int) string {
	if startTabsID < 0 {
		return src
	}
	srcContent := src
	firstH6 := strings.Index(src, "<h6 ")
	if firstH6 > 0 {
		src = src[firstH6+4:]
		firstH6Left := strings.Index(src, ">")
		firstH6Right := strings.Index(src, "</h6>")
		firstTitle := src[firstH6Left+1 : firstH6Right]
		if strings.HasPrefix(firstTitle, "&gt; ") {
			src = "<h6 " + src
			var contents []string
			var content string
			for {
				content, src = getNextContent(src, startTabsID)
				if content == "" {
					break
				}
				contents = append(contents, content)
			}
			if len(contents) == 0 {
				startTabsID = -1
				return srcContent
			}
			contents[0] = strings.Replace(contents[0], `/><label for="`, ` checked="checked"/><label for="`, 1)
			tpl := `<div class="tabs">%s</div>`
			tpl = fmt.Sprintf(tpl, strings.Join(contents, "\n"))
			startTabsID ++
			src = renderCodeTabs(src, startTabsID)
			srcContent = srcContent[:firstH6] + tpl + src
			return srcContent
		}
	}
	startTabsID = -1
	return srcContent
}

func getNextContent(src string, idx int) (content, cleft string) {
	left := strings.Index(src, "<")
	if left < 0 || left+4 > len(src) {
		return "", src
	}
	if src[left:left+4] != "<h6 " {
		return "", src
	}
	leftR := strings.Index(src, ">")
	h6Right := strings.Index(src, "</h6>")
	if h6Right <= 0 {
		return "", src
	}
	title := src[leftR+1 : h6Right]
	if !strings.HasPrefix(title, "&gt; ") {
		return "", src
	}
	src = src[h6Right+5:]
	tagContent, src := getNextTagContent(src)
	cleft = src
	label := title[5:]
	trackName := url.QueryEscape(strings.ToLower(label))
	t := strconv.FormatInt(time.Now().UnixNano()/10, 10)
	id := t[len(t)-10:]
	tpl := `<input track-name="%s" id="%s" type="radio" name="%d"/><label for="%s">%s</label><section>%s</section>`
	content = fmt.Sprintf(tpl, trackName, id, idx, id, label, tagContent)
	return
}

func getNextTagContent(src string) (content, cleft string) {
	left := strings.Index(src, "<")
	right := strings.Index(src, ">")
	if left < 0 || right < 0 || left+1 >= right {
		return "", src
	}
	tag := src[left+1 : right]
	if right > left {
		tag = src[left+1 : right]
	}
	if tag == "" {
		return "", src
	}
	rightTag := "</" + tag + ">"
	rightTagIndex := strings.Index(src, rightTag)
	if rightTagIndex > 0 {
		content = src[left:rightTagIndex] + rightTag
		cleft = src[rightTagIndex+len(rightTag):]
		return
	}
	return "", src
}

func getNextTag(src string) string {
	left := strings.Index(src, "<")
	right := strings.Index(src, ">")
	if right > left {
		return src[left+1 : right]
	}
	return ""
}

package mdrest

import (
	"bytes"
	"github.com/russross/blackfriday"
	"strings"
)

type HTMLRenderer struct {
	basePath    string
	location string
	blackfriday.Renderer
}

func (renderer *HTMLRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if bytes.HasPrefix(link,[]byte(".")) || !bytes.Contains(link, []byte("://")) {
		link = []byte(strings.Replace(AbsPath(renderer.basePath, renderer.location,string(link))," ","%20",-1))
	}
	renderer.Renderer.Link(out, link, title, content)
}

func (renderer *HTMLRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	if !(bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://"))) {
		link = []byte(AbsPath(renderer.basePath, renderer.location,string(link)))
	}
	out.WriteString(`<div class="image-package"><img src="`)
	out.Write(link)
	out.WriteString(`" alt="`)
	out.Write(alt)
	out.WriteString(`"`)
	if title != nil {
		out.WriteString(` title="`)
		out.Write(title)
		out.WriteString(`"/>`)
		out.WriteString(`<div class="caption">`)
		out.Write(title)
		out.WriteString(`</div>`)
	} else {
		out.WriteString(`/>`)
	}
	out.WriteString(`</div>`)
}

/***
<div class="tabs">
    <input id="515736620" type="radio" name="tab" checked="checked"/>
    <label for="515736620">Golang</label>
    <section>1</section>
    <input id="479272831" type="radio" name="tab"/>
    <label for="479272831">Bash</label>
    <section>2</section>
    <input id="479272834" type="radio" name="tab"/>
    <label for="479272834">Test</label>
    <section>3</section>
</div>
 */

const tabTpl = `<input id="%d" type="radio" name="tab" checked="checked"/><label for="%d">%s</label><section>%s</section>`

// ListItem adds task list support to the Blackfriday renderer.
func (renderer *HTMLRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
   if bytes.HasPrefix(text, []byte(`<p>[ ] `)) {
		right := bytes.Index(text,[]byte("</p>"))
		text = []byte(`[ ] ` + string(text[7:right]) + string(text[right+4:]))
	}
	switch {
	case bytes.HasPrefix(text, []byte("[ ] ")):
		text = append([]byte(`<input type="checkbox" disabled />`), text[3:]...)

	case bytes.HasPrefix(text, []byte("[x] ")) || bytes.HasPrefix(text, []byte("[X] ")):
		text = append([]byte(`<input type="checkbox" checked disabled />`), text[3:]...)
	}
	renderer.Renderer.ListItem(out, text, flags)
}

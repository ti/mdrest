package mdrest

import (
	"github.com/russross/blackfriday"
	"bytes"
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
	if bytes.HasPrefix(link,[]byte(".")) || !bytes.Contains(link, []byte("://")) {
		link = []byte(AbsPath(renderer.basePath, renderer.location,string(link)))
	}
	renderer.Renderer.Image(out, link, title, alt)
}

// ListItem adds task list support to the Blackfriday renderer.
func (renderer *HTMLRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	switch {
	case bytes.HasPrefix(text, []byte("[ ] ")):
		text = append([]byte(`<input type="checkbox" disabled class="task-list-item">`), text[3:]...)

	case bytes.HasPrefix(text, []byte("[x] ")) || bytes.HasPrefix(text, []byte("[X] ")):
		text = append([]byte(`<input type="checkbox" checked disabled class="task-list-item">`), text[3:]...)
	}
	renderer.Renderer.ListItem(out, text, flags)
}

// List adds task list support to the Blackfriday renderer.
func (renderer *HTMLRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	renderer.Renderer.List(out, text, flags)
	if out.Len() > marker {
		list := out.Bytes()[marker:]
		if bytes.Contains(list, []byte("task-list-item")) {
			// Rewrite the buffer from the marker
			out.Truncate(marker)
			// May be either dl, ul or ol
			list := append(list[:4], append([]byte(` class="task-list"`), list[4:]...)...)
			out.Write(list)
		}
	}
}


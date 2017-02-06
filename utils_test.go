package mdrest

import (
	"testing"
	"log"
	"strings"
)


/**

转换情形：
http://thiss
1. ![name](icon.png)
2. [name](../blob/master/LICENSE)
3. [logo]: ../this/dab.png

忽律的情形：
1. ``` 开头的跳过
2. 转换 过程中陪到 : 的
 */


func TestAbsPath(t *testing.T){
	loc := "/page/pages2/pagss/"
	path := "../../../aseets/icon.png"
	abs := "/aseets/icon.png"
	if abs != AbsPath("",loc,path) {
		t.Fail()
	}
}



const tstHTMLContent = "<!DOCTYPE html><html><head><script src=\"http://two/foobar.js\"></script></head><body><nav><ul><li hugo-nav=\"section_0\"></li><li hugo-nav=\"section_1\"></li></ul></nav><article>content <a href=\"http://two/foobar\">foobar</a>. Follow up</article><p>This is some text.<br>And some more.</p></body></html>"


func TestStripHTML(t *testing.T) {
	type test struct {
		input, expected string
	}
	data := []test{
		{"<h1>strip h1 tag</h1>\n", "strip h1 tag "},
		{"<p> strip p tag </p>", " strip p tag "},
		{"</br> strip br<br>", " strip br\n"},
		{"</br> strip br2<br />", " strip br2\n"},
		{"This <strong>is</strong> a\nnewline", "This is a newline"},
		{"No Tags", "No Tags"},
		{`<p>Summary Next Line.
<figure >

        <img src="/not/real" />


</figure>
.
More text here.</p>

<p>Some more text</p>`, "Summary Next Line.  . More text here.\nSome more text\n"},
	}
	for i, d := range data {
		output := StripHTML(d.input)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}
	}
}
func BenchmarkStripHTML(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StripHTML(tstHTMLContent)
	}
}

func TestTruncateWordsToWholeSentence(t *testing.T) {
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"This is a sentence.", "This is a sentence.", 5, false},
		{"This is also a sentence!", "This is also a sentence!", 1, false},
		{"To be. Or not to be. That's the question.", "To be.", 1, true},
		{" \nThis is not a sentence\nAnd this is another", "This is not a sentence", 4, true},
		{"", "", 10, false},
	}
	for i, d := range data {
		output, truncated := TruncateWordsToWholeSentence(d.input, d.max)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}

func TestTruncateWordsByRune(t *testing.T) {
	type test struct {
		input, expected string
		max             int
		truncated       bool
	}
	data := []test{
		{"", "", 1, false},
		{"a b c", "a b c", 12, false},
		{"a b c", "a b c", 3, false},
		{"a", "a", 1, false},
		{"Hello 中国", "", 0, true},
		{"这是中文，全中文。", "这是中文，", 5, true},
		{"Hello 中国", "Hello 中", 2, true},
		{"Hello 中国", "Hello 中国", 3, false},
		{"Hello中国 Good 好的", "Hello中国 Good 好", 9, true},
		{"This is a sentence.", "This is", 2, true},
		{"This is also a sentence!", "This", 1, true},
		{"To be. Or not to be. That's the question.", "To be. Or not", 4, true},
		{" \nThis is    not a sentence\n ", "This is not", 3, true},
	}
	for i, d := range data {
		output, truncated := TruncateWordsByRune(strings.Fields(d.input), d.max)
		if d.expected != output {
			t.Errorf("Test %d failed. Expected %q got %q", i, d.expected, output)
		}

		if d.truncated != truncated {
			t.Errorf("Test %d failed. Expected truncated=%t got %t", i, d.truncated, truncated)
		}
	}
}


func TestParseRelativePath(t *testing.T){
	loc := "page/pages2/pagss/xx"
	path := "../../../aseets/icon.png"
	result := "aseets/icon.png"

	if result != AbsPath("",loc,path) {
		t.Fail()
	}

	loc2 := "readme"

	path2 := "aseets/icon.png"


	result2 := "aseets/icon.png"

	rl := AbsPath("",loc2,path2)

	log.Println("说吧2",rl)

	if result2 != rl {
		log.Println("说吧2",rl)
		t.Fail()
	}
}

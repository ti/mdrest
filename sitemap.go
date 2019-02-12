package mdrest

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	Title    string `json:"title"`
	Location string `json:"location,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Children Nodes  `json:"children,omitempty"`
}

type Nodes []*Node

func (n *Node) Add(path []string, title, location, icon string) {
	if len(path) == 0 {
		n.Children = append(n.Children, &Node{Title: title, Location: location, Icon:icon})
		return
	}
	idx := -1
	for i, v := range n.Children {
		if v.Title == path[0] {
			idx = i
		}
	}
	if idx == -1 {
		idx = len(n.Children)
		n.Children = append(n.Children, &Node{Title: path[0]})
	}
	n.Children[idx].Add(path[1:], title, location, icon)
}

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

//Less
func (n Nodes) Less(i, j int) bool {
	//the Articles must have date key, if can not read from yaml header, the 'date' be replaced by file last modify time
	left := len(n[i].Children)
	right := len(n[j].Children)
	return left < right
}


//Less
func (this Articles) GetSiteMapMarkdownFromReadme() string {
	for _, article := range this {
		a := (*article)
		location := a[KeyLocation].(string)
		if location == "README" {
			content := string(a[KeyRawContent].([]byte))
			firstIndex := strings.Index(content, "\n# CONTENTS")
			lenContentTitle := len("\n# CONTENTS")
			if firstIndex > 0 && len(content) > firstIndex + 3 {
				remanent :=  content[firstIndex+lenContentTitle+1:]
				remaIndex := strings.Index(remanent, "\n#")
				if remaIndex > 0 {
					remanent = remanent[0:remaIndex]
				}
				if !strings.HasSuffix(remanent,"\n") {
					remanent += "\n"
				}
				return remanent
			}
		}
	}
	return ""

}
//Less
func (this Articles) GetSiteMap(deep int) Nodes {
	siteMap := Node{}
	for _, article := range this {
		a := (*article)
		location := a[KeyLocation].(string)
		sli := strings.SplitN(location, "/", deep+1)
		if len(sli) > deep {
			continue
		}
		var icon string
		if i, ok := a["icon"]; ok {
			icon = fmt.Sprint(i)
		}
		folders := sli[:len(sli)-1]
		location = strings.Replace(location, " ", "%20", -1)
		siteMap.Add(folders, a[KeyTitle].(string), location, icon)
	}
	sort.Sort(siteMap.Children)
	return siteMap.Children
}


var iconClass = "material-icons"

// Markdown convert nodes to Markdown format
func (n Nodes) ToMarkdown(intent ...string) (md string) {
	if len(n) == 0 {
		return ""
	}
	var iconClassTag string
	if iconClass != "" {
		iconClassTag = fmt.Sprintf(` class="%s"`, iconClass)
	}

	intents := strings.Join(intent, "")
	for _, c := range n {
		icon := c.Icon
		if icon != "" {
			icon = fmt.Sprintf(`<i%s>%s</i> `, iconClassTag, c.Icon)
		}
		if c.Location == "README" {
			continue
		}
		if c.Location == "" {
			c.Location = strings.Replace(c.Title, " ", "%20", -1) + "/"
		} else if !strings.HasSuffix(c.Location,"/") && !strings.HasSuffix(c.Location,".md"){
			c.Location += ".md"
		}
		line := fmt.Sprintf("%s* %s[%s](%s)\n", intents, icon, c.Title, c.Location)
		md += line
		if len(c.Children) > 0 {
			newIntent := append(intent, "\t")
			cmd := c.Children.ToMarkdown(newIntent...)
			md += cmd
		}
	}
	return
}

func markdownToNodes(md string,n *Nodes) error {
	md = strings.Replace(md,"    ", "\t",-1)
	md = strings.Replace(md,"  ", "\t",-1)
	lines := strings.Split(md, "\n")
	var pre *Node
	var preMd string
	for _, v := range lines {
		if strings.HasPrefix(v, "\t") {
			line := strings.TrimPrefix(v,"\t")
			preMd += line + "\n"
			continue
		} else  {
			if preMd != "" {
				if pre != nil {
					markdownToNodes(preMd, &pre.Children)
				}
				preMd = ""
			}
			pre = lineToNode(v)
			if pre != nil {
				*n = append(*n, pre)
			}
		}

	}
	return nil
}

func lineToNode(v string) *Node {
	if !(strings.HasPrefix(v,"* ") || strings.HasPrefix(v,"- ") || strings.HasPrefix(v,"+ ")) {
		return nil
	}
	indexOfRightSquare := strings.Index(v, "]")
	indexOfRightRound := strings.Index(v, ")")
	indexOfRightTip := strings.Index(v, "</i>")
	var icon string
	start := 3
	if indexOfRightTip > 0 {
		indexOfLeftTip := strings.LastIndex(v[0:indexOfRightTip],">")
		if indexOfLeftTip > 0 {
			icon = v[indexOfLeftTip+1:indexOfRightTip]
		}
		start = indexOfRightTip + 5
	}
	if indexOfRightRound > 0 && indexOfRightRound > indexOfRightSquare {
		n :=  &Node{
			Icon:icon,
			Title: v[start:indexOfRightSquare],
			Location: v[indexOfRightSquare+2:indexOfRightRound],
		}
		if strings.HasSuffix(n.Location, ".md") {
			n.Location = strings.TrimSuffix(n.Location,".md")
		}
		return n
	}
	return nil
}
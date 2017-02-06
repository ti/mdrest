package mdrest

import (
	"strings"
	"sort"
)

type Node struct {
	Title string `json:"title"`
	Location string `json:"location,omitempty"`
	Children Nodes `json:"children,omitempty"`
}

type Nodes []Node

func (n *Node) Add(path []string, title, location string){
	if len(path) == 0 {
		n.Children = append(n.Children, Node{Title:title, Location:location})
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
		n.Children = append(n.Children, Node{Title:path[0]})
	}
	n.Children[idx].Add(path[1:], title, location)
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
func (this Articles) GetSiteMap(deep int) Nodes {
	siteMap := Node{}
	for _, article := range this {
		a := (*article)
		location := a[KeyLocation].(string)
		sli := strings.SplitN(location, "/",deep + 1)
		if len(sli) > deep {
			continue
		}
		folders := sli[:len(sli)-1]
		siteMap.Add(folders,a[KeyTitle].(string),a[KeyLocation].(string))
	}
	sort.Sort(siteMap.Children)
	return siteMap.Children
}

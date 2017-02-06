package mdrest

import (
	"testing"
)


func TestSiteMapTree(t *testing.T) {
	siteMap := Node{Title:"root"}
	siteMap.Add([]string{},"My Name","xx")
	siteMap.Add([]string{"first","second"},"this is ","xx")
	siteMap.Add([]string{"first"},"how to","xx")
}

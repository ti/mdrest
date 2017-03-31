package mdrest

type Config struct {
	IsCJKLanguage bool    //for trip summery
	Watch    bool
	BasePath string        //the base path of you project, default is "", you can use "/" or "/docs/"
	SrcDir   string
	OutputType string      //you can put json or  html is default
	SiteMapDeep int
	DistDir  string
	NoLogging  bool
	NoIndex    bool
	NoSiteMap  bool
}
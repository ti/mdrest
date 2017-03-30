package mdrest

import (
	"os"
	"strings"
	"path"
	"log"
	"fmt"
	"path/filepath"
)

type MdRest struct {
	cfg  *Config
}

func New(config *Config) *MdRest {
	return &MdRest {
		cfg:config,
	}
}

//Write files
func (mj *MdRest) Do()  error {
	if mj.cfg.SrcDir == "" {
		mj.cfg.SrcDir = "./"
	}
	var err error
	if mj.cfg.SrcDir, err =  filepath.Abs(mj.cfg.SrcDir); err != nil {
		return err
	}
	if mj.cfg.BasePath != "" {
		if !strings.HasSuffix(mj.cfg.BasePath,"/") {
			mj.cfg.BasePath += "/"
		}
	}
	if !strings.HasSuffix(mj.cfg.SrcDir,"/") {
		mj.cfg.SrcDir += "/"
	}
	if mj.cfg.DistDir == "" {
		mj.cfg.DistDir  = mj.cfg.SrcDir + "assets/mdrest"
	}
	if !mj.cfg.NoLogging {
		log.Println("Genrating files to", mj.cfg.DistDir)
	}
	articles, err := ReadArticles(mj.cfg.SrcDir,mj.cfg.BasePath)
	if err != nil {
		return err
	}
	os.Remove(mj.cfg.DistDir)
	if err := os.MkdirAll(mj.cfg.DistDir, os.ModePerm); err != nil {
		return fmt.Errorf("Can not make dist dir %v", err)
	}
	if mj.cfg.OutputType == "" {
		mj.cfg.OutputType = "json"
	}
	articles.WriteAllFiles(mj.cfg.DistDir, mj.cfg.OutputType)
	if !mj.cfg.NoLogging {
		if len(articles) == 0 {
			log.Println("no articles to generated")
		} else {
			log.Printf("Write success %v articles data is generated \n", len(articles))
		}
	}
	if mj.cfg.Watch {
		if !mj.cfg.NoLogging {
			log.Printf("MdRest: Listening to changes in directories: %s except _* or .* folder", mj.cfg.SrcDir)
		}
		Watch(mj.cfg.SrcDir, []string{".md"}, func(event Event, fpath string) {
			location := strings.TrimPrefix(strings.TrimSuffix(fpath,path.Ext(fpath)),mj.cfg.SrcDir)
			log.Println(location,"is", event)
			if event == EventRemove {
				if err := os.Remove(mj.cfg.DistDir + "/" + location + ".json"); err != nil {
					os.Remove(mj.cfg.DistDir + "/" + location)
				}
				articles.Remove(location)
				return
			}
			articles, err := ReadArticles(mj.cfg.SrcDir,mj.cfg.BasePath)
			if err != nil {
				log.Println(err)
				return
			}
			os.Remove(mj.cfg.DistDir)
			os.MkdirAll(mj.cfg.DistDir, os.ModePerm)
			articles.WriteAllFiles(mj.cfg.DistDir, mj.cfg.OutputType)
		})
	}
	return  nil
}
func (mj *MdRest) Close()  error {
	if watcher != nil {
		return  watcher.Close()
	}
	return nil
}


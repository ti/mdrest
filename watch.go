package mdrest

import (
	"strings"
	"os"
	"github.com/fsnotify/fsnotify"
	"log"
	"path"
	"path/filepath"
)

type Event string

const (
	EventRemove Event = "REMOVE"
	EventUpsert Event = "UPSERT"
)


var watcher *fsnotify.Watcher

// WatchFunc type is an adapter to watch path changes
type WatchFunc func(event Event, path string)


func Watch(dir string, exts []string, f WatchFunc) (err error) {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	extMap := map[string]bool{}
	for _, v := range exts {
		extMap[v] = true
	}
	watcherDone := make(chan bool)
	var watchedDirs []string
	walkFunc := func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			return nil
		}
		baseName := path.Base(info.Name())
		if strings.HasPrefix(baseName, "_") || strings.HasPrefix(baseName, ".")  {
			return filepath.SkipDir
		}
		for _, v := range watchedDirs {
			if v == fpath {
				return nil
			}
		}
		watchedDirs = append(watchedDirs, fpath)
		return nil
	}
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if info, err := os.Stat(event.Name); err == nil {
					if info.IsDir() {
						if (event.Op&fsnotify.Remove == fsnotify.Remove) {
							return
						}
						filepath.Walk(event.Name, walkFunc)
						for _, watchedDir := range watchedDirs {
							watcher.Add(watchedDir)
						}
					}
				}
				//md file or may be dir
				if extMap[path.Ext(event.Name)] || !strings.Contains(path.Base(event.Name), ".") {
					//remove must be first
					if (event.Op&fsnotify.Remove == fsnotify.Remove) {
						f(EventRemove,event.Name)
					} else if (event.Op&fsnotify.Write == fsnotify.Write) {
						f(EventUpsert,event.Name)
					} else if (event.Op&fsnotify.Create == fsnotify.Create) {
						f(EventUpsert,event.Name)
					}
				}
			case err := <-watcher.Errors:
				log.Println("Got error:", err)
			}
		}
	}()

	filepath.Walk(dir, walkFunc)

	for _, watchedDir := range watchedDirs {
		watcher.Add(watchedDir)
	}
	<-watcherDone
	return nil
}

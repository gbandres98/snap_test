package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var path = filepath.FromSlash(fmt.Sprint(os.Getenv("SNAP_USER_DATA"), "/dummy.log"))
var quit = make(chan struct{})

func init() {
	fmt.Println(path)
}

func main() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for t := range ticker.C {
			_ = t // we don't print the ticker time, so assign this `t` variable to underscore `_` to avoid error
			if err := writeFile(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	<-quit
}

func writeFile() error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		return err
	}
	defer file.Close()

	line := fmt.Sprintln(time.Now(), " - Snap: ", os.Getenv("SNAP_NAME"), " - Ver: ", os.Getenv("SNAP_REVISION"))
	fmt.Print(line)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return nil
}

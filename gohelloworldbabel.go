package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type TimeResponse struct {
	Status    string `json:"status,omitempty"`
	Formatted string `json:"formatted,omitempty"`
}

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

	r, err := http.Get("http://api.timezonedb.com/v2.1/get-time-zone?key=50SG1ZNZM5CR&by=zone&zone=Europe/Madrid&format=json")
	var timeString string
	if err != nil {
		fmt.Println(err)
		timeString = time.Now().Format("2006-01-02 15:04:05")
	} else {
		defer r.Body.Close()
		timeResponse := &TimeResponse{}
		if err := json.NewDecoder(r.Body).Decode(timeResponse); err != nil {
			return err
		}
		timeString = timeResponse.Formatted
	}

	line := fmt.Sprintln(timeString, "- Snap: ", os.Getenv("SNAP_NAME"), "- Ver: ", os.Getenv("SNAP_REVISION"))
	fmt.Print(line)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return nil
}

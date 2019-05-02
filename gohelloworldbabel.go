package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var path = filepath.FromSlash(fmt.Sprint(os.Getenv("SNAP_USER_DATA"), "/dummy.log"))

func init() {
	fmt.Println(path)
}

func main() {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	line := fmt.Sprintln(time.Now(), " - Snap: ", os.Getenv("SNAP_NAME"), " - Ver: ", os.Getenv("SNAP_REVISION"))
	fmt.Println(line)
	if _, err := file.WriteString(line); err != nil {
		log.Fatal(err)
	}
}

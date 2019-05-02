package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var path = filepath.FromSlash(fmt.Sprint(os.Getenv("SNAP_DATA"), "/dummy.log"))

func init() {
	fmt.Println(path)
}

func main() {
	file, err := os.Open(path)

	if err != nil {
		file, err = os.Create(path)

		if err != nil {
			fmt.Println(err)
		}
	}
	defer file.Close()

	line := fmt.Sprint(time.Now(), " - Snap: ", os.Getenv("SNAP_NAME"), " - Ver: ", os.Getenv("SNAP_REVISION"))
	fmt.Println(line)
	_, _ = file.WriteString(line)
}

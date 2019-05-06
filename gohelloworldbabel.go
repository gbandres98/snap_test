package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
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
	longTicker := time.NewTicker(5 * time.Second)

	go func() {
		for t := range ticker.C {
			_ = t
			if err := writeFile(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	go func() {
		for t := range longTicker.C {
			_ = t
			if err := uploadFile(); err != nil {
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
			timeString = time.Now().Format("2006-01-02 15:04:05") + "*"
		} else {
			timeString = timeResponse.Formatted
		}
	}

	line := fmt.Sprintln(timeString, "- Snap: ", os.Getenv("SNAP_NAME"), "- Ver: ", os.Getenv("SNAP_REVISION"))
	fmt.Print(line)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	return nil
}

func uploadFile() error {
	var client *http.Client
	{
		//setup a mocked http client.
		ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := httputil.DumpRequest(r, true)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", b)
		}))
		defer ts.Close()
		client = ts.Client()
	}

	//prepare the reader instances to encode
	values := map[string]io.Reader{
		"file":  mustOpen(path), // lets assume its this file
		"other": strings.NewReader("hello world!"),
	}
	err := upload(client, values)
	if err != nil {
		panic(err)
	}

	_ = os.Remove(path)

	return nil
}

func upload(client *http.Client, values map[string]io.Reader) (err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an image file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	_ = w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", fmt.Sprint("http://", os.Getenv("GO_SERVER"), ":8080/log/upload"), &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		fmt.Print(err)
	}
	return r
}

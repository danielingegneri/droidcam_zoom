package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	host := flag.String("host", "", "Hostname or IP")
	action := flag.String("action", "", "in|out")
	port := flag.Int("port", 4747, "Port number")
	zoomFile := flag.String("zoom_file", "zoom.txt", "Zoom file")

	flag.Parse()

	if *host == "" || *port == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if *action != "in" && *action != "out" {
		println("Action must be in or out")
		os.Exit(1)
	}
	baseUrl := fmt.Sprintf("http://%s:%d/v2/camera", *host, *port)

	err := makeFileIfReq(*zoomFile)
	if err != nil {
		fmt.Printf("Error creating zoom file: %s\n", err)
		os.Exit(1)
	}

	zoomText, err := os.ReadFile(*zoomFile)
	if err != nil {
		fmt.Printf("Error reading zoom file: %s\n", err)
		os.Exit(1)
	}
	zoom, err := strconv.Atoi(string(zoomText))
	if err != nil {
		println("zoom val in file must be an integer")
		os.Exit(1)
	}

	if *action == "in" {
		zoom += 1
	} else {
		// Out
		zoom = int(math.Max(float64(zoom)-1, 1))
	}

	println("New zoom: " + strconv.Itoa(zoom))
	err = setZoom(baseUrl, zoom)
	if err != nil {
		fmt.Printf("Error setting zoom: %s\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*zoomFile, []byte(strconv.Itoa(zoom)), fs.ModeExclusive)
	if err != nil {
		fmt.Printf("Error writing zoom: %s\n", err)
		os.Exit(1)
	}

}

func makeFileIfReq(filename string) error {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		// Create file
		var zf *os.File
		zf, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer zf.Close()
		_, err = zf.WriteString("1")
		if err != nil {
			return err
		}
	}
	return nil
}

func setZoom(baseUrl string, zoom int) (err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/zoom_level/%d", baseUrl, zoom), nil)
	client.Timeout = time.Second * 3
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

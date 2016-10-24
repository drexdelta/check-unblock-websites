package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	inputFileName  = "input.txt"
	outputFileName = "output.txt"
	uaFileName     = "useragent.txt"
	total          = 0
	correct        = 0
	incorrect      = 0
)

func in(s string, list []string) bool {
	for _, ele := range list {
		if strings.Contains(s, ele) {
			return true
		}
	}
	return false
}

func main() {
	b, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		log.Fatal(err)
	}
	rdr := bufio.NewReader(bytes.NewReader(b))
	//File Contents to io Reader

	userAgent, err := ioutil.ReadFile(uaFileName)
	if err != nil {
		log.Fatal(err)
	}

	//Remove Output if exist
	if _, err := os.Stat("output.txt"); err == nil {
		os.Remove(outputFileName)
	}
	oFile, err := os.OpenFile(outputFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0667)
	if err != nil {
		log.Fatal(err)
	}
	defer oFile.Close()
	//Ouput Pointer to File

	client := &http.Client{}
	//Client Setup

	t1 := time.Now()
	for {
		line, _, err := rdr.ReadLine()
		if err == io.EOF {
			break
		}
		lineMod := strings.ToLower(string(line))
		url := "http://" + lineMod

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("Accept", "*/*;q=0.5")
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("Connection", "close")
		req.Header.Add("Pragma", "no-cache")
		req.Header.Add("User-Agent", string(userAgent))
		//Setting Request Headers

		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				correct++
				oFile.WriteString(url + "\n")
			}
		} else {
			url = "https://" + lineMod
			req, _ := http.NewRequest("GET", url, nil)
			resp, err = client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				if resp.StatusCode == 200 {
					correct++
					oFile.WriteString(url + "\n")
				}
			} else {
				incorrect++
			}
		}
		total++
		fmt.Printf("\rAvailabe/N.A./Total: %d/%d/%d", correct, incorrect, total)
	}
	t2 := time.Since(t1)
	fmt.Println("\r\n", t2)
}

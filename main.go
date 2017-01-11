package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	URL      string
	Method   string
	Data     string
	Interval string
	Header   header
}
type header []string

func (h *header) String() string {
	return fmt.Sprint(*h)
}

func (h *header) Set(value string) error {
	for _, headerString := range strings.Split(value, ",") {
		*h = append(*h, headerString)
	}
	return nil
}

func ParseData(data string) io.Reader {
	if data[0] == '@' {
		content, err := ioutil.ReadFile(data[1:])
		if err != nil {
			fmt.Printf("Impossible d'ouvrir [%s]\n", data[1:])
			os.Exit(1)
		}
		return bytes.NewBuffer(content)

	} else {
		return bytes.NewBufferString(data)
	}
}

func main() {
	var config = Config{}

	flag.StringVar(&config.Method, "X", "GET", "Methode HTTP (POST PUT GET ...)")
	flag.StringVar(&config.Data, "d", "", "Les datas a envoyé")
	flag.StringVar(&config.Interval, "interval", "10ms", "Interval entre les appels")
	flag.Var(&config.Header, "H", "Header")

	flag.Parse()
	config.URL = flag.Arg(0)

	fmt.Println(config)

	if config.URL == "" {
		fmt.Println("Vous devez spécifier une URL")
		os.Exit(1)
	}

	status := make(map[int]int)
	client := &http.Client{}

	duration, err := time.ParseDuration(config.Interval)
	if err != nil {
		fmt.Println("interval invalid")
		os.Exit(1)
	}

	for {
		req, _ := http.NewRequest(config.Method, config.URL, ParseData(config.Data))
		for _, headerString := range config.Header {
			header := strings.Split(headerString, ":")
			if len(header) != 2 {
				fmt.Printf("Header invalid [%s]\n", headerString)
				os.Exit(1)
			}
			req.Header.Add(header[0], header[1])
		}
		resp, err := client.Do(req)
		if err != nil {
			status[0] = status[0] + 1
		} else {
			status[resp.StatusCode] = status[resp.StatusCode] + 1
		}
		if err == nil {
			resp.Body.Close()
		}

		fmt.Printf("\r")
		for code, count := range status {
			fmt.Printf("Status %d : %d ", code, count)
		}
		time.Sleep(duration)
	}

}

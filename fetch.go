package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"mime"
	"net/http"
	urls "net/url"
	"os"
	"path"
)

var (
	green = color.New(color.FgGreen)
	red   = color.New(color.FgRed)
)

func main() {
	save := flag.Bool("d", false, "download contents to a file")
	exec := flag.Bool("x", false, "if -d is set, mark the file as executable")
	help := flag.Bool("h", false, "show usage")
	dflo := flag.Bool("r", false, "don't follow redirects")
	iger := flag.Bool("e", false, "ignore errors and fetch page on status code != 200")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var client *http.Client
	if *dflo {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	} else {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			},
		}
	}

	errors := 0
	for _, url := range flag.Args() {
		if *save {
			green.Printf("Fetching ")
			fmt.Println(url)
		}
		resp, err := client.Get(url)
		check(err)

		if *save {
			if resp.StatusCode != 200 {
				red.Printf("Status ")
			} else {
				green.Printf("Status ")
			}
			fmt.Println(resp.Status)
		}

		if !*iger && resp.StatusCode != 200 {
			if *save {
				red.Printf("Failed ")
				fmt.Printf("with HTTP status code %v\n", resp.StatusCode)
			}
			errors += 1
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		check(err)

		if *save {
			var perm os.FileMode = 0664
			if *exec {
				perm = 0774
			}
			check(err)
			disposition := resp.Header.Get("Content-Disposition")
			_, params, _ := mime.ParseMediaType(disposition)
			filename := params["filename"]
			if filename == "" {
				parsedUrl, _ := urls.Parse(url)
				filename = path.Base(parsedUrl.Path)
				if filename == "." {
					filename = "out"
				}
				mediatype, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
				extensions, err := mime.ExtensionsByType(mediatype)
				var extension string
				if extensions == nil || err != nil {
					extension = fmt.Sprintf(".%v", path.Base(mediatype))
				} else {
					extension = extensions[0]
				}
				filename += extension
			}
			ioutil.WriteFile(filename, body, perm)
		} else {
			fmt.Printf("%s", body)
		}
	}

	if *save {
		if errors > 0 {
			color.Yellow(fmt.Sprintf("%v file(s) failed to download due to non-OK HTTP status codes.", errors))
			color.Yellow("Try running with the '-e' flag to ignore errors and fetch the file(s) anyway.")
		}
	}
}

func check(err error) {
	if err != nil {
		red.Printf("fetch: %v\n", err)
		os.Exit(1)
	}
}

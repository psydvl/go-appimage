// The following directive is necessary to make the package coherent:

//go:build ignore
// +build ignore

// This program generates exclude.go. It can be invoked by running
// go generate
package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	const url = "https://raw.githubusercontent.com/AppImage/pkg2appimage/master/excludelist"

	libraries := getExcludedLibs(url)

	f, err := os.Create("exclude.go")
	die(err)
	defer f.Close()

	packageTemplate.Execute(f, struct {
		Timestamp time.Time
		URL       string
		Libraries []string
	}{
		Timestamp: time.Now(),
		URL:       url,
		Libraries: libraries,
	})

}

func getExcludedLibs(url string) []string {
	var excludeListLibs []string
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var bodyString string
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return (excludeListLibs)
		}
		bodyString = string(bodyBytes)
	} else {
		die(errors.New("error getting excludelist"))
	}

	lines := strings.Split(bodyString, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		line = strings.Split(line, "#")[0]
		excludeListLibs = append(excludeListLibs, strings.TrimSpace(line))
	}
	return (excludeListLibs)
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var packageTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots at
// {{ .Timestamp }}
// using data from
// {{ .URL }}
package exclude

var ExcludedLibraries = []string{
{{- range .Libraries }}
	{{ printf "%q" . }},
{{- end }}
}
`))

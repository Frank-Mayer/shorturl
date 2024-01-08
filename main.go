package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	filename = "urls.txt"
)

type entry struct {
	id  string
	url string
}

func entries(filename string) ([]entry, error) {
	// Read file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries := []entry{}

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := strings.TrimSpace(scanner.Text())

		// comment
		if strings.HasPrefix(content, "#") || strings.HasPrefix(content, "//") {
			continue
		}

		// split: {id}={url}
		parts := strings.Split(scanner.Text(), "=")

		// id is the first part
		id := parts[0]

		// url is the rest (might contain =)
		url := strings.Join(parts[1:], "=")

		e := entry{id: id, url: url}
		entries = append(entries, e)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Warn("no PORT environment variable found, defaulting to " + port)
	}
	return port
}

func logLevel() log.Level {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	}
	log.Warn("invalid LOG_LEVEL environment variable found, defaulting to info")
	return log.InfoLevel
}

func main() {
	log.SetLevel(logLevel())
	entries, err := entries(filename)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		e := e
		http.HandleFunc("/"+e.id, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, e.url, http.StatusTemporaryRedirect)
			log.Debug("redirecting", "id", e.id, "url", e.url)
		})
		log.Info("registered url", "id", e.id, "url", e.url)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("redirecting", "id", "root", "url", "https://frankmayer.dev")
		http.Redirect(w, r, "https://frankmayer.dev", http.StatusTemporaryRedirect)
	})

	p := port()
	log.Info("listening", "port", p)
	err = http.ListenAndServe("0.0.0.0:"+p, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("exiting")
}

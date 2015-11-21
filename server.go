package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"howtotip/models"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const dateFmt = "2006-01-02"
const timeFmt = "2006-01-02 15:04:05"

type config map[string]string

func jsonResponder(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		enc.Encode(map[string]string{"error": err.Error()})
	} else {
		kind := reflect.TypeOf(data).Kind()
		if (kind == reflect.Slice || kind == reflect.Map || kind == reflect.Ptr) && reflect.ValueOf(data).IsNil() {
			if kind == reflect.Slice {
				fmt.Fprintf(w, "[]")
			} else {
				fmt.Fprintf(w, "{}")
			}
		} else {
			enc.Encode(data)
		}
	}
}

func successResponder(w http.ResponseWriter, r *http.Request, err error) {
	jsonResponder(w, r, map[string]bool{"success": err == nil}, err)
}

func countriesHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var data []models.Country

	defer func() { jsonResponder(w, r, data, err) }()

	countries := models.GetCountries()
	data = countries
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
	})
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var data models.Country

	defer func() { jsonResponder(w, r, data, err) }()

	r.ParseForm()

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		return
	}

	data = models.GetCountry(id)
}

func readConfig(filename string) (map[string]string, error) {
	config := map[string]string{
		"dbhost": "localhost",
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("Couldn't read config " + filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(strings.SplitN(scanner.Text(), "#", 2)[0])

		if line == "" {
			continue
		}

		pieces := strings.SplitN(line, "=", 2)
		if len(pieces) != 2 {
			return nil, errors.New(fmt.Sprintf("Couldn't parse line \"%s\"", line))
		}

		config[strings.TrimSpace(pieces[0])] = strings.TrimSpace(pieces[1])
	}

	return config, nil
}

func main() {
	var err error
	var config_file string
	flag.StringVar(&config_file, "config", "./config.cfg", "config file location")
	flag.Parse()
	cfg, err := readConfig(fmt.Sprintf("%v", config_file))
	if err != nil {
		log.Fatal(err.Error())
	}

	if lf := cfg["log_file"]; lf != "" {
		logfile, err := os.OpenFile(lf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0666))
		if err != nil {
			log.Println(err.Error())
		} else {
			log.SetOutput(logfile)
		}
	}

	if pidfile := cfg["pid_file"]; pidfile != "" {
		file, err := os.Create(pidfile)
		if err != nil {
			log.Fatal(err.Error())
		}

		file.WriteString(fmt.Sprintln(os.Getpid()))
		file.Close()
	}

	models.ConnectDB(cfg["dbhost"], cfg["dbport"], cfg["dbuser"], cfg["dbpass"], cfg["dbname"])

	http.HandleFunc("/countries.json", countriesHandler)
	http.HandleFunc("/countries/show.json", countryHandler)

	listen := "127.0.0.1:8080"
	fmt.Println(fmt.Sprintf("listening on %s", listen))
	err = http.ListenAndServe(listen, Log(http.DefaultServeMux))
	if err != nil {
		panic("http.ListenAndServe: " + err.Error())
	}
}

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"howtotip/helpers"
	"howtotip/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
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

	data = models.GetCountries()
}

func countryHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var data models.Country

	defer func() { jsonResponder(w, r, data, err) }()

	r.ParseForm()

	slug := r.FormValue("slug")
	if slug == "" {
		return
	}

	data = models.GetCountry(slug)
}

type Testdata struct {
	Name string
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	var name string
	var data interface{}

	if r.URL.Path == "/" {
		name = "home"
		data = models.GetCountries()
	}	else {
		name = "page"
		slug := strings.Split(r.URL.Path, "/")[1]
		data = models.GetCountry(slug)
	}

	// layout := path.Join("templates", "layout.html")
	// page := path.Join("templates", fmt.Sprintf("%s.html", name))
	// t, _ := template.ParseFiles(layout, page)
	// t.ExecuteTemplate(w, "layout", &data)

	page := path.Join("templates", fmt.Sprintf("%s.html", name))
	t, _ := template.ParseFiles(page)
	t.Execute(w, &data)
}

func routeHandler(router helpers.RegexpRouter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		router.ServeHTTP(w, r)
		log.Printf("%s %s %s %v", r.RemoteAddr, r.Method, r.URL, time.Since(start))
	})
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

	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	router := new(helpers.RegexpRouter)
	router.AddRoute("/countries.json", countriesHandler)
	router.AddRoute("/countries/show.json", countryHandler)
	router.AddRoute("/.*", pageHandler)

	listen := "127.0.0.1:8080"
	fmt.Println(fmt.Sprintf("listening on %s", listen))
	err = http.ListenAndServe(listen, routeHandler(*router))
	if err != nil {
		panic("http.ListenAndServe: " + err.Error())
	}
}

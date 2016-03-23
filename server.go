package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"howtotip/models"
	"howtotip/server"
	"log"
	"net/http"
	"os"
	"strings"
)

const dateFmt = "2006-01-02"
const timeFmt = "2006-01-02 15:04:05"

type config map[string]string

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

	assetServer := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assetServer))

	router := new(server.RegexpRouter)
	router.AddRoute("/api/countries.json", server.CountriesHandler)
	router.AddRoute("/api/countries/show.json", server.CountryHandler)

	router.AddRoute("/countries", server.GetCountriesHandler)
	router.AddRoute("/country_get", server.GetCountryHandler)
	router.AddRoute("/country_post", server.PostCountryHandler)

	router.AddRoute("/.*", server.PageHandler)

	http.Handle("/", server.RouteHandler(*router))

	listen := "127.0.0.1:8080"
	fmt.Println(fmt.Sprintf("listening on %s", listen))
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		panic("http.ListenAndServe: " + err.Error())
	}
}

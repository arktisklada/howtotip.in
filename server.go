package main

import (
	"fmt"
	"howtotip.in/models"
	"howtotip.in/helpers"
	"net/http"
)

const dateFmt = "2006-01-02"
const timeFmt = "2006-01-02 15:04:05"

func main() {
	var err error

	config := helpers.ReadConfig()

	models.ConnectDB(config["dbhost"], config["dbport"], config["dbuser"], config["dbpass"], config["dbname"])

	assetServer := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assetServer))

	router := new(helpers.RegexpRouter)
	router.AddRoute("/api/countries.json", helpers.CountriesHandler)
	router.AddRoute("/api/countries/show.json", helpers.CountryHandler)

	router.AddRoute("/countries", helpers.GetCountriesHandler)
	router.AddRoute("/country_get", helpers.GetCountryHandler)
	router.AddRoute("/country_post", helpers.PostCountryHandler)

	router.AddRoute("/.*", helpers.PageHandler)

	http.Handle("/", helpers.RouteHandler(*router))

	listen := config["listen_host"] + ":" + config["listen_port"]
	fmt.Println(fmt.Sprintf("listening on %s", listen))
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		panic("http.ListenAndServe: " + err.Error())
	}
}

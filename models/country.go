package models

import (
	"database/sql"
)

type Country struct {
	Country string `json:"country"`
	Slug    string `json:"slug"`
	Caption string `json:"caption"`
	Body    string `json:"body"`
}

func GetCountries() (countries []Country) {
	statement, err := db.Prepare("SELECT country, slug, caption, body FROM countries")
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		country, err := buildCountry(rows)
		if err == nil {
			countries = append(countries, country)
		}
	}
	return countries
}

func GetCountry(slug string) (country Country) {
	statement, err := db.Prepare("SELECT country, slug, caption, body FROM countries WHERE slug = $1")
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	rows, err := statement.Query(slug)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		country, err = buildCountry(rows)
		if err != nil {
			panic(err)
		}
	}

	return
}

func buildCountry(row *sql.Rows) (obj Country, err error) {
	var country string
	var slug string
	var caption string
	var body string

	err = row.Scan(&country, &slug, &caption, &body)
	if err == nil {
		obj = Country{
			country,
			slug,
			caption,
			body,
		}
	}

	return
}

package models

import (
	"database/sql"
)

type Country struct {
	Name    string `json:"country"`
	Slug    string `json:"slug"`
	Caption string `json:"caption"`
	Body    string `json:"body"`
}

func GetCountries() (countries []Country) {
	statement, err := db.Prepare("SELECT name, slug, caption, body FROM countries ORDER BY name")
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
	statement, err := db.Prepare("SELECT name, slug, caption, body FROM countries WHERE slug = $1")
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
	var name string
	var slug string
	var caption string
	var body string

	err = row.Scan(&name, &slug, &caption, &body)
	if err == nil {
		obj = Country{
			name,
			slug,
			caption,
			body,
		}
	}

	return
}

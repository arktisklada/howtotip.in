package models

import (
	"database/sql"
)

type Country struct {
	Id      int    `json:"id"`
	Country string `json:"country"`
	Slug    string `json:"slug"`
	Caption string `json:"caption"`
	Body    string `json:"body"`
}

func GetCountries() (countries []Country) {
	statement, err := db.Prepare("SELECT id, country, slug, caption, body FROM countries")
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

func GetCountry(id int) (country Country) {
	statement, err := db.Prepare("SELECT id, country, slug, caption, body FROM countries WHERE id = $1")
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	rows, err := statement.Query(id)
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
	var id int
	var country string
	var slug string
	var caption string
	var body string

	err = row.Scan(&id, &country, &slug, &caption, &body)
	if err == nil {
		obj = Country{
			id,
			country,
			slug,
			caption,
			body,
		}
	}

	return
}

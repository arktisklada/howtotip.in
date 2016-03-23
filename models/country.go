package models

import (
	"database/sql"
)

type Country struct {
	Name    string `json:"country"`
	Slug    string `json:"slug"`
	Caption string `json:"caption"`
	Body    string `json:"body"`
	Live		bool
}

func GetCountries() (countries []Country) {
	statement, err := db.Prepare("SELECT name, slug, caption, body, live FROM countries ORDER BY name")
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
	statement, err := db.Prepare("SELECT name, slug, caption, body, live FROM countries WHERE slug = $1")
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

func SaveCountry(slug, name, caption, live string) (country Country) {
	statement, err := db.Prepare("UPDATE countries SET name=$1, caption=$2, live=$3 WHERE slug = $4")
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	_, err = statement.Exec(name, caption, live, slug)
	if err != nil {
		panic(err)
	}

	country = GetCountry(slug)
	return
}

func buildCountry(row *sql.Rows) (obj Country, err error) {
	var name string
	var slug string
	var caption string
	var body string
	var live bool

	err = row.Scan(&name, &slug, &caption, &body, &live)
	if err == nil {
		obj = Country{
			name,
			slug,
			caption,
			body,
			live,
		}
	}

	return
}

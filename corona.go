package main

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/gometalinter/_linters/src/gopkg.in/yaml.v2"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	//"github.com/google/uuid"
	pb "github.com/acstech/doppler-events/rpc/eventAPI" //c meaning client call
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	clientID string
	eventIDs []string
	c        pb.EventAPIClient
	stop     bool
)

const (
	insertDataStatement string = `INSERT INTO data (id, country, cases, cases_today, deaths, deaths_today, recovered, 
									active, critical, cases_per_one_million, deaths_per_one_million, updated) 
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	insertInfoStatement string = `INSERT INTO info (id, api_id, latitude, longitude, data_id, updated) VALUES ($1, $2, $3, $4, $5, $6)`
)

// Coronavirus struct
type Coronavirus struct {
	Country             string              `json:"country"`
	CountriesInfo       CountryInfo         `json:"countryInfo"`
	Cases               int                 `json:"cases"`
	CasesToday          int                 `json:"todayCases"`
	Deaths              int                 `json:"deaths"`
	DeathsToday         int                 `json:"todayDeaths"`
	Recovered           int                 `json:"recovered"`
	Active              int                 `json:"active"`
	Critical            int                 `json:"critical"`
	CasesPerOneMillion  float32             `json:"casesPerOneMillion"`
	DeathsPerOneMillion float32             `json:"deathsPerOneMillion"`
	Updated             int64 				`json:"updated"`
}

// CountryInfo struct
type CountryInfo struct {
	ID        int     `json:"_id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

// Postgres info
type Postgres struct {
	Host string  `yaml:"host"`
	Port int     `yaml:"port"`
	User string  `yaml:"username"`
	DBName string`yaml:"databaseName"`
}

func main() {

	url := "https://corona.lmao.ninja/countries?sort=country"

	coronaClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	clientID = "Administrator" //In order for test to work, couchbase must contain all 3 clients
	eventIDs = []string{"physical check in", "recent login", "test event"}
	var err error

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "corona-counts")

	res, getErr := coronaClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var countries []Coronavirus

	if err := json.Unmarshal([]byte(body), &countries); err != nil {
		log.Fatal(err)
	}

	var pg Postgres
	values := pg.getPostgres()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		values.Host, values.Port, values.User, values.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	for _, country := range countries {
		id := uuid.New()

		_, err = db.Exec(insertDataStatement, id, country.Country, country.Cases, country.CasesToday, country.Deaths, country.DeathsToday,
			country.Recovered, country.Active, country.Critical, country.CasesPerOneMillion, country.DeathsPerOneMillion,
			time.Unix(country.Updated/1000, 0))
		if err != nil {
			panic(err)
		}

		_, err = db.Exec(insertInfoStatement, uuid.New(), country.CountriesInfo.ID, country.CountriesInfo.Latitude, country.CountriesInfo.Longitude,
			id, time.Unix(country.Updated/1000, 0))
		if err != nil {
			panic(err)
		}
	}
}

func (c *Postgres) getPostgres() *Postgres {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

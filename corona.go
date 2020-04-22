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
	"database/sql"
	//"github.com/google/uuid"
	pb "github.com/acstech/doppler-events/rpc/eventAPI" //c meaning client call

	_ "github.com/lib/pq"
)

var (
	clientID string
	eventIDs []string
	c        pb.EventAPIClient
	stop     bool
	Countries []Coronavirus
	pg Postgres
	db *sql.DB
	err error
)

const (
	insertDataStatement string = `INSERT INTO data (id, country, cases, cases_today, deaths, deaths_today, recovered, 
									active, critical, cases_per_one_million, deaths_per_one_million, updated, time_ran) 
									VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	insertInfoStatement string = `INSERT INTO info (id, api_id, latitude, longitude, data_id, updated, time_ran) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	cleanupCountryData string = `DELETE FROM data a using data b where a.time_ran < b.time_ran and a.cases = b.cases and a.country = b.country;`

	cleanupInfoData string = `DELETE FROM info a using info b where a.time_ran < b.time_ran and a.api_id = b.api_id and a.latitude = b.latitude and a.longitude = b.longitude;`
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

func setUpPostgres() (*sql.DB, error) {
	values := pg.getPostgres()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		values.Host, values.Port, values.User, values.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}


	err = db.Ping()
	if err != nil {
		return db, err
	}

	fmt.Println("Successfully connected! to Postgres")
	return db, nil
}

func requestAPI() error {
	url := "https://corona.lmao.ninja/v2/countries"

	coronaClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Access-Control-Allow-Origin")

	res, getErr := coronaClient.Do(req)
	if getErr != nil {
		return getErr
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return err
	}

	if err := json.Unmarshal([]byte(body), &Countries); err != nil {
		log.Fatal(err)
	}
	return nil
}

func insertDataToPostgres() error {
	for _, country := range Countries {
		id := uuid.New()

		_, err = db.Exec(insertDataStatement, id, country.Country, country.Cases, country.CasesToday, country.Deaths, country.DeathsToday,
			country.Recovered, country.Active, country.Critical, country.CasesPerOneMillion, country.DeathsPerOneMillion,
			time.Unix(country.Updated/1000, 0), time.Now())
		if err != nil {
			return err
		}

		_, err = db.Exec(insertInfoStatement, uuid.New(), country.CountriesInfo.ID, country.CountriesInfo.Latitude, country.CountriesInfo.Longitude,
			id, time.Unix(country.Updated/1000, 0), time.Now())
		if err != nil {
			return err
		}
	}
	fmt.Println("Data Inserted at: " + time.Now().String())
	return nil
}

func cleanupData() error {
	_, err = db.Exec(cleanupCountryData)
	if err != nil {
		return err
	}

	_, err = db.Exec(cleanupInfoData)
	if err != nil {
		return err
	}

	fmt.Println("Data Cleaned at: " + time.Now().String())
	return nil
}

func apiGetTimer() {
	for {
		time.Sleep(1 * time.Hour)
		go callingData()
	}
}

func initialRun() error {
	// call to database to setup
	db, err = setUpPostgres()
	if err != nil {
		return err
	}
	defer db.Close()

	// func to call API to get JSON data for first initial run
	err := requestAPI()
	if err != nil {
		return err
	}

	// inserts data into postgres for first initial run
	err = insertDataToPostgres()
	if err != nil {
		return err
	}

	return nil
}

func callingData() {
	// func to call API to get JSON data
	err := requestAPI()
	if err != nil {
		log.Fatal(err)
	}

	// inserts data each hour
	err = insertDataToPostgres()
	if err != nil {
		log.Fatal(err)
	}

	// cleans up data each hour
	err = cleanupData()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err = initialRun()
	if err != nil {
		log.Fatal(err)
	}

	// call polling function to recall API GET call after an hour
	apiGetTimer()
}



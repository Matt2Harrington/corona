package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

// Coronavirus struct
type Data struct {
	ID 					uuid.UUID
	Country             string
	Cases               int
	CasesToday          int
	Deaths              int
	DeathsToday         int
	Recovered           int
	Active              int
	Critical            int
	CasesPerOneMillion  float32
	DeathsPerOneMillion float32
	Updated             time.Time
	TimeRan  	  		time.Time
	//Info 				Info

}

// CountryInfo struct
type Info struct {
	ID 		  	  uuid.UUID
	APIID        int
	Latitude  	  float64
	Longitude 	  float64
	DataID    	  uuid.UUID
	Updated  	  time.Time
	TimeRan  	  time.Time
}

type dataInfo struct {
	DataList []Data
	InfoList []Info
}

func homePage(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/corona", indexHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

// parseParams accepts a req and returns the `num` path tokens found after the `prefix`.
// returns an error if the number of tokens are less or more than expected
func parseParams(req *http.Request, prefix string, num int) ([]string, error) {
	url := strings.TrimPrefix(req.URL.Path, prefix)
	params := strings.Split(url, "/")
	if len(params) != num || len(params[0]) == 0 || len(params[1]) == 0 {
		return nil, fmt.Errorf("Bad format. Expecting exactly %d params", num)
	}
	return params, nil
}

// indexHandler calls `queryRepos()` and marshals the result as JSON
func indexHandler(w http.ResponseWriter, req *http.Request) {
	data := dataInfo{}

	err = queryData(&data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}
// queryRepos first fetches the repositories data from the db
func queryData(dataInfo *dataInfo) error {
	rows, err := db.Query(`
		SELECT * FROM data ORDER BY cases DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		data := Data{}
		//info := Info{}
		err = rows.Scan(
			&data.ID,
			&data.Country,
			&data.Cases,
			&data.CasesToday,
			&data.Deaths,
			&data.DeathsToday,
			&data.Recovered,
			&data.Active,
			&data.Critical,
			&data.CasesPerOneMillion,
			&data.DeathsPerOneMillion,
			&data.Updated,
			&data.TimeRan,
			//&info.ID,
			//&info.APIID,
			//&info.Latitude,
			//&info.Longitude,
			//&info.DataID,
			//&info.Updated,
			//&info.TimeRan,
		)
		if err != nil {
			return err
		}
		dataInfo.DataList = append(dataInfo.DataList, data)
		//dataInfo.InfoList = append(dataInfo.InfoList, info)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}


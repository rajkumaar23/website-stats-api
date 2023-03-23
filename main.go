package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var HttpClient = &http.Client{}
var PlausibleAPIToken = "iRx0iQDTT7S9fB2dKFByLne4C-Ee_Z-X_sYlxNiuX2lg1Oty2hLXzqcNI8bHvhCu"

type PlausibleAPIResponse struct {
	Results struct {
		Visitors struct {
			Value int `json:"value"`
		} `json:"visitors"`
	} `json:"results"`
}

type GoatCounterAPIResponse struct {
	Count string `json:"count_unique"`
}

func GetGoatCounterStats(page string) int {
	goatCounterURL := fmt.Sprintf("https://rajkumaar23.goatcounter.com/counter/%v.json", strings.TrimSuffix(page, "/"))
	req, _ := http.NewRequest(http.MethodGet, goatCounterURL, nil)
	res, err := HttpClient.Do(req)

	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var response GoatCounterAPIResponse
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	count, err := strconv.Atoi(strings.ReplaceAll(response.Count, "\u202f", ""))
	if err != nil {
		fmt.Println(err)
		return 0
	}

	return count
}

func GetPlausibleStats(page string) int {
	queryParams := url.Values{
		"site_id": {"rajkumaar.co.in"},
		"period":  {"custom"},
		"date":    {fmt.Sprintf("2022-01-01,%v", time.Now().Format(time.DateOnly))},
		"filters": {fmt.Sprintf("event:page==%v", page)},
	}

	apiUrl := fmt.Sprintf("https://plausible.pi.rajkumaar.co.in/api/v1/stats/aggregate?%v", queryParams.Encode())
	req, _ := http.NewRequest(http.MethodGet, apiUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", PlausibleAPIToken))
	res, err := HttpClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	var response PlausibleAPIResponse
	if err = json.Unmarshal(body, &response); err != nil {
		panic(err)
	}

	return response.Results.Visitors.Value
}

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		page := request.URL.Query().Get("page")
		fmt.Fprintf(writer, strconv.Itoa(GetGoatCounterStats(page)+GetPlausibleStats(page)))
		return
	})

	err := http.ListenAndServe(":2304", nil)
	if err != nil {
		panic(err)
	}
}

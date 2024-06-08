package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type ReportingStructure struct {
	ReportingPlans    []ReportingPlan   `json:"reporting_plans"`
	InNetworkFiles    []FileLocation    `json:"in_network_files"`
	AllowedAmountFile AllowedAmountFile `json:"allowed_amount_file"`
}

type ReportingPlan struct {
	PlanName       string `json:"plan_name"`
	PlanIdType     string `json:"plan_id_type"`
	PlanId         string `json:"plan_id"`
	PlanMarketType string `json:"plan_market_type"`
}

type FileLocation struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}

type AllowedAmountFile struct {
	Description string `json:"description"`
	Location    string `json:"location"`
}

func main() {
	startTime := time.Now()

	path := "2024-06-01_anthem_index.json.gz"
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	gz, err := gzip.NewReader(file)

	if err != nil {
		log.Fatal(err)
	}

	defer gz.Close()

	dec := json.NewDecoder(gz)

	var urls []string

	// NOTE: The following tags/identifiers appear to correspond to Anthem PPOs in NY state
	identifiers := []string{"_39F0_", "_71B0_", "_72B0_", "_42F0_", "_39B0_", "_42B0_", "_71A0_"}

	for dec.More() {
		t, _ := dec.Token()
		if t == "reporting_structure" {
			_, _ = dec.Token()

			for dec.More() {
				var structure ReportingStructure

				err := dec.Decode(&structure)
				if err != nil {
					log.Fatal(err)
				}

				if len(structure.InNetworkFiles) > 0 {
					for _, file := range structure.InNetworkFiles {
						for _, id := range identifiers {
							if strings.Contains(file.Location, id) {
								urls = append(urls, file.Location)
								break
							}
						}
					}
				}

				urls = uniq(urls)
			}
		}
	}

	out, err := os.OpenFile("urls.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	writer := bufio.NewWriter(out)

	defer writer.Flush()

	for _, url := range urls {
		_, _ = writer.WriteString(url + "\n")
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	fmt.Println("duration: ", duration)
}

func uniq(slice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}

	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}

	return list
}

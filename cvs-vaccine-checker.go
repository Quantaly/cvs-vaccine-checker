package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Quantaly/cvs-vaccine-checker/structs"
)

const apiEndpoint string = "https://www.cvs.com/immunizations/covid-19-vaccine.vaccine-status.CO.json?vaccineinfo"

// for me, this was defined as "a ~30-minute-or-less drive according to Google Maps"
var nearbyCities []string = []string{
	"AURORA",
	"BOULDER",
	"BRIGHTON",
	"DENVER",
	"EDGEWATER",
	"GLENDALE",
	"LAKEWOOD",
	"SHERIDAN",
	"THORNTON",
}

func main() {
	_, forceDump := os.LookupEnv("FORCE_DUMP")
	lastState := structs.Persistence{
		LastTimestamp: "yeet",
		HadVaccines:   false,
	}
	stateJson, err := ioutil.ReadFile(".cvs-vaccine-checker-persistence")
	if err == nil {
		json.Unmarshal(stateJson, &lastState) // don't care about the error here
	}

	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("Referer", "https://www.cvs.com/immunizations/covid-19-vaccine")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var data structs.Response
	if err = json.Unmarshal(body, &data); err != nil {
		log.Fatalln(err)
	}
	if forceDump || data.ResponsePayloadData.CurrentTime != lastState.LastTimestamp {
		fields := make([]structs.Field, 0, len(nearbyCities))
		for _, city := range data.ResponsePayloadData.Data.CO {
			if forceDump || city.Status != "Fully Booked" {
				for _, nearby := range nearbyCities {
					if city.City == nearby {
						fields = append(fields, structs.Field{
							Name:   city.City,
							Value:  city.Status,
							Inline: true,
						})
						break
					}
				}
			}
		}
		if forceDump || len(fields) > 0 {
			content := "A vaccine might be available!!"
			if forceDump {
				content = "Here's the current status:"
			}
			message := structs.Message{
				Content: content,
				Embeds: []structs.Embed{{
					Title:       "CVS Vaccine Availability",
					Type:        "rich",
					Description: "COVID-19 Vaccines in Colorado",
					Url:         "https://www.cvs.com/immunizations/covid-19-vaccine",
					Timestamp:   data.ResponsePayloadData.CurrentTime + "-0700",
					Color:       0xcc0000,
					Fields:      fields,
				}},
			}
			if err = message.Send(); err != nil {
				log.Fatalln(err)
			}
		} else if lastState.HadVaccines {
			message := structs.Message{
				Content: "Looks like they've run out.",
				Embeds:  nil,
			}
			if err = message.Send(); err != nil {
				log.Fatalln(err)
			}
		}

		saveState := structs.Persistence{
			LastTimestamp: data.ResponsePayloadData.CurrentTime,
			HadVaccines:   !forceDump && len(fields) > 0,
		}
		save, err := json.Marshal(saveState)
		if err != nil {
			log.Fatalln(err)
		}
		if err = ioutil.WriteFile(".cvs-vaccine-checker-persistence", save, 0o644); err != nil {
			log.Fatalln(err)
		}
	}
}

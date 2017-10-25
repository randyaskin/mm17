package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

type savingsReport struct {
	InitialCost           int         `json:"InitialCost"`
	SolarEnergyGeneration int         `json:"SolarEnergyGenerationKWH"`
	PredictedConsumption  float64     `json:"PredictedConsumptionKWH"`
	SolarOffset           float64     `json:"SolarOffsetKWH"`
	HistoricalConsumption []DataPoint `json:"HistoricalConsumptionKWH"`
}

var sdEnergyHistorical = []float64{14395.6602, 13768.9679, 14501.41184, 14421.33019, 14733.81518, 14898.06997, 15359.51847, 15979.35501, 15489.09043, 13984.35539, 17475.61581, 16741.03086, 17120.1648, 17772.47179, 18706.83082, 18714.43401, 19664.75597, 19638.47435, 19994.82863, 19515.30326, 18978.25227, 19022.59425, 19562.03057, 19425.7796, 19903.80127, 19781.17809}
var r *regression.Regression

type DataPoint struct {
	Year   int     `json:"Year"`
	Energy float64 `json:"Energy"`
}

const (
	resCostPerInstall = 27648  // 8 kwH
	comCostPerInstall = 331776 // 96 kWh
	indCostPerInstall = 442368 // 128 kWh

	resEnergyPerInstall = 6800
	comEnergyPerInstall = 81600
	indEnergyPerInstall = 108800
)

var historicalConsumption []DataPoint
var currentYear = 2015

func main() {
	// train model on startup
	r = new(regression.Regression)
	r.SetObserved("City of San Diego Energy usage")
	r.SetVar(0, "Year")
	r.SetVar(1, "Energy consumption in GWh")
	for i, year := range sdEnergyHistorical {
		r.Train(regression.DataPoint(year, []float64{float64(i + 1990)}))
	}
	r.Run()

	startYear := 1990
	for i, val := range sdEnergyHistorical {
		historicalConsumption = append(historicalConsumption, DataPoint{
			Year:   startYear + i,
			Energy: val,
		})

	}

	http.HandleFunc("/v1/savings", savingsHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func savingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	comCount, err := strconv.Atoi(r.URL.Query().Get("com"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}
	resCount, err := strconv.Atoi(r.URL.Query().Get("res"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}
	indCount, err := strconv.Atoi(r.URL.Query().Get("ind"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}
	targetYear, err := strconv.Atoi(r.URL.Query().Get("targetYear"))
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}

	solarGeneration := getSolarGeneration(comCount, resCount, indCount)
	predictedConsumption, err := getPredictedConsumption(targetYear)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}
	sr := savingsReport{
		InitialCost:           getInitialCost(comCount, resCount, indCount),
		SolarEnergyGeneration: solarGeneration,
		PredictedConsumption:  predictedConsumption,
		SolarOffset:           predictedConsumption - float64(solarGeneration),
		// HistoricalConsumption: historicalConsumption,
	}

	respBytes, err := json.Marshal(sr)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}

	w.Write(respBytes)
}

func getInitialCost(resCount, comCount, indCount int) int {
	return (resCount * resCostPerInstall) + (comCount * comCostPerInstall) + (indCount * indCostPerInstall)
}

func getSolarGeneration(resCount, comCount, indCount int) int {
	return (resCount * resEnergyPerInstall) + (comCount * comEnergyPerInstall) + (indCount * indEnergyPerInstall)
}

func getPredictedConsumption(targetYear int) (float64, error) {
	// years := targetYear - currentYear
	// var forecastedYears []DataPoint
	// for i := 0; i < years; i++ {
	// 	year := currentYear + i
	// 	fc, err := r.Predict([]float64{float64(year)})
	// 	if err != nil {
	// 		return 0, err
	// 	}
	// 	forecastedYears = append(forecastedYears, DataPoint{
	// 		Year:   year,
	// 		Energy: fc,
	// 	})
	// }

	fc, err := r.Predict([]float64{float64(targetYear)})
	if err != nil {
		return 0, err
	}

	// convert to kWh
	fc = fc * 1000000

	return fc, nil
}

// old code to get total consumption
// var forecastedYears []DataPoint
// for i := 0; i < years; i++ {
// 	year := currentYear + i
// 	fc, err := r.Predict([]float64{float64(year)})
// 	if err != nil {
// 		return nil, err
// 	}
// 	forecastedYears = append(forecastedYears, DataPoint{
// 		Year:   year,
// 		Energy: fc,
// 	})
// }

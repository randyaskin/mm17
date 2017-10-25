package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	Subsidies             []float64   `json:"Subsidies"`
	TotalSubsidy          float64     `json:"TotalSubsidy"`
	AvgSubsidy            float64     `json:"AvgSubsidy"`
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

	bracket1A     = 1200
	bracket1B     = 1500000000
	bracket1Scale = 10000

	bracket2A     = 16000
	bracket2B     = 1500000000
	bracket2Scale = 9000

	bracket3A     = 26000
	bracket3B     = 1500000000
	bracket3Scale = 8000

	bracket4A     = 37000
	bracket4B     = 1500000000
	bracket4Scale = 7000

	bracket5A     = 31000
	bracket5B     = 1500000000
	bracket5Scale = 5000
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

	var bracketCounts []float64
	var subsidies []float64
	for i := 1; i <= 9; i++ {
		var subsidy float64
		bracketCount, err := strconv.Atoi(r.URL.Query().Get(fmt.Sprintf("p%v", i)))
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
			return
		}
		switch i {
		case 5:
			subsidy = bellCurve(bracket1Scale, bracket1A, bracket1B, float64(bracketCount))
		case 6:
			subsidy = bellCurve(bracket2Scale, bracket2A, bracket2B, float64(bracketCount))
		case 7:
			subsidy = bellCurve(bracket3Scale, bracket3A, bracket3B, float64(bracketCount))
		case 8:
			subsidy = bellCurve(bracket4Scale, bracket4A, bracket4B, float64(bracketCount))
		case 9:
			subsidy = bellCurve(bracket5Scale, bracket5A, bracket5B, float64(bracketCount))
		default:
			subsidy = 0
		}
		subsidies = append(subsidies, subsidy)
		bracketCounts = append(bracketCounts, float64(bracketCount))
	}

	solarGeneration := getSolarGeneration(comCount, resCount, indCount)
	predictedConsumption, err := getPredictedConsumption(targetYear)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error occurred: %v", err)))
		return
	}

	var totalSubsidy float64
	for i, s := range subsidies {
		totalSubsidy = totalSubsidy + s*bracketCounts[i]
	}

	var totalPeople float64
	for _, bc := range bracketCounts {
		totalPeople = totalPeople + bc
	}
	sr := savingsReport{
		InitialCost:           getInitialCost(comCount, resCount, indCount),
		SolarEnergyGeneration: solarGeneration,
		PredictedConsumption:  predictedConsumption,
		SolarOffset:           predictedConsumption - float64(solarGeneration),
		Subsidies:             subsidies,
		TotalSubsidy:          totalSubsidy,
		AvgSubsidy:            totalSubsidy / totalPeople,
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

// func getZeroEnergyOffset(totalConsumption int, cnDate int, subsidy float64) int {
// 	dateRange := cnDate - 2017

// 	totPop := 1073089

// 	percByIncome := [5.9,4.8,8.7,8.2,12,16.9,13,15.8,7.4]

// 	popByIncome := []
// 	for i, val := range percByIncome {
//         popByIncome = append(popByIncome, DataPoint{
//             Year:   startYear + i,
//             Energy: val,
//         })
//     }

// 	return 0
// }

func bellCurve(scale, a, b, x float64) float64 {
	//	return 5000e^(-(x-5000)^2/(3000))
	// fmt.Println("ANS", math.Pow(2.71828182, (x-a)*(x-a)*-1))
	y := ((x - a) * (x - a) * -1) / b
	val := scale * math.Pow(math.E, y)
	fmt.Println("scale", scale)
	fmt.Println("a", a)
	fmt.Println("b", b)
	fmt.Println("x", x)
	fmt.Println("y", y)
	fmt.Println("val", val)
	fmt.Println()
	return val
}

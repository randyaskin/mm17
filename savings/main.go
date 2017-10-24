package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type savingsReport struct {
	InitialCost           int `json:"InitialCost"`
	TotalCityConsumption  int `json:"TotalCityConsumption"`
	SolarEnergyGeneration int `json:"SolarEnergyGeneration"`
	EnergyOffset          int `json:"EnergyOffset"`
}

const (
	resCostPerInstall = 13824 // 4 kwH
	comCostPerInstall = 55296 // 16 kWh
	indCostPerInstall = 82944 // 24 kWh

	resEnergyPerInstall = 3400
	comEnergyPerInstall = 13600
	indEnergyPerInstall = 20400
)

func main() {
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

	totalConsumption := 2000000000
	solarGeneration := getSolarGeneration(comCount, resCount, indCount)
	sr := savingsReport{
		InitialCost:           getInitialCost(comCount, resCount, indCount),
		TotalCityConsumption:  totalConsumption,
		SolarEnergyGeneration: solarGeneration,
		EnergyOffset:          totalConsumption - solarGeneration,
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

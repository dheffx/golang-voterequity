package main

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"sort"
	"strconv"
)

type LoadResources struct {
	Directory      string
	RawDataFile string
}

func (this *LoadResources) LoadData() (data ResourceData) {
	f, err := ioutil.ReadFile(this.GetResourceFile())
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(f, &data)
	data.SetTotalPopulation()
	data.SetTotalVotes()
	data.SetAveragePopulation()
	data.SetAverageElectoralVotes()
	return
}

type State struct {
	Name string
	Population int
	ElectoralVotes int
	VoteEquity float32
}

func (this *State) ToString() (pretty string) {
	pretty = this.Name + "," + strconv.FormatFloat(float64(this.VoteEquity), 'f', 6, 32)
	return
}

type SortByEquity []State

func (s SortByEquity) Len() int {
	return len(s)
}

func (s SortByEquity) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortByEquity) Less(i, j int) bool {
	return s[i].VoteEquity > s[j].VoteEquity
}

type ResourceData struct {
	States []State
	TotalPopulation int
	TotalElectoralVotes int
	AveragePopulation float32
	AverageElectoralVotes float32
}

func (this *ResourceData) SetTotalPopulation() {
	for _, state := range (this.States) {
		this.TotalPopulation += state.Population
	}
}

func (this *ResourceData) SetTotalVotes() {
	for _, state := range (this.States) {
		this.TotalElectoralVotes += state.ElectoralVotes
	}
}

/**

 VoteEquity =
 	(State.ElectoralVotes / TotalElectoralVotes) /
 	(State.Population / TotalPopulation)

 */
func (this *ResourceData) CalculateVoteEquity(state *State) (equity float32){
	equity = this.ElectoralVoteRate(state.ElectoralVotes) /
			 this.PopulationRate(state.Population)
	state.VoteEquity = equity
	return
}

func (this *ResourceData) SetAveragePopulation() {
	this.AveragePopulation = float32(this.TotalPopulation) / float32(len(this.States))
	return
}

func (this *ResourceData) SetAverageElectoralVotes() {
	this.AverageElectoralVotes = float32(this.TotalElectoralVotes) / float32(len(this.States))
	return
}

func (this *ResourceData) PopulationRate(population int) (rate float32) {
	rate = float32(population) / float32(this.TotalPopulation)
	return
}

func (this *ResourceData) ElectoralVoteRate(population int) (rate float32) {
	rate = float32(population) / float32(this.TotalElectoralVotes)
	return
}

func (this *ResourceData) Calculate() {
	for i, state := range this.States {
		this.CalculateVoteEquity(&state)
		this.States[i] = state
	}
}

func (this *ResourceData) SortByVoteEquity() {
	sort.Sort(SortByEquity(this.States))
}

func (this *ResourceData) ToString() (pretty string) {
	for i, state := range this.States {
		pretty += strconv.Itoa(this.StateRank(i)) + "," + state.ToString() + "\n"
	}
	return
}

func (this *ResourceData) ToJsonFile(filename string) {
	jsonString, _ := json.Marshal(this)
	err := ioutil.WriteFile(filename, jsonString, 0644)
	if err != nil {
		fmt.Println("Could not write to file:", filename)
	}
}

func (this *ResourceData) StateRank(order int) (rank int) {
	rank = len(this.States) - order
	return
}

func (this *LoadResources) GetResourceFile() (path string) {
	path, _ = filepath.Abs(filepath.Join(this.Directory, this.RawDataFile))
	return
}

func main() {
	var resources = &LoadResources{Directory: "data", RawDataFile: "state_data.json"}
	var data = resources.LoadData()
	data.Calculate()
	data.SortByVoteEquity()
	fmt.Println(data.ToString())
	data.ToJsonFile("stateequity.json")
	avg_equity := data.AverageElectoralVotes / data.AveragePopulation

	fmt.Println(
		strconv.FormatFloat(float64(data.AverageElectoralVotes), 'f', 8, 32),
		strconv.FormatFloat(float64(data.AveragePopulation), 'f', 8, 32))

	fmt.Println(strconv.FormatFloat(float64(avg_equity), 'f', 8, 32))
}


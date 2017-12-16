package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
)

type day int

const (
	WD day = iota
	SUN
	SAT
)

const data_dir = "./data"

func services(d day) (map[string]struct{}, error) {
	var col int
	switch d {
	case WD:
		col = 1
	case SUN:
		col = 7
	case SAT:
		col = 6
	default:
		panic("fake day")
	}
	cal_f, err := os.Open(data_dir + "/calendar.txt")
	if err != nil {
		return nil, fmt.Errorf("Error opening calendar file: %s", err)
	}
	defer cal_f.Close()
	cal := csv.NewReader(cal_f)
	cal.ReuseRecord = true
	var line []string
	ret := make(map[string]struct{})
	for line, err = cal.Read(); err == nil; line, err = cal.Read() {
		if line[col] == "1" {
			ret[line[0]] = struct{}{}
		}
	}
	if err != io.EOF {
		return nil, fmt.Errorf("Error reading calendar file: %s", err)
	}
	return ret, nil
}

func trips(route string, d day) (map[string]struct{}, error) {
	serviceList, err := services(d)
	if err != nil {
		return nil, err
	}
	tripList_f, err := os.Open(data_dir + "/trips.txt")
	if err != nil {
		return nil, fmt.Errorf("Error opening trips file: %s", err)
	}
	defer tripList_f.Close()
	tripList := csv.NewReader(tripList_f)
	tripList.ReuseRecord = true
	var line []string
	ret := make(map[string]struct{})
	for line, err = tripList.Read(); err == nil; line, err = tripList.Read() {
		if line[0] == route {
			if _, ok := serviceList[line[1]]; ok {
				ret[line[2]] = struct{}{}
			}
		}
	}
	if err != io.EOF {
		return nil, fmt.Errorf("Error reading trips file: %s", err)
	}
	return ret, nil
}

func times(route, stop_id string, d day) ([]string, error) {
	tripList, err := trips(route, d)
	if err != nil {
		return nil, err
	}
	timesList_f, err := os.Open(data_dir + "/stop_times.txt")
	if err != nil {
		return nil, fmt.Errorf("Error opening stop_times file: %s", err)
	}
	defer timesList_f.Close()
	timesList := csv.NewReader(timesList_f)
	timesList.ReuseRecord = true
	var line []string
	var ret []string
	for line, err = timesList.Read(); err == nil; line, err = timesList.Read() {
		if line[3] == stop_id {
			if _, ok := tripList[line[0]]; ok {
				ret = append(ret, line[2])
			}
		}
	}
	if err != io.EOF {
		return nil, fmt.Errorf("Error reading stop_times file: %s", err)
	}
	return ret, nil
}

func main() {
	t, err := times("84", "1102", SUN)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	sort.Strings(t)
	fmt.Println(t)
}

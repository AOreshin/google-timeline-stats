package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

type result struct {
	distance float64
	rides    int
}

func main() {
	input := flag.String("input", "", "")
	flag.Parse()
	paths := strings.Split(*input, ",")
	totalResult := map[string]result{}
	for _, path := range paths {
		fmt.Println(path)
		json := getString(path)
		activityTypes := getActivityTypes(json)
		for _, activityType := range activityTypes {
			distance, rides := calculateDistance(json, activityType)
			fmt.Printf("Distance traveled by %s is %.3f kilometers in %d rides\n", activityType, distance, rides)
			res := totalResult[activityType]
			res.distance += distance
			res.rides += rides
			totalResult[activityType] = res
		}
	}
	fmt.Println("Total")
	for activityType, result := range totalResult {
		fmt.Printf("Distance traveled by %s is %.3f kilometers in %d rides\n", activityType, result.distance, result.rides)
	}
}

func calculateDistance(json, activityType string) (float64, int) {
	distancePath := fmt.Sprintf("timelineObjects.#(activitySegment.activityType=\"%s\")#.activitySegment.distance", activityType)
	distanceResult := gjson.Get(json, distancePath)
	if !distanceResult.IsArray() {
		panic(fmt.Sprintf("%s is not an array", distancePath))
	}
	sum := float64(0)
	for _, distance := range distanceResult.Array() {
		sum += distance.Float()
	}
	return sum / 1000, len(distanceResult.Array())
}

func getActivityTypes(json string) []string {
	activityTypePath := "timelineObjects.#.activitySegment.activityType"
	activityTypeResult := gjson.Get(json, activityTypePath)
	if !activityTypeResult.IsArray() {
		panic(fmt.Sprintf("%s is not an array", activityTypePath))
	}
	activityTypes := map[string]bool{}
	for _, result := range activityTypeResult.Array() {
		if activityTypes[result.String()] == false {
			activityTypes[result.String()] = true
		}
	}
	result := []string{}
	for key := range activityTypes {
		result = append(result, key)
	}
	return result
}

func getString(path string) string {
	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

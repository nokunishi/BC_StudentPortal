package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseCsv(p string) []Student {
	f, err := os.Open(p)
	if err != nil {
		log.Fatal("unable to open input csv file")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	fs, err := csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	var students []Student

	for _, line := range fs {
		s := NewModel(line)
		students = append(students, s)
	}

	return students
}

func main() {
	db, err := NewDriver("../data")

	if err != nil {
		fmt.Println("Error", err)
	}

	students := parseCsv("../input.csv")

	for _, value := range students {
		name := value.Name.Family + "_" + value.Name.Given
		school := strings.ToLower(value.School)

		if err := db.Write(school, name, Student{
			Name:    value.Name,
			Age:     value.Age,
			Phone:   value.Phone,
			School:  value.School,
			Address: value.Address,
		}); err != nil {
			fmt.Println(err)
		}
	}

	if _, err := db.ReadAll("UBCV"); err != nil {
		fmt.Println("Error:", err)
	}
}

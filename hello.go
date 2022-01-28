package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type csvLine struct {
	ID    string `json:"ID"`
	Value string `json:"value"`
}

func main() {
	router := gin.Default()
	router.GET("/data", getAll)
	router.GET("/data/:id", getById)

	router.Run("localhost:8080")
}

func getAll(c *gin.Context) {
	records := readCsvFile("./data.csv")
	c.IndentedJSON(http.StatusOK, records)
}

func getById(c *gin.Context) {
	records := readCsvFile("./data.csv")
	id := c.Param("id")
	for _, a := range records {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Id not found"})
}

func readCsvFile(filePath string) []csvLine {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	var csvData = []csvLine{}
	for _, value := range records {
		newLine := csvLine{ID: value[0], Value: value[1]}
		csvData = append(csvData, newLine)
	}
	return csvData
}

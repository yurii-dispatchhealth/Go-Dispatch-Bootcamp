package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"sync"

	"Go-Dispatch-Bootcamp/types"
)

type readerService struct{}

func New() *readerService {
	return &readerService{}
}

func (ts *readerService) readCsvFromFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New("can not open file")
	}
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("read file error")
	}
	return records, nil
}

func (ts *readerService) FetchCsvFromRemote(feedUrl string) ([][]string, error) {
	url := feedUrl

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	data, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("can not read csv from remote")
	}

	return data, nil
}

func (ts *readerService) Update(users *[]types.User, dataFileName string) (bool, error) {
	file, err := os.Create(dataFileName)
	defer file.Close()
	if err != nil {
		return false, errors.New(fmt.Sprintf("can not open file: %v", dataFileName))
	}
	log.Println(users)

	var csvUsers [][]string

	for _, user := range *users {
		csvUser := []string{
			strconv.Itoa(user.Id),
			user.Username,
			user.Identifier,
			user.FirstName,
			user.LastName,
		}

		csvUsers = append(csvUsers, csvUser)
	}

	w := csv.NewWriter(file)
	defer w.Flush()
	err = w.WriteAll(csvUsers)
	if err != nil {
		return false, err
	}

	log.Println(fmt.Sprintf("%v was updated", dataFileName))

	return true, nil
}

func (ts *readerService) Get(dataFileName string) (*[]types.User, error) {
	records, err := ts.readCsvFromFile(dataFileName)
	if err != nil {
		return nil, err
	}

	var users []types.User

	for _, line := range records {
		id, err := strconv.Atoi(line[0])

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Id '%v' is not a number", line[0]))
		}

		users = append(users, types.User{
			Id:         id,
			Username:   line[1],
			Identifier: line[2],
			FirstName:  line[3],
			LastName:   line[4],
		})
	}

	return &users, nil
}

func (ts *readerService) GetConcurrently(dataFileName string, idType string, items int, itemsPerWorker int) (*[]types.User, error) {
	if !contains([]string{"odd", "even"}, idType) {
		idType = "both"
	}

	file, err := os.Open(dataFileName)
	if err != nil {
		return nil, errors.New("can not open file")
	}
	reader := csv.NewReader(file)

	numOfCors := int(math.Min(float64(runtime.NumCPU()), float64(items)))
	runtime.GOMAXPROCS(numOfCors)

	resultChannel := make(chan types.User)
	endChannel := make(chan int)

	itemsPerWorker = int(math.Min(float64(itemsPerWorker), float64(items/numOfCors)))
	for i := 0; i < numOfCors; i++ {
		go process(resultChannel, endChannel, idType, reader, itemsPerWorker)
	}

	var users []types.User

	for i := 0; i < itemsPerWorker*numOfCors; i++ {
		select {
		case user := <-resultChannel:
			users = append(users, user)
		case left := <-endChannel:
			i += left - 1
		}
	}

	close(resultChannel)
	close(endChannel)

	return &users, nil
}

func process(resultChannel chan<- types.User, endChannel chan<- int, idType string, reader *csv.Reader, itemsPerWorker int) {
	var mu sync.Mutex
	for i := 0; i < itemsPerWorker; {
		mu.Lock()
		record, err := reader.Read()
		if err != nil {
			fmt.Printf("read file error: %v", err)
			endChannel <- itemsPerWorker - i
			break
		}
		mu.Unlock()

		id, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Printf("Id '%v' is not a number", record[0])
			continue
		}

		if idType != "both" {
			var userIdType string
			if id%2 == 0 {
				userIdType = "even"
			} else {
				userIdType = "odd"
			}

			if idType != userIdType {
				continue
			}
		}

		user := types.User{
			Id:         id,
			Username:   record[1],
			Identifier: record[2],
			FirstName:  record[3],
			LastName:   record[4],
		}

		resultChannel <- user
		i++
	}
}

func (ts *readerService) GetMap(dataFileName string) (map[int]types.User, error) {
	records, err := ts.readCsvFromFile(dataFileName)
	if err != nil {
		return nil, err
	}

	users := make(map[int]types.User, len(records))

	for _, line := range records {
		id, err := strconv.Atoi(line[0])

		if err != nil {
			return nil, errors.New(fmt.Sprintf("Id '%v' is not a number", line[0]))
		}

		users[id] = types.User{
			Id:         id,
			Username:   line[1],
			Identifier: line[2],
			FirstName:  line[3],
			LastName:   line[4],
		}
	}

	return users, nil
}

func (ts *readerService) GetFeed(feedFileName string) ([][]string, error) {
	return ts.readCsvFromFile(feedFileName)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

package usecase

import (
	"errors"
	"fmt"

	"Go-Dispatch-Bootcamp/types"
)

const apiUrl = "http://localhost:8080/api/v1/feed"
const dataFileName = "data/data.csv"
const feedFileName = "data/feed.csv"

type readerService interface {
	Get(string) (*[]types.User, error)
	GetConcurrently(string, string, int, int) (*[]types.User, error)
	GetMap(string) (map[int]types.User, error)
	GetFeed(string) ([][]string, error)
	FetchCsvFromRemote(string) ([][]string, error)
	Update(*[]types.User, string) (bool, error)
}

type readerUsecase struct {
	service readerService
}

func New(s readerService) *readerUsecase {
	return &readerUsecase{
		service: s,
	}
}

func (tu *readerUsecase) Fetch() (*[]types.User, error) {
	users, err := tu.service.Get(dataFileName)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (tu *readerUsecase) FetchConcurrently(idType string, items int, itemsPerWorker int) (*[]types.User, error) {
	users, err := tu.service.GetConcurrently(dataFileName, idType, items, itemsPerWorker)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (tu *readerUsecase) FetchById(id int) (*types.User, error) {
	users, err := tu.service.GetMap(dataFileName)

	if err != nil {
		return nil, err
	}

	result, ok := users[id]

	if !ok {
		return nil, errors.New(fmt.Sprintf("User with id: %v doesn't exist", id))
	}

	return &result, nil
}

func (tu *readerUsecase) Feed() ([][]string, error) {
	return tu.service.GetFeed(feedFileName)
}

func (tu *readerUsecase) UpdateUsersFromFeed() (bool, error) {
	csvUsers, err := tu.service.FetchCsvFromRemote(apiUrl)

	if err != nil {
		return false, err
	}

	fmt.Println(csvUsers)
	var users []types.User

	for i, csvUser := range csvUsers {
		// skipping title
		if i == 0 {
			continue
		}

		user := types.User{
			Id:         i,
			Username:   csvUser[0],
			Identifier: csvUser[1],
			FirstName:  csvUser[2],
			LastName:   csvUser[3],
		}

		users = append(users, user)
	}

	return tu.service.Update(&users, dataFileName)
}

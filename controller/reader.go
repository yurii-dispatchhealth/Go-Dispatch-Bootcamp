package controller

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Go-Dispatch-Bootcamp/types"

	"github.com/gorilla/mux"
)

type usecase interface {
	Fetch() (*[]types.User, error)
	FetchConcurrently(string, int, int) (*[]types.User, error)
	FetchById(int) (*types.User, error)
	Feed() ([][]string, error)
	UpdateUsersFromFeed() (bool, error)
}

type readerController struct {
	usecase usecase
}

func (tc *readerController) Fetch(w http.ResponseWriter, r *http.Request) {
	users, err := tc.usecase.Fetch()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Fetch error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	data, _ := json.Marshal(users)
	w.Write(data)
}

func (tc *readerController) FetchConcurrently(w http.ResponseWriter, r *http.Request) {
	idType := r.URL.Query().Get("type")

	items, err := strconv.Atoi(r.URL.Query().Get("items"))
	if err != nil {
		fmt.Println("Items parameter is invalid")
		items = 2
	}

	itemsPerWorkers, err := strconv.Atoi(r.URL.Query().Get("items_per_workers"))
	if err != nil {
		fmt.Println("ItemsPerWorkers parameter is invalid")
		itemsPerWorkers = 2
	}

	users, err := tc.usecase.FetchConcurrently(idType, items, itemsPerWorkers)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "FetchConcurrently error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	data, _ := json.Marshal(users)
	w.Write(data)
}

func (tc *readerController) FetchById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringId, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "FetchById error: id is not defined")
		return
	}

	id, err := strconv.Atoi(stringId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("FetchById error: Id '%v' is not a number", stringId))
		return
	}

	user, err := tc.usecase.FetchById(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("FetchById error: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	data, _ := json.Marshal(user)
	w.Write(data)
}

func (tc *readerController) Feed(w http.ResponseWriter, r *http.Request) {
	users, err := tc.usecase.Feed()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Feed error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	b := new(bytes.Buffer)
	csvWriter := csv.NewWriter(b)
	csvWriter.WriteAll(users)
	w.Write(b.Bytes())
}

func (tc *readerController) UpdateUsersFromFeed(w http.ResponseWriter, r *http.Request) {
	success, err := tc.usecase.UpdateUsersFromFeed()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "UpdateUsersFromFeed error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{ success: { %v } }", success)))
}

func New(uc usecase) *readerController {
	return &readerController{
		usecase: uc,
	}
}

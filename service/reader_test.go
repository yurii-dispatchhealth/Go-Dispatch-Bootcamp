package service

import (
	"log"
	"reflect"
	"strings"
	"testing"

	"Go-Dispatch-Bootcamp/types"

	"github.com/stretchr/testify/assert"
)

var (
	fileData = [][]string{
		{"1", "booker12", "9012", "Rachel", "Booker"},
		{"2", "grey07", "2070", "Laura", "Grey"},
		{"3", "johnson81", "4081", "Craig", "Johnson"},
		{"4", "jenkins46", "9346", "Mary", "Jenkins"},
		{"5", "smith79", "5079", "Jamie", "Smith"},
	}

	feedFileData = [][]string{
		{"username", "id", "first_name", "last_name"},
		{"booker12", "9012", "Rachel", "Booker"},
		{"grey07", "2070", "Laura", "Grey"},
		{"johnson81", "4081", "Craig", "Johnson"},
		{"jenkins46", "9346", "Mary", "Jenkins"},
		{"smith79", "5079", "Jamie", "Smith"},
	}

	users = []types.User{
		{1, "booker12", "9012", "Rachel", "Booker"},
		{2, "grey07", "2070", "Laura", "Grey"},
		{3, "johnson81", "4081", "Craig", "Johnson"},
		{4, "jenkins46", "9346", "Mary", "Jenkins"},
		{5, "smith79", "5079", "Jamie", "Smith"},
	}

	usersMap = map[int]types.User{
		1: {1, "booker12", "9012", "Rachel", "Booker"},
		2: {2, "grey07", "2070", "Laura", "Grey"},
		3: {3, "johnson81", "4081", "Craig", "Johnson"},
		4: {4, "jenkins46", "9346", "Mary", "Jenkins"},
		5: {5, "smith79", "5079", "Jamie", "Smith"},
	}

	feedUrl      = "http://localhost:8080/api/v1/feed"
	dataFileName = "../data/data.csv"
	feedFileName = "../data/feed.csv"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *readerService
	}{
		{
			name: "New controller test",
			want: &readerService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readerService_readCsvFromFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		ts   *readerService
		want [][]string
	}{
		{
			name: "Read csv from file. Success story.",
			args: args{path: dataFileName},
			ts:   &readerService{},
			want: fileData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.ts.readCsvFromFile(tt.args.path)

			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

func Test_readerService_Update(t *testing.T) {
	type args struct {
		users *[]types.User
	}
	tests := []struct {
		name string
		ts   *readerService
		args args
		want bool
	}{
		{
			name: "Update users. Success story.",
			args: args{users: &users},
			ts:   &readerService{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.ts.Update(tt.args.users, dataFileName)

			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

func Test_readerService_Get(t *testing.T) {
	tests := []struct {
		name string
		ts   *readerService
		want *[]types.User
	}{
		{
			name: "Get users. Success story.",
			ts:   &readerService{},
			want: &users,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.ts.Get(dataFileName)

			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

func Test_readerService_GetMap(t *testing.T) {
	tests := []struct {
		name string
		ts   *readerService
		want map[int]types.User
	}{
		{
			name: "Get users map. Success story.",
			ts:   &readerService{},
			want: usersMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.ts.GetMap(dataFileName)

			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

func Test_readerService_GetFeed(t *testing.T) {
	tests := []struct {
		name string
		ts   *readerService
		want [][]string
	}{
		{
			name: "Get users feed. Success story.",
			ts:   &readerService{},
			want: feedFileData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.ts.GetFeed(feedFileName)

			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

// To make this test pass you need to run the server
func Test_readerService_FetchCsvFromRemote(t *testing.T) {
	tests := []struct {
		name string
		ts   *readerService
		want [][]string
	}{
		{
			name: "Fetch csv from remote. Success story.",
			ts:   &readerService{},
			want: feedFileData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.ts.FetchCsvFromRemote(feedUrl)

			if err != nil && strings.Contains(err.Error(), "connection refused") {
				log.Printf("error: %v", err)
				log.Fatalf("Make sure the server is up on port 8080")
			}

			assert.Equal(t, tt.want, result)
			assert.NoError(t, err)
		})
	}
}

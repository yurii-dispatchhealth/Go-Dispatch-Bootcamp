package router

import (
	"net/http"
	"testing"

	"Go-Dispatch-Bootcamp/mocks"

	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	type args struct {
		c readerController
	}
	tests := []struct {
		name      string
		routeName string
		path      string
		methods   []string
	}{
		{
			name:      "Fetch route",
			routeName: "Fetch",
			path:      "api/v1/fetch",
			methods:   []string{http.MethodGet},
		},
		{
			name:      "FetchById route",
			routeName: "FetchById",
			path:      "api/v1/fetch/{id}",
			methods:   []string{http.MethodGet},
		},
		{
			name:      "Feed route",
			routeName: "Feed",
			path:      "api/v1/feed",
			methods:   []string{http.MethodGet},
		},
		{
			name:      "UpdateUsersFromFeed route",
			routeName: "UpdateUsersFromFeed",
			path:      "api/v1/run-update-users-from-feed",
			methods:   []string{http.MethodGet},
		},
	}

	r := Setup(&mocks.ReaderController{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			route := r.Get(tt.routeName)

			assert.NotNil(t, route)
			assert.Equal(t, route.GetName(), tt.routeName)

			methods, _ := route.GetMethods()
			assert.Equal(t, methods, tt.methods)
		})
	}
}

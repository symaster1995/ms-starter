package http_test

import (
	"bytes"
	"context"
	"github.com/symaster1995/ms-starter/internal/models"
	"net/http"
	"testing"
	"time"
)

func TestDialIndex(t *testing.T) {
	// Start the mocked HTTP test server.
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	ctx0 := context.Background()

	//Mock data
	item := &models.Item{
		ID:        69,
		Name:      "nai sas",
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	//Mock Get Items function
	s.ItemService.FindItemsFn = func(ctx context.Context, filter models.ItemFilter) ([]*models.Item, int, error) {
		return []*models.Item{item}, 1, nil
	}



	t.Run("GET ITEMS", func(t *testing.T) {
		var filter = []byte(`{"name":"", "limit": 0, "offset": 0}`)

		resp, err := http.DefaultClient.Do(s.MustNewRequest(t, ctx0, "GET", "/items", bytes.NewBuffer(filter)))
		if err != nil {
			t.Fatal(err)
		} else if got, want := resp.StatusCode, http.StatusOK; got != want {
			t.Fatalf("StatusCode=%v, want %v", got, want)
		}
	})
}

package http_test

import (
	"context"
	"github.com/google/go-cmp/cmp"
	productsModel "github.com/symaster1995/ms-starter/internal/products/models"
	"testing"
	"time"
)

func TestItemIndex(t *testing.T) {
	// Start the mocked HTTP test server.
	s := MustOpenServer(t)
	defer MustCloseServer(t, s)

	ctx0 := context.Background()

	//Mock data
	item := &productsModel.Item{
		ID:        69,
		Name:      "tony",
		CreatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	//Mock Get Items function
	s.ItemService.FindItemsFn = func(ctx context.Context, filter productsModel.ItemFilter) ([]*productsModel.Item, int, error) {
		if filter.Name != nil && *filter.Name != item.Name {
			return []*productsModel.Item{}, 0, nil
		}
		return []*productsModel.Item{item}, 1, nil
	}

	type args struct {
		Name   string
		Limit  int
		Offset int
	}

	type wants struct {
		err   error
		items []*productsModel.Item
	}

	//collection of tests
	tests := []struct {
		name string
		wants
		args
	}{
		{
			name: "find_all_items",
			wants: wants{
				items: []*productsModel.Item{item},
			},
			args: args{},
		},
		{
			name: "find_non_existing_item_filtered",
			wants: wants{
				items: []*productsModel.Item{},
			},
			args: args{
				Name: "parker",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filter = productsModel.ItemFilter{Name: nil, Limit: 0, Offset: 0}

			if tt.args.Name != "" {
				filter.Name = &tt.args.Name
			}

			items, _, err := s.ItemService.FindItems(ctx0, filter)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(items, tt.wants.items); diff != "" {
				t.Errorf("items are different -got/+want\ndiff %s", diff)
			}
		})
	}
}

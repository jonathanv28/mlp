package pagination

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginationOptions(t *testing.T) {
	one, two, three, ten := int32(1), int32(2), int32(3), int32(10)
	paginator := NewPaginator(1, 10, 50)

	tests := map[string]struct {
		page     *int32
		pageSize *int32
		expected Options
	}{
		"defaults": {
			expected: Options{
				Page:     &one,
				PageSize: &ten,
			},
		},
		"missing page": {
			pageSize: &three,
			expected: Options{
				Page:     &one,
				PageSize: &three,
			},
		},
		"missing page size": {
			page: &two,
			expected: Options{
				Page:     &two,
				PageSize: &ten,
			},
		},
		"new values": {
			page:     &two,
			pageSize: &three,
			expected: Options{
				Page:     &two,
				PageSize: &three,
			},
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, data.expected, paginator.NewPaginationOptions(data.page, data.pageSize))
		})
	}
}

func TestValidatePaginationParams(t *testing.T) {
	zero, one, two, hundred := int32(0), int32(1), int32(2), int32(100)
	paginator := NewPaginator(1, 10, 50)

	tests := map[string]struct {
		page          *int32
		pageSize      *int32
		expectedError string
	}{
		"zero page size": {
			pageSize:      &zero,
			expectedError: "page size must be within range (0 < page_size <= 50) or unset",
		},
		"large page size": {
			pageSize:      &hundred,
			expectedError: "page size must be within range (0 < page_size <= 50) or unset",
		},
		"zero page": {
			page:          &zero,
			expectedError: "page must be > 0 or unset",
		},
		"success - all values unset": {},
		"success - all values set": {
			page:     &one,
			pageSize: &two,
		},
	}

	for name, data := range tests {
		t.Run(name, func(t *testing.T) {
			err := paginator.ValidatePaginationParams(data.page, data.pageSize)
			if data.expectedError == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, data.expectedError)
			}
		})
	}
}

func TestPaging(t *testing.T) {
	tests := []struct {
		page     int32
		pageSize int32
		count    int
		expected Paging
	}{
		{
			page:     int32(1),
			pageSize: int32(10),
			count:    5,
			expected: Paging{
				Page:  1,
				Pages: 1,
				Total: 5,
			},
		},
		{
			page:     int32(2),
			pageSize: int32(3),
			count:    7,
			expected: Paging{
				Page:  2,
				Pages: 3,
				Total: 7,
			},
		},
	}

	for idx, data := range tests {
		t.Run(fmt.Sprintf("case %d", idx), func(t *testing.T) {
			actual := ToPaging(Options{Page: &data.page, PageSize: &data.pageSize}, data.count)
			assert.Equal(t, data.expected, *actual)
		})
	}
}

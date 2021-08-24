package service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alexkaplun/books-test/service/api"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8080/books"
)

var (
	client = http.DefaultClient
	book   = &api.Book{
		Title:       "1",
		Author:      "2",
		Publisher:   "3",
		PublishDate: "2021-07-07",
		Rating:      2,
		Status:      "CheckedOut",
	}
)

// Since we are running tests against dockerized app, we are not able to use httptest package and
// test handlers separately
//
// We may also want to set up test data prior to running test and clean up after

func TestCreateBook(t *testing.T) {
	cases := map[string]struct {
		payload      string
		expectedCode int
	}{
		"bad json": {
			payload: `{ aa
						"author": "some author",
						"publisher": "some publisher",
						"rating": 2,
						"status": "CheckedOut"
					}`,
			expectedCode: http.StatusBadRequest,
		},
		"no title": {
			payload: `{
						"author": "some author",
						"publisher": "some publisher",
						"rating": 2,
						"status": "CheckedOut"
					}`,
			expectedCode: http.StatusBadRequest,
		},
		"no author": {
			payload: `{
						"title": "some title",
						"publisher": "some publisher",
						"rating": 2,
						"status": "CheckedOut"
					}`,
			expectedCode: http.StatusBadRequest,
		},
		"rating out of range": {
			payload: `{
						"title": "some title",
						"author": "some author",
						"publisher": "some publisher",
						"rating": 4,
						"status": "CheckedOut"
					}`,
			expectedCode: http.StatusBadRequest,
		},
		"invalid status": {
			payload: `{
						"title": "some title",
						"author": "some author",
						"publisher": "some publisher",
						"rating": 1,
						"status": "NotCheckedIn"
					}`,
			expectedCode: http.StatusBadRequest,
		},
		"valid payload": {
			payload: `{
						"title": "some title",
						"author": "some author",
						"publisher": "some publisher",
						"rating": 1,
						"status": "CheckedIn"
					}`,
			expectedCode: http.StatusOK,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(test.payload))
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, resp.StatusCode)
			if test.expectedCode == http.StatusOK {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				// unmarshal response
				var createResp api.CreateBookResponse
				require.NoError(t, json.Unmarshal(body, &createResp))
				assert.NotNil(t, createResp.ID)
			}
		})
	}
}

func TestGetBook(t *testing.T) {
	// create a book first
	id, err := createBook(book)
	require.NoError(t, err)

	cases := map[string]struct {
		bookID       string
		expectedCode int
		expectedData *api.Book
	}{
		"bad id": {
			bookID:       "i am bad id",
			expectedCode: http.StatusBadRequest,
		},
		"missing id": {
			bookID:       uuid.New().String(),
			expectedCode: http.StatusNotFound,
		},
		"valid": {
			bookID:       id.String(),
			expectedCode: http.StatusOK,
			expectedData: book,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			url := fmt.Sprintf("%s/%s", baseURL, test.bookID)
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, resp.StatusCode)
			if test.expectedCode == http.StatusOK {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)

				// unmarshal response
				var getBook api.Book
				require.NoError(t, json.Unmarshal(body, &getBook))

				if test.expectedData == nil {
					t.Fatal("test expectedData should not be empty")
				}

				d := test.expectedData
				assert.Equal(t, test.bookID, getBook.ID.String())
				assert.Equal(t, d.Title, getBook.Title)
				assert.Equal(t, d.Author, getBook.Author)
				assert.Equal(t, d.Publisher, getBook.Publisher)
				assert.Equal(t, d.PublishDate, getBook.PublishDate)
				assert.Equal(t, d.Rating, getBook.Rating)
				assert.Equal(t, d.Status, getBook.Status)
				assert.NotEmpty(t, getBook.CreatedAt)
				assert.NotEmpty(t, getBook.UpdatedAt)
			}
		})
	}
}

func TestListBooks(t *testing.T) {
	// create a book first to make sure at least one object exists
	_, err := createBook(book)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	// unmarshal response
	var getResp []*api.Book
	require.NoError(t, json.Unmarshal(body, &getResp))

	assert.NotEmpty(t, getResp)
}

func TestDeleteBook(t *testing.T) {
	// create a book first
	id, err := createBook(book)
	require.NoError(t, err)

	cases := map[string]struct {
		bookID       string
		expectedCode int
	}{
		"bad id": {
			bookID:       "i am bad id",
			expectedCode: http.StatusBadRequest,
		},
		// we expect 200 on missing id as DELETE method is implemented as an idempotent operation
		"missing id": {
			bookID:       uuid.New().String(),
			expectedCode: http.StatusOK,
		},
		// in this test we don't check that the book was actually deleted
		"valid": {
			bookID:       id.String(),
			expectedCode: http.StatusOK,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			url := fmt.Sprintf("%s/%s", baseURL, test.bookID)
			req, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, resp.StatusCode)
		})
	}
}

// we don't check all the 400 cases here - they are the same as in TestCreateBook
func TestUpdateBook(t *testing.T) {
	// create a book that will be updated
	id, err := createBook(book)
	require.NoError(t, err)

	// sleep to make sure updated_at will change as it's precision is seconds
	time.Sleep(1 * time.Second)

	updatedBook := &api.UpsertBookRequest{
		Title:       "s1 ome title",
		Author:      "s1 ome author",
		Publisher:   "some publisher",
		PublishDate: "2017-01-12",
		Rating:      3,
		Status:      "CheckedIn",
	}

	// marshall updatedBook payload
	var payload []byte
	payload, err = json.Marshal(updatedBook)
	require.NoError(t, err)

	cases := map[string]struct {
		bookID       string
		expectedCode int
		payload      string
	}{
		"bad id": {
			bookID:       "i am bad id",
			expectedCode: http.StatusBadRequest,
			payload:      string(payload),
		},
		"missing id": {
			bookID:       uuid.New().String(),
			expectedCode: http.StatusNotFound,
			payload:      string(payload),
		},
		"valid": {
			bookID:       id.String(),
			expectedCode: http.StatusOK,
			payload:      string(payload),
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {

			url := fmt.Sprintf("%s/%s", baseURL, test.bookID)
			req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(test.payload))
			require.NoError(t, err)
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, test.expectedCode, resp.StatusCode)
			if test.expectedCode == http.StatusOK {
				// get the updated book and check all data was updated

				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				getResp, err := client.Do(req)
				require.NoError(t, err)

				defer getResp.Body.Close()
				body, err := ioutil.ReadAll(getResp.Body)
				require.NoError(t, err)

				var getBook api.Book
				require.NoError(t, json.Unmarshal(body, &getBook))

				assert.Equal(t, test.bookID, getBook.ID.String())
				assert.Equal(t, updatedBook.Title, getBook.Title)
				assert.Equal(t, updatedBook.Author, getBook.Author)
				assert.Equal(t, updatedBook.Publisher, getBook.Publisher)
				assert.Equal(t, updatedBook.PublishDate, getBook.PublishDate)
				assert.Equal(t, updatedBook.Rating, getBook.Rating)
				assert.Equal(t, updatedBook.Status, getBook.Status)
				assert.NotEmpty(t, getBook.CreatedAt)
				assert.NotEmpty(t, getBook.UpdatedAt)
				assert.NotEqual(t, getBook.CreatedAt, getBook.UpdatedAt)
			}
		})
	}
}

func createBook(book *api.Book) (*uuid.UUID, error) {
	upsertBookRequest := &api.UpsertBookRequest{
		Title:       book.Title,
		Author:      book.Author,
		Publisher:   book.Publisher,
		PublishDate: book.PublishDate,
		Rating:      book.Rating,
		Status:      book.Status,
	}

	payload, err := json.Marshal(upsertBookRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got non 200 status: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var createResp api.CreateBookResponse
	if err := json.Unmarshal(body, &createResp); err != nil {
		return nil, err
	}

	return createResp.ID, nil
}

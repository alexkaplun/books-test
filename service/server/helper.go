package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/alexkaplun/books-test/service/api"
	"github.com/alexkaplun/books-test/storage/models"
)

func jsonOK(w http.ResponseWriter, resp interface{}) {
	var (
		payload []byte
		err     error
	)
	if resp == nil {
		payload = []byte("OK")
	} else {
		payload, err = json.Marshal(resp)
		if err != nil {
			log.Printf("failed to marshal response body. err: %v\n", err)
			http.Error(w, "failed to marshal response body", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func parseBody(body io.ReadCloser, dest interface{}) error {
	defer body.Close()
	if body == http.NoBody {
		return errors.New("empty body")
	}

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return err
	}

	return nil
}

func convertBookToDB(in *api.UpsertBookRequest) (*models.Book, error) {
	book := &models.Book{
		Title:     in.Title,
		Author:    in.Author,
		Publisher: in.Publisher,
		Rating:    in.Rating,
		Status:    models.BookStatus(in.Status),
	}

	if len(in.PublishDate) != 0 {
		publishDate, err := time.Parse("2006-01-02", in.PublishDate)
		if err != nil {
			return nil, err
		}

		book.PublishDate = &publishDate
	}

	return book, nil
}

func convertBookFromDB(in *models.Book) *api.Book {
	book := &api.Book{
		ID:        in.ID,
		Title:     in.Title,
		Author:    in.Author,
		Publisher: in.Publisher,
		Rating:    in.Rating,
		Status:    string(in.Status),
		CreatedAt: in.CreatedAt.Format(time.RFC3339),
		UpdatedAt: in.UpdatedAt.Format(time.RFC3339),
	}

	if in.PublishDate != nil {
		book.PublishDate = in.PublishDate.Format("2006-01-02")
	}

	return book
}

func convertBooksFromDB(in []*models.Book) []*api.Book {
	books := make([]*api.Book, len(in))
	for i, v := range in {
		books[i] = convertBookFromDB(v)
	}
	return books
}

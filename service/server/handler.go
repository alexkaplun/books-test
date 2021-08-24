package server

import (
	"log"
	"net/http"

	"github.com/alexkaplun/books-test/storage/models"

	"github.com/google/uuid"

	"github.com/alexkaplun/books-test/service/api"
	"github.com/alexkaplun/books-test/storage"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	storage storage.Storage
}

type HandlerParams struct {
	Storage storage.Storage
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		storage: params.Storage,
	}
}

func (h *Handler) createBookHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer h.guardPanic()

	var req api.UpsertBookRequest
	if err := parseBody(r.Body, &req); err != nil {
		log.Printf("failed to parse request body. err: %v\n", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("failed to validate request. err: %v\n", err)
		http.Error(w, "failed to validate request", http.StatusBadRequest)
		return
	}

	book, err := convertBookToDB(&req)
	if err != nil {
		log.Printf("failed to convert book request to DB. err: %v\n", err)
		http.Error(w, "failed to convert book request to DB", http.StatusBadRequest)
		return
	}

	id, err := h.storage.CreateBook(r.Context(), book)
	if err != nil {
		log.Printf("failed to save book to DB. err: %v\n", err)
		http.Error(w, "failed to save book to DB", http.StatusInternalServerError)
		return
	}

	resp := api.CreateBookResponse{
		ID: id,
	}

	jsonOK(w, &resp)
}

func (h *Handler) deleteBookHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer h.guardPanic()

	bookIDStr := p.ByName("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		log.Printf("failed to parse book id. err: %v\n", err)
		http.Error(w, "failed to parse book id", http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteBook(r.Context(), bookID); err != nil {
		log.Printf("failed to delete book. err: %v\n", err)
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	jsonOK(w, nil)
}

func (h *Handler) updateBookHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer h.guardPanic()

	bookIDStr := p.ByName("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		log.Printf("failed to parse book id. err: %v\n", err)
		http.Error(w, "failed to parse book id", http.StatusBadRequest)
		return
	}

	var req api.UpsertBookRequest
	if err = parseBody(r.Body, &req); err != nil {
		log.Printf("failed to parse request body. err: %v\n", err)
		http.Error(w, "failed to parse request body", http.StatusBadRequest)
		return
	}

	if err = req.Validate(); err != nil {
		log.Printf("failed to validate request. err: %v\n", err)
		http.Error(w, "failed to validate request", http.StatusBadRequest)
		return
	}

	var book *models.Book
	book, err = convertBookToDB(&req)
	if err != nil {
		log.Printf("failed to convert book request to DB. err: %v\n", err)
		http.Error(w, "failed to convert book request to DB", http.StatusBadRequest)
		return
	}

	book.ID = bookID

	if err = h.storage.UpdateBook(r.Context(), book); err != nil {
		log.Printf("failed to update book. err: %v\n", err)
		if err == storage.ErrBookNotFound {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	jsonOK(w, nil)
}

func (h *Handler) getBookHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	defer h.guardPanic()

	bookIDStr := p.ByName("id")
	bookID, err := uuid.Parse(bookIDStr)
	if err != nil {
		log.Printf("failed to parse book id. err: %v\n", err)
		http.Error(w, "failed to parse book id", http.StatusBadRequest)
		return
	}

	book, err := h.storage.GetBook(r.Context(), bookID)
	if err != nil {
		log.Printf("failed to find book. err: %v\n", err)
		if err == storage.ErrBookNotFound {
			http.Error(w, "book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to find book", http.StatusInternalServerError)
		return
	}

	jsonOK(w, convertBookFromDB(book))
}

// TODO: paging
func (h *Handler) listBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer h.guardPanic()

	books, err := h.storage.ListBooks(r.Context())
	if err != nil {
		log.Printf("failed to list book. err: %v\n", err)
		http.Error(w, "failed to list books", http.StatusInternalServerError)
		return
	}

	jsonOK(w, convertBooksFromDB(books))
}

func (h *Handler) guardPanic() {
	if p := recover(); p != nil {
		log.Printf("caught panic: %v\n", p)
	}
}

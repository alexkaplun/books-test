package server

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	http.Handler
}

type RouterParams struct {
	Handler *Handler
}

func NewRouter(params RouterParams) *Router {
	router := httprouter.New()
	router.PanicHandler = panicHandler
	h := params.Handler

	router.POST("/books", h.createBookHandler)
	router.DELETE("/books/:id", h.deleteBookHandler)
	router.PUT("/books/:id", h.updateBookHandler)
	router.GET("/books/:id", h.getBookHandler)
	router.GET("/books", h.listBooks)

	return &Router{
		Handler: router,
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	log.Println(r.URL.Path, err)
	debug.PrintStack()
	w.WriteHeader(http.StatusInternalServerError)
}

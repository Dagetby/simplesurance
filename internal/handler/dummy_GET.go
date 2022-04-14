package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type counter interface {
	CountRequests(requestTime time.Time) (int, error)
}

type Handler struct {
	counter counter
}

func New(c counter) *Handler {
	return &Handler{counter: c}
}

func (h *Handler) DummyHandler(w http.ResponseWriter, _ *http.Request) {
	count, err := h.counter.CountRequests(time.Now())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Sorry, we have trouble( \n Try later!"))
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(strconv.Itoa(count)))
	if err != nil {
		log.Println(err)
	}
}

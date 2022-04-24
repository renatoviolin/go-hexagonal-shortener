package chiAdapter

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/renatoviolin/shortener/adapter/serializer"
	"github.com/renatoviolin/shortener/application/entity"
	"github.com/renatoviolin/shortener/application/shortener"
	"github.com/renatoviolin/shortener/ports"
)

type ChiHandler struct {
	UseCase   shortener.UseCaseShortener
	ChiRouter *chi.Mux
}

func NewChiHandler(useCase shortener.UseCaseShortener) *ChiHandler {
	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.Logger)
	return &ChiHandler{
		UseCase:   useCase,
		ChiRouter: chiRouter,
	}
}

func (h *ChiHandler) Run(address string) {
	fmt.Printf("CHI  listening on %s\n", address)
	log.Fatal(http.ListenAndServe(address, h.ChiRouter))
}

func (h *ChiHandler) SetupRoutes() {
	h.ChiRouter.Get("/{code}", func(w http.ResponseWriter, r *http.Request) {
		handlerGet(w, r, &h.UseCase)
	})
	h.ChiRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		handlerPost(w, r, &h.UseCase)
	})
}

func handlerGet(w http.ResponseWriter, r *http.Request, useCase *shortener.UseCaseShortener) {
	code := chi.URLParam(r, "code")
	redirect, err := useCase.CodeToUrl(code)
	if err != nil {
		if err == entity.ErrRedirectNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

func handlerPost(w http.ResponseWriter, r *http.Request, useCase *shortener.UseCaseShortener) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	redirectIn, err := getSerializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	redirectOut, err := useCase.UrlToCode(redirectIn.URL)
	if err != nil {
		if err == entity.ErrRedirectInvalid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := getSerializer(contentType).Encode(redirectOut)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	returnResponse(w, contentType, responseBody, http.StatusCreated)
}

func getSerializer(contentType string) ports.RedirectSerializer {
	if contentType == "application/x-msgpack" {
		return &serializer.MsgPackSerializer{}
	} else {
		return &serializer.JSONSerializer{}
	}
}

func returnResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

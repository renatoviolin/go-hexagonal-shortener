package ginAdapter

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renatoviolin/shortener/application/entity"
	"github.com/renatoviolin/shortener/application/shortener"
)

type GinHandler struct {
	UseCase   shortener.UseCaseShortener
	GinRouter *gin.Engine
}

func NewGinHandler(useCase shortener.UseCaseShortener) *GinHandler {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	return &GinHandler{
		UseCase:   useCase,
		GinRouter: ginRouter,
	}
}

func (h *GinHandler) Run(address string) {
	fmt.Printf("GIN  listening on %s\n", address)
	log.Fatal(http.ListenAndServe(address, h.GinRouter))
}

func (h *GinHandler) SetupRoutes() {
	h.GinRouter.GET("/:code", func(ctx *gin.Context) {
		handlerGet(ctx, h.UseCase)
	})
	h.GinRouter.POST("/", func(ctx *gin.Context) {
		handlerPost(ctx, h.UseCase)
	})
}

func handlerGet(ctx *gin.Context, useCase shortener.UseCaseShortener) {
	code := ctx.Param("code")
	redirect, err := useCase.CodeToUrl(code)
	if err != nil {
		if err == entity.ErrRedirectNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.Redirect(http.StatusPermanentRedirect, redirect.URL)
}

type inputRequest struct {
	Url string `json:"url" binding:"required"`
}

func handlerPost(ctx *gin.Context, useCase shortener.UseCaseShortener) {
	var input = inputRequest{}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	redirect, err := useCase.UrlToCode(input.Url)
	if err != nil {
		if err == entity.ErrRedirectInvalid {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, redirect)
}

package shortener

import (
	"github.com/renatoviolin/shortener/application/entity"
)

type UseCaseShortener struct {
	services *RedirectService
}

func NewUseCaseShortener(redirectServices *RedirectService) *UseCaseShortener {
	return &UseCaseShortener{
		services: redirectServices,
	}
}

func (uc *UseCaseShortener) UrlToCode(url string) (*entity.Redirect, error) {
	return uc.services.Store(url)
}

func (uc *UseCaseShortener) CodeToUrl(code string) (*entity.Redirect, error) {
	return uc.services.Find(code)
}

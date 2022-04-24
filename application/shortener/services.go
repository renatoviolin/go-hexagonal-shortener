package shortener

import (
	"time"

	"github.com/renatoviolin/shortener/application/entity"
	"github.com/renatoviolin/shortener/ports"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

type RedirectService struct {
	repository ports.RedirectRepository
}

func NewRedirectService(redirectRepo ports.RedirectRepository) *RedirectService {
	return &RedirectService{
		repository: redirectRepo,
	}
}

func (uc *RedirectService) Find(code string) (*entity.Redirect, error) {
	return uc.repository.Find(code)
}

func (uc *RedirectService) Store(url string) (*entity.Redirect, error) {
	redirect := &entity.Redirect{
		URL: url,
	}
	if err := validate.Validate(redirect); err != nil {
		return nil, entity.ErrRedirectInvalid
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	err := uc.repository.Store(redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}

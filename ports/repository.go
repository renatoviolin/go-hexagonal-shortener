package ports

import "github.com/renatoviolin/shortener/application/entity"

type RedirectRepository interface {
	Find(code string) (*entity.Redirect, error)
	Store(redirect *entity.Redirect) error
}

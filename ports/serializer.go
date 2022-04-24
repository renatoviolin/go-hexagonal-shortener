package ports

import "github.com/renatoviolin/shortener/application/entity"

type RedirectSerializer interface {
	Decode(input []byte) (*entity.Redirect, error)
	Encode(input *entity.Redirect) ([]byte, error)
}

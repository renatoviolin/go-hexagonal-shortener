package serializer

import (
	"github.com/renatoviolin/shortener/application/entity"
	"github.com/vmihailenco/msgpack/v5"
)

type MsgPackSerializer struct{}

func (m *MsgPackSerializer) Decode(input []byte) (*entity.Redirect, error) {
	redirect := &entity.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, err
	}

	return redirect, nil
}

func (m *MsgPackSerializer) Encode(input *entity.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, err
	}

	return rawMsg, nil
}

package id

import (
	"context"

	"github.com/speps/go-hashids"
	"ymir.com/pkg/xerr"
)

var hash *HashID

func init() {
	hash = NewHashID()
}

type HashID struct {
	hash *hashids.HashID
}

func NewHashID() *HashID {
	return &HashID{
		hash: NewHashIds(),
	}
}

func NewHashIds() *hashids.HashID {
	var hd = hashids.NewData()
	hd.Salt = "123poiasdmnmb"
	hd.MinLength = 8
	hd.Alphabet = "qwertyuiopasdfghjklzxcvbnm"
	h, _ := hashids.NewWithData(hd)
	return h
}

func (h *HashID) EncodedId(id int64) (string, error) {
	encodedId, err := h.hash.EncodeInt64([]int64{id})
	if err != nil {
		return "", err
	}

	return encodedId, nil
}

func (h *HashID) DecodedId(id string) int64 {
	return h.hash.DecodeInt64(id)[0]
}

func GetDecodedUserId(ctx context.Context) (int64, error) {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return -1, xerr.NewErrCode(xerr.UnauthorizedError)
	}

	return hash.DecodedId(userId), nil
}

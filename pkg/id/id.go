package id

import (
	"context"

	"github.com/speps/go-hashids"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"
)

var hash *hashids.HashID

func init() {
	hash = NewHashIds()
}

// type HashID struct {
// 	hash *hashids.HashID
// }

// func NewHashID() *HashID {
// 	return &HashID{
// 		hash: NewHashIds(),
// 	}
// }

func NewHashIds() *hashids.HashID {
	var hd = hashids.NewData()
	hd.Salt = "123poiasdmnmb"
	hd.MinLength = 8
	hd.Alphabet = "qwertyuiopasdfghjklzxcvbnm"
	h, _ := hashids.NewWithData(hd)
	return h
}

// func (h *HashID) EncodedId(id int64) (string, error) {
// 	encodedId, err := h.hash.EncodeInt64([]int64{id})
// 	if err != nil {
// 		return "", err
// 	}

// 	return encodedId, nil
// }

// func (h *HashID) DecodedId(id string) int64 {
// 	return h.hash.DecodeInt64(id)[0]
// }

func GetDecodedUserId(ctx context.Context) (int64, error) {
	userId, ok := ctx.Value(vars.UserIdKey).(string)
	if !ok {
		return 0, xerr.NewErrCode(xerr.UnauthorizedError)
	}
	userIdDecoded, err := DecodeId(userId)
	if err != nil {
		return 0, err
	}

	return userIdDecoded, nil
}

func EncodeId(id int64) string {
	ide, _ := hash.EncodeInt64([]int64{id})
	return ide
}

func DecodeId(id string) (int64, error) {
	res := hash.DecodeInt64(id)
	return res[0], nil
}

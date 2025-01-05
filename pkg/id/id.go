package id

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"ymir.com/pkg/vars"
	"ymir.com/pkg/xerr"
)

// var Hash *HashID

// func init() {
// 	Hash = NewHashID()
// }

// type HashID struct {
// 	hash *hashids.HashID
// }

// func NewHashID() *HashID {
// 	return &HashID{
// 		hash: NewHashIds(),
// 	}
// }

// func NewHashIds() *hashids.HashID {
// 	var hd = hashids.NewData()
// 	hd.Salt = "123poiasdmnmb"
// 	hd.MinLength = 8
// 	hd.Alphabet = "qwertyuiopasdfghjklzxcvbnm"
// 	h, _ := hashids.NewWithData(hd)
// 	return h
// }

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
	userId, err := ctx.Value(vars.UserIdKey).(json.Number).Int64()
	if err != nil {
		return -1, errors.Wrap(xerr.NewErrCode(xerr.UnauthorizedError), "empty id")
	}

	return userId, nil
}

func EncodeId(id int64) string {
	return strconv.FormatInt(id, 10)
}

func DecodeId(id string) (int64, error) {
	res, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, xerr.NewErrCode(xerr.ReuqestParamError)
	}
	return res, nil
}

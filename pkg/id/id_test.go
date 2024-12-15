package id

import (
	"testing"

	"github.com/google/uuid"
)

func TestId(t *testing.T) {
	var hd = NewHashID()
	t.Log(hd.EncodedId(1))
}

func TestUUID(t *testing.T) {
	id, _ := uuid.NewRandom()
	t.Log(id)
}

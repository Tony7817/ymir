package id

import (
	"testing"

	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	id, _ := uuid.NewRandom()
	t.Log(id)
}

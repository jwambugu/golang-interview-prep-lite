package testutil

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"github.com/jaswdr/faker"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/model"
	"math/rand"
)

func Seed() int64 {
	var b [8]byte
	_, err := cryptorand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

func NewUser() *model.User {
	f := faker.NewWithSeed(rand.NewSource(Seed()))

	return &model.User{
		Username: f.Person().FirstName() + f.Numerify("####"),
		Password: "password",
	}
}

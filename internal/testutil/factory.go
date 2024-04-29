package testutil

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"github.com/jaswdr/faker"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/user"
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

func NewUser() *user.User {
	f := faker.NewWithSeed(rand.NewSource(Seed()))

	return &user.User{
		Username: f.Person().FirstName() + f.Numerify("####"),
		Password: "password",
	}
}

package uuidhelper

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"

	guuid "github.com/google/uuid"
)

var letters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var numerics = []rune("0123456789")

func NewUuidV4() guuid.UUID {
	return guuid.New()
}

func NewUuidV4String() string {
	// return guuid.New().String()
	return guuid.NewString()
}

func NewUuidV7String() string {
	uuidV7Value, err := guuid.NewV7()
	if err != nil {
		return ""
	}

	return uuidV7Value.String()
}

func ParseUuidV4(uuidV4String string) (guuid.UUID, error) {
	return guuid.Parse(uuidV4String)
}

func CreateTraceId() string {
	return guuid.NewString()
}

func CreateTraceId12() string {
	dest, _ := hex.DecodeString(fmt.Sprintf("%d", nowAsUnixSecond()))
	var id strings.Builder
	encode := base64.StdEncoding.EncodeToString(dest)
	// rand.Seed(time.Now().UnixNano())
	id.WriteString(encode)
	id.WriteString(RandString(4))
	return strings.Replace(id.String(), "=", RandString(1), 1)
}

func CreateRefreshToken() string {
	dest, _ := hex.DecodeString(fmt.Sprintf("%d", nowAsUnixSecond()))
	var id strings.Builder
	encode := base64.StdEncoding.EncodeToString(dest)
	// rand.Seed(time.Now().UnixNano())
	id.WriteString(encode)
	id.WriteString(RandString(4))
	return strings.Replace(id.String(), "=", RandString(1), 1)
}

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandIntn(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numerics[rand.Intn(len(numerics))]
	}
	return string(b)
}

func nowAsUnixSecond() int64 {
	return time.Now().UnixNano() / 1e9
}

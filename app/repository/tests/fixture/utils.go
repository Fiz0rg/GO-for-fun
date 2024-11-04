package fixture

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
)

func GetTimeNow() int64 {
	res := time.Now()
	return res.Unix()
}

func GenUUID() string {
	return uuid.New().String()
}

func genName() string {
	return gofakeit.Name()
}

func genString() string {
	return gofakeit.Letter()
}

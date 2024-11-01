package fixture

import (
	"time"

	"github.com/google/uuid"
)

func GetTimeNow() int64 {
	res := time.Now()
	return res.Unix()
}

func genUUID() string {
	return uuid.New().String()
}

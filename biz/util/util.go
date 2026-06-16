package util

import "github.com/oklog/ulid/v2"

func GenerateTxID() string {
	return ulid.Make().String()
}

package idgen

import (
	"crypto/rand"
	"encoding/hex"
)

type IdGenerater struct {
	db []string
}

func (gen *IdGenerater) Clear() {
	gen.db = make([]string, 0)
}

func (gen *IdGenerater) GenerateID() string {
	for {
		math := true
		hex, _ := randomHex(2)
		for i := range gen.db {
			if gen.db[i] == hex {
				math = false
				break
			}
		}

		if math {
			gen.db = append(gen.db, hex)
			return hex
		}
	}
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"math/big"
)

func Hash(s any) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(s)
	sum := sha256.Sum256(b.Bytes())
	return fmt.Sprintf("%x", sum)
}

func GetRandomNumberInt64(until int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(until))
	if err != nil {
		return -1
	}
	return nBig.Int64()
}

func GetRandomNumberUint64(until uint64) uint64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(until)))
	if err != nil {
		return 0
	}
	return nBig.Uint64()
}

func GetRandomNumberBool() bool {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1))
	if err != nil {
		return false
	}
	return nBig.Int64() == 1
}

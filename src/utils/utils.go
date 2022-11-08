package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

func Hash(s any) string {
	var b bytes.Buffer
	gob.NewEncoder(&b).Encode(s)
	sum := sha256.Sum256(b.Bytes())
	return fmt.Sprintf("%x", sum)
}

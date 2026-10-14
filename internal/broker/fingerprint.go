package broker

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

func Fingerprint(queue string,priority int, payload []byte) string {
	h := sha256.New()
	h.Write([]byte(queue))
	h.Write([]byte{0})
	h.Write([]byte(strconv.Itoa(priority)))
	h.Write([]byte{0})
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
	//fingerprint : sha256(queue + 0 + priority + 0 + payload)
}
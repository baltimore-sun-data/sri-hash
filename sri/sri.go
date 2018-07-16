package sri

import (
	"crypto/sha512"
	"encoding/base64"
	"io"
)

func FromReader(r io.Reader) (string, error) {
	h := sha512.New384()
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	b := h.Sum(nil)
	return "sha384-" + base64.StdEncoding.EncodeToString(b), nil
}

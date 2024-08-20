package utils

import "github.com/hidevopsio/hiboot/pkg/utils/crypto/rsa"

func Decrypt(param string, pk string) (retVal string) {
	if len(param) != 0 {
		b, err := rsa.DecryptBase64([]byte(param), []byte(pk))
		if err == nil {
			retVal = string(b)
		}
	}
	return
}

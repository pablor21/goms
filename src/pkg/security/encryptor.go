package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"github.com/pablor21/goms/app/config"
)

func Encrypt(data string) (res string, err error) {
	aes, err := aes.NewCipher([]byte(config.GetConfig().Security.Encryption.Key))
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	// to base64 encode the ciphertext
	return base64.URLEncoding.EncodeToString(ciphertext), nil

	// return string(ciphertext), nil
}

func Decrypt(data string) (res string, err error) {

	// decode the base64 encoded ciphertext
	ciphertext, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return
	}

	aes, err := aes.NewCipher([]byte(config.GetConfig().Security.Encryption.Key))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext), nil
}

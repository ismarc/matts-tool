package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// unpackG unpacks the provided byte array into its constituent parts.  This
// unpacks the legacy format from the Ruby slosilo package for compatibility.
// FORMAT: {VERSION_MAGIC}{tag}{nonce}{ciphertext}
// SIZES: {1}{16}{12}{rest...}
func unpackG(data []byte) (nonce []byte, ciphertext []byte, err error) {
	versionIndex := 1
	version := data[:versionIndex]

	// Verify that data indicates it's in a matching format
	if string(version) != "G" {
		err = fmt.Errorf("version didn't match %s", string(version))
		return
	}

	tagIndex := versionIndex + 16
	nonceIndex := tagIndex + 12
	tag := data[versionIndex:tagIndex]
	nonce = data[tagIndex:nonceIndex]
	ciphertext = append(data[nonceIndex:], tag...)

	return
}

// packG repacks the data to conform to the Ruby slosilo package legacy package
// (G format).
// INCOMING FORMAT: {nonce}{cipher}{tag}
// SIZES: {12}{...}{16}
// OUTPUT_FORMAT: {VERSION_MAGIC}{tag}{nonce}{ciphertext}
func packG(data []byte) (ciphertext []byte) {
	nonceSize := 12
	tagSize := 16
	nonce := data[:nonceSize]
	tag := data[len(data)-tagSize:]
	cipher := data[nonceSize : len(data)-tagSize]
	ciphertext = append(append(append([]byte("G"), tag...), nonce...), cipher...)
	return
}

// unpack separates the provided byte array into its constituent parts.  This
// unpacks the golang related format.
// FORMAT: {VERSION_MAGIC}{nonce}{ciphertext}{tag}
// SIZES: {1}{12}{rest...}
func unpack(data []byte) (nonce []byte, ciphertext []byte, err error) {
	versionIndex := 1
	version := data[:versionIndex]

	nonceIndex := versionIndex + 12

	if string(version) != "H" {
		err = fmt.Errorf("version didn't match %s", string(version))
		return
	}

	nonce = data[versionIndex:nonceIndex]
	ciphertext = data[nonceIndex:]

	return
}

// pack repacks the data to conform to the golang related format
// INCOMING FORMAT: {nonce}{cipher}{tag}
// OUTPUT_FORMAT: {VERSION_MAGIC}{nonce}{cipher}{tag}
func pack(data []byte) (ciphertext []byte) {
	ciphertext = append([]byte("H"), data...)
	return
}

// RandomKey generates a secure random key
func RandomKey(length ...int) (key string, err error) {
	// AES GCM 256 key length is 32
	keyLength := 32
	if len(length) == 1 {
		if length[0] > 0 {
			keyLength = length[0]
		} else {
			err = fmt.Errorf("invalid key length: %+v", length[0])
			return
		}
	} else if len(length) > 1 {
		err = fmt.Errorf("too many parameters: %+v", length)
		return
	}

	bytes := make([]byte, keyLength)

	_, err = rand.Read(bytes)
	if err != nil {
		return
	}

	key = base64.StdEncoding.EncodeToString(bytes)
	return
}

// AES256GCMEncrypt encrypts the supplied plaintext using the given key and
// additional authentication data
func AES256GCMEncrypt(stringKey string, rawPlaintext string, additionalData string) (ciphertext string, err error) {

	key, err := base64.StdEncoding.DecodeString(stringKey)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	plaintext := []byte(rawPlaintext)

	rawCiphertext := gcm.Seal(nonce, nonce, plaintext, []byte(additionalData))
	packed := pack(rawCiphertext)
	ciphertext = hex.EncodeToString(packed)

	return
}

// AES256GCMDecrypt decrypt supplied ciphertext with the given key authenticated
// against the supplied additional data
func AES256GCMDecrypt(stringKey string, rawCipherText string, additionalData string) (plaintext []byte, err error) {
	key, err := base64.StdEncoding.DecodeString(stringKey)
	if err != nil {
		return
	}

	rawCipher, err := hex.DecodeString(rawCipherText)
	if err != nil {
		return
	}

	var nonce, ciphertext []byte

	switch string(rawCipher[:1]) {
	case "G":
		nonce, ciphertext, err = unpackG(rawCipher)
	case "H":
		nonce, ciphertext, err = unpack(rawCipher)
	default:
		err = fmt.Errorf("unknown format indicator '%v'", rawCipher[:1])
	}

	if err != nil {
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	plaintext, err = gcm.Open(nil, nonce, ciphertext, []byte(additionalData))
	return
}

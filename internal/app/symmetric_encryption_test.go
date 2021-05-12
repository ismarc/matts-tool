package app

import (
	"encoding/hex"

	"github.com/stretchr/testify/assert"
)

func (s *Suite) TestAES256GCM() {
	randKey, err := RandomKey()
	assert.Nil(s.T(), err)

	randKey24, err := RandomKey(24)
	assert.Nil(s.T(), err)

	plaintext := "You can't trust the weatherman, not in the summer."

	// Uses random nonce, cipher includes the nonce
	// To verify behavior, two runs should generate different output
	// given the same inputs, but both decrypt to the same value

	// Encrypt/Decrypt with no additional data
	cipherNoAdditional, err := AES256GCMEncrypt(randKey, plaintext, "")
	assert.Nil(s.T(), err)
	cipherNoAdditional2, err := AES256GCMEncrypt(randKey, plaintext, "")
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), cipherNoAdditional, cipherNoAdditional2)

	decryptedNoAdditional, err := AES256GCMDecrypt(randKey, cipherNoAdditional, "")
	assert.Nil(s.T(), err)
	decryptedNoAdditional2, err := AES256GCMDecrypt(randKey, cipherNoAdditional2, "")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), decryptedNoAdditional, decryptedNoAdditional2)

	// Encrypt/Decrypt with additional data
	cipher, err := AES256GCMEncrypt(randKey, plaintext, "hey data")
	assert.Nil(s.T(), err)
	cipher2, err := AES256GCMEncrypt(randKey, plaintext, "hey data")
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), cipher, cipher2)

	decrypted, err := AES256GCMDecrypt(randKey, cipher, "hey data")
	assert.Nil(s.T(), err)
	decrypted2, err := AES256GCMDecrypt(randKey, cipher2, "hey data")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), decrypted, decrypted2)

	// Validate key length of 24
	cipher24, err := AES256GCMEncrypt(randKey24, plaintext, "")
	assert.Nil(s.T(), err)
	decrypted24, err := AES256GCMDecrypt(randKey24, cipher24, "")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(decrypted24), plaintext)
}

func (s *Suite) TestAES256GCMG() {
	// Test decryption from slosilo compatible format encrypted values
	key := "6QrDHLBWYXieY5FM5DlRWRXX/wA8hefCuwMciHQ5ms0="
	encrypted := "4760d5309fdd5acd030976dc826d8ab9fbe31360f60cc636aa58764bec63cad3f9bdffa5255649bbb5581f25f4677b47cfb9b8d037c44843a39d086a2540a396f26d31340bf9d2bbe1a2062fe6a76c7ab3"
	additional := "myConjurAccount:user:admin"
	expected := "13cntq93tw9kmzrw757w22j88phb2yadbxea64gzhng4j27mnzev"

	// Check pack/unpack
	rawBytes, err := hex.DecodeString(encrypted)
	assert.Nil(s.T(), err)
	nonce, ciphertext, err := unpackG(rawBytes)
	assert.Nil(s.T(), err)

	packed := packG(append(nonce, ciphertext...))
	packedString := hex.EncodeToString(packed)
	assert.Equal(s.T(), packedString, encrypted)

	// Check G format decryption
	plaintext, err := AES256GCMDecrypt(key, encrypted, additional)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), string(plaintext), expected)
}

func (s *Suite) TestErrorPaths() {
	key := "6QrDHLBWYXieY5FM5DlRWRXX/wA8hefCuwMciHQ5ms0="
	badVersion := "8880d5309fdd5acd0309"

	// Validate value unpacking error cases
	_, _, err := unpackG([]byte(badVersion))
	assert.NotNil(s.T(), err)
	err = nil
	_, _, err = unpack([]byte(badVersion))
	assert.NotNil(s.T(), err)

	// Validate RandomKey error cases
	err = nil
	_, err = RandomKey(-1)
	assert.NotNil(s.T(), err)
	err = nil
	_, err = RandomKey(5, 10)
	assert.NotNil(s.T(), err)

	// Validate decrypt error cases
	err = nil
	_, err = AES256GCMDecrypt("BADKEY", "plaintext", "")
	assert.NotNil(s.T(), err)
	err = nil
	_, err = AES256GCMDecrypt(key, "BADHEXDATA", "")
	assert.NotNil(s.T(), err)
	err = nil
	_, err = AES256GCMDecrypt(key, badVersion, "")
	assert.NotNil(s.T(), err)

	// Validate encrypt error cases
	err = nil
	_, err = AES256GCMEncrypt("BADKEY", "plaintext", "")
	assert.NotNil(s.T(), err)
	err = nil
	randKey6, err := RandomKey(6)
	assert.Nil(s.T(), err)
	_, err = AES256GCMEncrypt(randKey6, "plaintext", "")
	assert.NotNil(s.T(), err)
}

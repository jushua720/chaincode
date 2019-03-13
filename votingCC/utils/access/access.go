package access

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"math/big"

	c "../constants"
	enc "github.com/btcsuite/btcutil/base58"
)

type Keys struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

//@note ecdsa alg key is shorter than rsa

func GenerateKeys() (*Keys, error) {

	var keys Keys

	reader := rand.Reader

	// @note bitSize : 512, 1024, 2048
	bitSize := 512
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return &keys, errors.New(err.Error())
	}

	keys.PrivateKey, keys.PublicKey, err = encodeKeys(key, key.PublicKey)

	return &keys, nil
}

func encodeKeys(privKey *rsa.PrivateKey, pubKey rsa.PublicKey) (string, string, error) {

	privateKey := b58encode(x509.MarshalPKCS1PrivateKey(privKey))

	publicKeyByte, err := asn1.Marshal(pubKey)
	if err != nil {
		return "", "", errors.New(err.Error())
	}

	publicKey := enc.Encode(publicKeyByte)

	return privateKey, publicKey, nil
}

func b58encode(b []byte) string {

	// @note : see https://en.bitcoin.it/wiki/Base58Check_encoding

	x := new(big.Int).SetBytes(b)

	r := new(big.Int)
	m := big.NewInt(58)
	zero := big.NewInt(0)
	str := ""

	for x.Cmp(zero) > 0 {
		/* x, r = (x / 58, x % 58) */
		x.QuoRem(x, m, r)
		/* Prepend ASCII character */
		str = string(c.Base58Table[r.Int64()]) + str
	}

	return str
}

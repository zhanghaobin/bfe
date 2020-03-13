// Copyright (c) 2019 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// AES Key Wrap with default initial value
// The default initial value (IV) is defined to be the hexadecimal
// constant: A[0] = IV = A6A6A6A6A6A6A6A6
// see: https://tools.ietf.org/html/rfc3394
package jwa

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"github.com/baidu/bfe/bfe_modules/mod_auth_jwt/jwk"
)

const IV uint64 = 12008468691120727718 // hex: A6A6A6A6A6A6A6A6

type AKW struct {
	block cipher.Block
}

func (akw *AKW) unwrap(eCek []byte) (iv uint64, r []uint64, err error) {
	defer CatchCryptoPanic(&err) // prevent from panic may be caused by AES decryption
	length := len(eCek)
	if length%8 != 0 || length < 24 {
		return 0, nil, fmt.Errorf("invalid bit length of encrypted CEK: %d", len(eCek)*8)
	}

	n := length/8 - 1
	aiv := binary.BigEndian.Uint64(eCek[:8])
	r = make([]uint64, n)

	for i := 0; i < n; i++ {
		r[i] = binary.BigEndian.Uint64(eCek[(i+1)*8:])
	}

	axt, dst := make([]byte, 16), make([]byte, akw.block.BlockSize())
	for j := 5; j >= 0; j-- {
		for i := n; i > 0; i-- {
			binary.BigEndian.PutUint64(axt, aiv^uint64(n*j+i))
			binary.BigEndian.PutUint64(axt[8:], r[i-1])
			akw.block.Decrypt(dst, axt) // this may cause panic

			// update iv and register
			aiv = binary.BigEndian.Uint64(dst[:8])
			r[i-1] = binary.BigEndian.Uint64(dst[len(dst)-8:])
		}
	}

	return aiv, r, nil
}

func (akw *AKW) Decrypt(eCek []byte) (cek []byte, err error) {
	iv, r, err := akw.unwrap(eCek)
	if err != nil {
		return nil, err
	}

	if iv != IV {
		return nil, fmt.Errorf("decrypted iv not default iv")
	}

	kLength := len(r)
	cek = make([]byte, kLength*8)

	// covert []uint64 to []byte
	for i := 0; i < kLength; i++ {
		binary.BigEndian.PutUint64(cek[i*8:], r[i])
	}

	return cek, nil
}

func NewAKW(kBit int, mJWK *jwk.JWK) (akw JWEAlg, err error) {
	if mJWK.Kty != jwk.OCT {
		return nil, fmt.Errorf("unsupported algorithm: A%dKW", kBit)
	}

	kLen := len(mJWK.Symmetric.K.Decoded)
	if kBit/8 != kLen {
		return nil, fmt.Errorf("invalid key length for algorithm A%dKW: %d", kBit, kLen*8)
	}

	block, err := aes.NewCipher(mJWK.Symmetric.K.Decoded)
	if err != nil {
		return nil, err
	}

	return &AKW{block}, nil
}

func NewA128KW(mJWK *jwk.JWK, _ map[string]interface{}) (akw JWEAlg, err error) {
	return NewAKW(128, mJWK)
}

func NewA192KW(mJWK *jwk.JWK, _ map[string]interface{}) (akw JWEAlg, err error) {
	return NewAKW(192, mJWK)
}

func NewA256KW(mJWK *jwk.JWK, _ map[string]interface{}) (akw JWEAlg, err error) {
	return NewAKW(256, mJWK)
}

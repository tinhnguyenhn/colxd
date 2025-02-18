// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2016 The Dash developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcec_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/fastsha256"
	"github.com/tinhnguyenhn/colxd/btcec"
)

type signatureTest struct {
	name    string
	sig     []byte
	der     bool
	isValid bool
}

// decodeHex decodes the passed hex string and returns the resulting bytes.  It
// panics if an error occurs.  This is only used in the tests as a helper since
// the only way it can fail is if there is an error in the test source code.
func decodeHex(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic("invalid hex string in test source: err " + err.Error() +
			", hex: " + hexStr)
	}

	return b
}

var signatureTests = []signatureTest{
	// signatures from bitcoin blockchain tx
	// 0437cd7f8525ceed2324359c2d0ba26006d92d85
	{
		name: "valid signature.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: true,
	},
	{
		name:    "empty.",
		sig:     []byte{},
		isValid: false,
	},
	{
		name: "bad magic.",
		sig: []byte{0x31, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "bad 1st int marker magic.",
		sig: []byte{0x30, 0x44, 0x03, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "bad 2nd int marker.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x03, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "short len",
		sig: []byte{0x30, 0x43, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "long len",
		sig: []byte{0x30, 0x45, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "long X",
		sig: []byte{0x30, 0x44, 0x02, 0x42, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "long Y",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x21, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "short Y",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x19, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "trailing crap.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09, 0x01,
		},
		der: true,

		// This test is now passing (used to be failing) because there
		// are signatures in the blockchain that have trailing zero
		// bytes before the hashtype. So ParseSignature was fixed to
		// permit buffers with trailing nonsense after the actual
		// signature.
		isValid: true,
	},
	{
		name: "X == N ",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48,
			0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "X == N ",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48,
			0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41,
			0x42, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     false,
		isValid: false,
	},
	{
		name: "Y == N",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B,
			0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x41,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "Y > N",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			0xFE, 0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B,
			0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x42,
		},
		der:     false,
		isValid: false,
	},
	{
		name: "0 len X.",
		sig: []byte{0x30, 0x24, 0x02, 0x00, 0x02, 0x20, 0x18, 0x15,
			0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60, 0xa4,
			0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c, 0xc5, 0x6c,
			0xbb, 0xac, 0x46, 0x22, 0x08, 0x22, 0x21, 0xa8, 0x76,
			0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "0 len Y.",
		sig: []byte{0x30, 0x24, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x00,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "extra R padding.",
		sig: []byte{0x30, 0x45, 0x02, 0x21, 0x00, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	{
		name: "extra S padding.",
		sig: []byte{0x30, 0x45, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x21, 0x00, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
	// Standard checks (in BER format, without checking for 'canonical' DER
	// signatures) don't test for negative numbers here because there isn't
	// a way that is the same between openssl and go that will mark a number
	// as negative. The Go ASN.1 parser marks numbers as negative when
	// openssl does not (it doesn't handle negative numbers that I can tell
	// at all. When not parsing DER signatures, which is done by by bitcoind
	// when accepting transactions into its mempool, we otherwise only check
	// for the coordinates being zero.
	{
		name: "X == 0",
		sig: []byte{0x30, 0x25, 0x02, 0x01, 0x00, 0x02, 0x20, 0x18,
			0x15, 0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60,
			0xa4, 0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c, 0xc5,
			0x6c, 0xbb, 0xac, 0x46, 0x22, 0x08, 0x22, 0x21, 0xa8,
			0x76, 0x8d, 0x1d, 0x09,
		},
		der:     false,
		isValid: false,
	},
	{
		name: "Y == 0.",
		sig: []byte{0x30, 0x25, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x01, 0x00,
		},
		der:     false,
		isValid: false,
	},
}

func TestSignatures(t *testing.T) {
	for _, test := range signatureTests {
		var err error
		if test.der {
			_, err = btcec.ParseDERSignature(test.sig, btcec.S256())
		} else {
			_, err = btcec.ParseSignature(test.sig, btcec.S256())
		}
		if err != nil {
			if test.isValid {
				t.Errorf("%s signature failed when shouldn't %v",
					test.name, err)
			} /* else {
				t.Errorf("%s got error %v", test.name, err)
			} */
			continue
		}
		if !test.isValid {
			t.Errorf("%s counted as valid when it should fail",
				test.name)
		}
	}
}

// TestSignatureSerialize ensures that serializing signatures works as expected.
func TestSignatureSerialize(t *testing.T) {
	tests := []struct {
		name     string
		ecsig    *btcec.Signature
		expected []byte
	}{
		// signature from bitcoin blockchain tx
		// 0437cd7f8525ceed2324359c2d0ba26006d92d85
		{
			"valid 1 - r and s most significant bits are zero",
			&btcec.Signature{
				R: fromHex("4e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd41"),
				S: fromHex("181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d09"),
			},
			[]byte{
				0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
				0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3,
				0xa1, 0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32,
				0xe9, 0xd6, 0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab,
				0x5f, 0xb8, 0xcd, 0x41, 0x02, 0x20, 0x18, 0x15,
				0x22, 0xec, 0x8e, 0xca, 0x07, 0xde, 0x48, 0x60,
				0xa4, 0xac, 0xdd, 0x12, 0x90, 0x9d, 0x83, 0x1c,
				0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22, 0x08, 0x22,
				0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
			},
		},
		// signature from bitcoin blockchain tx
		// cb00f8a0573b18faa8c4f467b049f5d202bf1101d9ef2633bc611be70376a4b4
		{
			"valid 2 - r most significant bit is one",
			&btcec.Signature{
				R: fromHex("0082235e21a2300022738dabb8e1bbd9d19cfb1e7ab8c30a23b0afbb8d178abcf3"),
				S: fromHex("24bf68e256c534ddfaf966bf908deb944305596f7bdcc38d69acad7f9c868724"),
			},
			[]byte{
				0x30, 0x45, 0x02, 0x21, 0x00, 0x82, 0x23, 0x5e,
				0x21, 0xa2, 0x30, 0x00, 0x22, 0x73, 0x8d, 0xab,
				0xb8, 0xe1, 0xbb, 0xd9, 0xd1, 0x9c, 0xfb, 0x1e,
				0x7a, 0xb8, 0xc3, 0x0a, 0x23, 0xb0, 0xaf, 0xbb,
				0x8d, 0x17, 0x8a, 0xbc, 0xf3, 0x02, 0x20, 0x24,
				0xbf, 0x68, 0xe2, 0x56, 0xc5, 0x34, 0xdd, 0xfa,
				0xf9, 0x66, 0xbf, 0x90, 0x8d, 0xeb, 0x94, 0x43,
				0x05, 0x59, 0x6f, 0x7b, 0xdc, 0xc3, 0x8d, 0x69,
				0xac, 0xad, 0x7f, 0x9c, 0x86, 0x87, 0x24,
			},
		},
		// signature from bitcoin blockchain tx
		// fda204502a3345e08afd6af27377c052e77f1fefeaeb31bdd45f1e1237ca5470
		{
			"valid 3 - s most significant bit is one",
			&btcec.Signature{
				R: fromHex("1cadddc2838598fee7dc35a12b340c6bde8b389f7bfd19a1252a17c4b5ed2d71"),
				S: new(big.Int).Add(fromHex("00c1a251bbecb14b058a8bd77f65de87e51c47e95904f4c0e9d52eddc21c1415ac"), btcec.S256().N),
			},
			[]byte{
				0x30, 0x45, 0x02, 0x20, 0x1c, 0xad, 0xdd, 0xc2,
				0x83, 0x85, 0x98, 0xfe, 0xe7, 0xdc, 0x35, 0xa1,
				0x2b, 0x34, 0x0c, 0x6b, 0xde, 0x8b, 0x38, 0x9f,
				0x7b, 0xfd, 0x19, 0xa1, 0x25, 0x2a, 0x17, 0xc4,
				0xb5, 0xed, 0x2d, 0x71, 0x02, 0x21, 0x00, 0xc1,
				0xa2, 0x51, 0xbb, 0xec, 0xb1, 0x4b, 0x05, 0x8a,
				0x8b, 0xd7, 0x7f, 0x65, 0xde, 0x87, 0xe5, 0x1c,
				0x47, 0xe9, 0x59, 0x04, 0xf4, 0xc0, 0xe9, 0xd5,
				0x2e, 0xdd, 0xc2, 0x1c, 0x14, 0x15, 0xac,
			},
		},
		{
			"zero signature",
			&btcec.Signature{
				R: big.NewInt(0),
				S: big.NewInt(0),
			},
			[]byte{0x30, 0x06, 0x02, 0x01, 0x00, 0x02, 0x01, 0x00},
		},
	}

	for i, test := range tests {
		result := test.ecsig.Serialize()
		if !bytes.Equal(result, test.expected) {
			t.Errorf("Serialize #%d (%s) unexpected result:\n"+
				"got:  %x\nwant: %x", i, test.name, result,
				test.expected)
		}
	}
}

func testSignCompact(t *testing.T, tag string, curve *btcec.KoblitzCurve,
	data []byte, isCompressed bool) {
	tmp, _ := btcec.NewPrivateKey(curve)
	priv := (*btcec.PrivateKey)(tmp)

	hashed := []byte("testing")
	sig, err := btcec.SignCompact(curve, priv, hashed, isCompressed)
	if err != nil {
		t.Errorf("%s: error signing: %s", tag, err)
		return
	}

	pk, wasCompressed, err := btcec.RecoverCompact(curve, sig, hashed)
	if err != nil {
		t.Errorf("%s: error recovering: %s", tag, err)
		return
	}
	if pk.X.Cmp(priv.X) != 0 || pk.Y.Cmp(priv.Y) != 0 {
		t.Errorf("%s: recovered pubkey doesn't match original "+
			"(%v,%v) vs (%v,%v) ", tag, pk.X, pk.Y, priv.X, priv.Y)
		return
	}
	if wasCompressed != isCompressed {
		t.Errorf("%s: recovered pubkey doesn't match compressed state "+
			"(%v vs %v)", tag, isCompressed, wasCompressed)
		return
	}

	// If we change the compressed bit we should get the same key back,
	// but the compressed flag should be reversed.
	if isCompressed {
		sig[0] -= 4
	} else {
		sig[0] += 4
	}

	pk, wasCompressed, err = btcec.RecoverCompact(curve, sig, hashed)
	if err != nil {
		t.Errorf("%s: error recovering (2): %s", tag, err)
		return
	}
	if pk.X.Cmp(priv.X) != 0 || pk.Y.Cmp(priv.Y) != 0 {
		t.Errorf("%s: recovered pubkey (2) doesn't match original "+
			"(%v,%v) vs (%v,%v) ", tag, pk.X, pk.Y, priv.X, priv.Y)
		return
	}
	if wasCompressed == isCompressed {
		t.Errorf("%s: recovered pubkey doesn't match reversed "+
			"compressed state (%v vs %v)", tag, isCompressed,
			wasCompressed)
		return
	}
}

func TestSignCompact(t *testing.T) {
	for i := 0; i < 256; i++ {
		name := fmt.Sprintf("test %d", i)
		data := make([]byte, 32)
		_, err := rand.Read(data)
		if err != nil {
			t.Errorf("failed to read random data for %s", name)
			continue
		}
		compressed := i%2 != 0
		testSignCompact(t, name, btcec.S256(), data, compressed)
	}
}

func TestRFC6979(t *testing.T) {
	// Test vectors matching Trezor and CoreBitcoin implementations.
	// - https://github.com/trezor/trezor-crypto/blob/9fea8f8ab377dc514e40c6fd1f7c89a74c1d8dc6/tests.c#L432-L453
	// - https://github.com/oleganza/CoreBitcoin/blob/e93dd71207861b5bf044415db5fa72405e7d8fbc/CoreBitcoin/BTCKey%2BTests.m#L23-L49
	tests := []struct {
		key       string
		msg       string
		nonce     string
		signature string
	}{
		{
			"cca9fbcc1b41e5a95d369eaa6ddcff73b61a4efaa279cfc6567e8daa39cbaf50",
			"sample",
			"2df40ca70e639d89528a6b670d9d48d9165fdc0febc0974056bdce192b8e16a3",
			"3045022100af340daf02cc15c8d5d08d7735dfe6b98a474ed373bdb5fbecf7571be52b384202205009fb27f37034a9b24b707b7c6b79ca23ddef9e25f7282e8a797efe53a8f124",
		},
		{
			// This signature hits the case when S is higher than halforder.
			// If S is not canonicalized (lowered by halforder), this test will fail.
			"0000000000000000000000000000000000000000000000000000000000000001",
			"Satoshi Nakamoto",
			"8f8a276c19f4149656b280621e358cce24f5f52542772691ee69063b74f15d15",
			"3045022100934b1ea10a4b3c1757e2b0c017d0b6143ce3c9a7e6a4a49860d7a6ab210ee3d802202442ce9d2b916064108014783e923ec36b49743e2ffa1c4496f01a512aafd9e5",
		},
		{
			"fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364140",
			"Satoshi Nakamoto",
			"33a19b60e25fb6f4435af53a3d42d493644827367e6453928554f43e49aa6f90",
			"3045022100fd567d121db66e382991534ada77a6bd3106f0a1098c231e47993447cd6af2d002206b39cd0eb1bc8603e159ef5c20a5c8ad685a45b06ce9bebed3f153d10d93bed5",
		},
		{
			"f8b8af8ce3c7cca5e300d33939540c10d45ce001b8f252bfbc57ba0342904181",
			"Alan Turing",
			"525a82b70e67874398067543fd84c83d30c175fdc45fdeee082fe13b1d7cfdf1",
			"304402207063ae83e7f62bbb171798131b4a0564b956930092b33b07b395615d9ec7e15c022058dfcc1e00a35e1572f366ffe34ba0fc47db1e7189759b9fb233c5b05ab388ea",
		},
		{
			"0000000000000000000000000000000000000000000000000000000000000001",
			"All those moments will be lost in time, like tears in rain. Time to die...",
			"38aa22d72376b4dbc472e06c3ba403ee0a394da63fc58d88686c611aba98d6b3",
			"30450221008600dbd41e348fe5c9465ab92d23e3db8b98b873beecd930736488696438cb6b0220547fe64427496db33bf66019dacbf0039c04199abb0122918601db38a72cfc21",
		},
		{
			"e91671c46231f833a6406ccbea0e3e392c76c167bac1cb013f6f1013980455c2",
			"There is a computer disease that anybody who works with computers knows about. It's a very serious disease and it interferes completely with the work. The trouble with computers is that you 'play' with them!",
			"1f4b84c23a86a221d233f2521be018d9318639d5b8bbd6374a8a59232d16ad3d",
			"3045022100b552edd27580141f3b2a5463048cb7cd3e047b97c9f98076c32dbdf85a68718b0220279fa72dd19bfae05577e06c7c0c1900c371fcd5893f7e1d56a37d30174671f6",
		},
	}

	for i, test := range tests {
		privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), decodeHex(test.key))
		hash := fastsha256.Sum256([]byte(test.msg))

		// Ensure deterministically generated nonce is the expected value.
		gotNonce := btcec.TstNonceRFC6979(privKey.D, hash[:]).Bytes()
		wantNonce := decodeHex(test.nonce)
		if !bytes.Equal(gotNonce, wantNonce) {
			t.Errorf("NonceRFC6979 #%d (%s): Nonce is incorrect: "+
				"%x (expected %x)", i, test.msg, gotNonce,
				wantNonce)
			continue
		}

		// Ensure deterministically generated signature is the expected value.
		gotSig, err := privKey.Sign(hash[:])
		if err != nil {
			t.Errorf("Sign #%d (%s): unexpected error: %v", i,
				test.msg, err)
			continue
		}
		gotSigBytes := gotSig.Serialize()
		wantSigBytes := decodeHex(test.signature)
		if !bytes.Equal(gotSigBytes, wantSigBytes) {
			t.Errorf("Sign #%d (%s): mismatched signature: %x "+
				"(expected %x)", i, test.msg, gotSigBytes,
				wantSigBytes)
			continue
		}
	}
}

func TestSignatureIsEqual(t *testing.T) {
	sig1 := &btcec.Signature{
		R: fromHex("0082235e21a2300022738dabb8e1bbd9d19cfb1e7ab8c30a23b0afbb8d178abcf3"),
		S: fromHex("24bf68e256c534ddfaf966bf908deb944305596f7bdcc38d69acad7f9c868724"),
	}
	sig2 := &btcec.Signature{
		R: fromHex("4e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd41"),
		S: fromHex("181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d09"),
	}

	if !sig1.IsEqual(sig1) {
		t.Fatalf("value of IsEqual is incorrect, %v is "+
			"equal to %v", sig1, sig1)
	}

	if sig1.IsEqual(sig2) {
		t.Fatalf("value of IsEqual is incorrect, %v is not "+
			"equal to %v", sig1, sig2)
	}
}

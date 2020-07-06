// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet

// FLnet is an interface that defines FL-net messages.
type FLnet interface {
	MarshalBinary() ([]byte, error)
	MarshalTo(b []byte) error
	MarshalLen() int
	UnmarshalBinary(b []byte) error
}

// Token is a token frame of FA Link frame.
type Token struct {
	Header *FALinkHeader
}

// NewToken creates a new Token.
func NewToken() *Token {
	t := &Token{
		Header: NewFALinkHeader(
			[4]uint8{0x46, 0x41, 0x43, 0x4e}, // H_TYPE
			0x40,                             // TFL
			0x1,                              // SA
			0x55,                             // DA
			0,                                // V_SEQ
			0,                                // SEQ
			false, false, false,              // M_CTL
			0, 0, // ULS, M_SZ
			0,       // M_ADD
			0, 0, 0, // MFT, M_RLT, reserved
			TcdToken, 0, // TCD, VER
			0, 0, // C_AD1, C_SZ1
			0, 0, // C_AD2, C_SZ2
			0, 3, true, 0x80, 0, // MODE, P_TYPE, PRI
			1, 1, 0x40, // CBN, TBN, BSIZE
			0, 0x32, 0, // LKS, TW, RCT
			[]byte{},
		),
	}
	return t
}

// MarshalBinary returns the byte sequence generated from a Token.
func (t *Token) MarshalBinary() ([]byte, error) {
	b, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalLen returns the serial length of Token.
func (t *Token) MarshalLen() int {
	return 64
}

// UnmarshalBinary sets the values retrieved from byte sequence in a FL-net common header.
func (t *Token) UnmarshalBinary(b []byte) error {
	err := t.Header.UnmarshalBinary(b)
	if err != nil {
		return err
	}
	return nil
}

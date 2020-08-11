// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet

import (
	"io"
)

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
			0x1,                              // SNA
			0x55,                             // DNA
			0,                                // V_SEQ
			0,                                // SEQ
			false, false, false,              // M_CTL
			0, 0, // ULS, M_SZ
			0,       // M_ADD
			0, 0, 0, // MFT, M_RLT, reserved
			TCDToken, 0, // TCD, VER
			0, 0, // C_AD1, C_SZ1
			0, 0, // C_AD2, C_SZ2
			0, 3, true, 0x80, 0, // MODE, P_TYPE, PRI
			1, 1, 0x40, // CBN, TBN, BSIZE
			0, 0x32, 0, // LKS, TW, RCT
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
	return t.Header.MarshalLen()
}

// UnmarshalBinary sets the values retrieved from byte sequence in a FL-net common header.
func (t *Token) UnmarshalBinary(b []byte) error {
	err := t.Header.UnmarshalBinary(b)
	if err != nil {
		return err
	}
	return nil
}

// ParticipationHeader is used for the commands of token participation.
type ParticipationHeader struct {
	Header  *FALinkHeader
	NDN     [10]uint8
	VDN     [10]uint8
	MSN     [10]uint8
	Reserve uint16
}

// Trigger is a trigger frame of FA Link frame.
type Trigger ParticipationHeader

// NewTrigger creates a new Trigger.
func NewTrigger(sna, dna uint8, vseq, seq uint32, ndn, vdn, msn string) *Trigger {
	phdr := &ParticipationHeader{}
	phdr.Header = NewFALinkHeader(
		[4]uint8{0x46, 0x41, 0x43, 0x4e}, // H_TYPE
		0,                                // TFL
		sna,                              // SNA
		dna,                              // DNA
		vseq,                             // V_SEQ
		seq,                              // SEQ
		false, false, false,              // M_CTL
		0, 0, // ULS, M_SZ
		0,          // M_ADD
		0x0a, 0, 0, // MFT, M_RLT, reserved
		TCDTrigger, 0, // TCD, VER
		0, 0x04, // C_AD1, C_SZ1
		0, 0x40, // C_AD2, C_SZ2
		0, 3, true, 0x80, 0, // MODE, P_TYPE, PRI
		1, 1, 0x40, // CBN, TBN, BSIZE
		0, 0x32, 0, // LKS, TW, RCT
	)
	t := Trigger(*phdr)
	t.Header.TFL = uint32(t.MarshalLen())
	t.Header.BSize = uint16(t.MarshalLen())

	t.NDN = [10]uint8{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	t.VDN = [10]uint8{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	t.MSN = [10]uint8{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}

	copy(t.NDN[:len(ndn)], []uint8(ndn))
	copy(t.VDN[:], []uint8(vdn))
	copy(t.MSN[:], []uint8(msn))

	return &t
}

// MarshalBinary returns the byte sequence generated from a Trigger.
func (t *Trigger) MarshalBinary() ([]byte, error) {
	b := make([]byte, t.MarshalLen())
	if err := t.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

//MarshalTo puts the byte sequence in the byte array given as b.
func (t *Trigger) MarshalTo(b []byte) error {
	l := len(b)
	if l < t.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	err := t.Header.MarshalTo(b)
	if err != nil {
		return err
	}

	offset := t.Header.MarshalLen()
	copy(b[offset:], t.NDN[:])
	offset += len(t.NDN)

	copy(b[offset:], t.VDN[:])
	offset += len(t.VDN)

	copy(b[offset:], t.MSN[:])

	return nil
}

// MarshalLen returns the serial length of Trigger
func (t *Trigger) MarshalLen() int {
	return 96
}

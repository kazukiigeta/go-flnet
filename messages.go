// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// FLnet is an interface that defines FL-net messages.
type FLnet interface {
	MarshalBinary() ([]byte, error)
	MarshalTo(b []byte) error
	MarshalLen() int
	UnmarshalBinary(b []byte) error
}

// Parse decodes the given bytes.
// This function checks the TCD.
func Parse(b []byte) (FLnet, error) {
	if len(b) < 64 {
		return nil, ErrTooShortToParse
	}

	var f FLnet
	t := binary.BigEndian.Uint16(b[40:42])

	switch t {
	// Transfer Messages
	case TCDToken:
		f = &Token{
			Header: &FALinkHeader{},
		}
	case TCDTrigger:
		f = &Trigger{
			ParticipationHeader: &ParticipationHeader{
				Header: &FALinkHeader{},
			},
		}
	default:
		// If the combination of class and type is unknown or not supported, *Generic is used.
		return nil, ErrNotImplemented
	}

	if err := f.UnmarshalBinary(b); err != nil {
		return nil, errors.Wrap(err, "failed to decode FLnet")
	}
	return f, nil
}

// Token is a token frame of FA Link frame.
type Token struct {
	Header *FALinkHeader
}

// NewToken creates a new Token.
func NewToken() *Token {
	t := &Token{
		Header: NewFALinkHeader(
			[4]byte{0x46, 0x41, 0x43, 0x4e}, // H_TYPE
			0x40,                            // TFL
			0x1,                             // SNA
			0x55,                            // DNA
			0,                               // V_SEQ
			0,                               // SEQ
			false, false, false,             // M_CTL
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

//MarshalTo puts the byte sequence in the byte array given as b.
func (t *Token) MarshalTo(b []byte) error {
	if err := t.Header.MarshalTo(b); err != nil {
		return err
	}

	return nil
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
	NDN     [10]byte
	VDN     [10]byte
	MSN     [10]byte
	Reserve uint16
}

// MarshalLen returns the serial length of Trigger
func (p *ParticipationHeader) MarshalLen() int {
	return 96
}

// MarshalBinary returns the byte sequence generated from a participation header.
func (p *ParticipationHeader) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.MarshalLen())
	if err := p.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

//MarshalTo puts the byte sequence in the byte array given as b.
func (p *ParticipationHeader) MarshalTo(b []byte) error {
	l := len(b)
	if l < p.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	err := p.Header.MarshalTo(b)
	if err != nil {
		return err
	}

	offset := p.Header.MarshalLen()
	copy(b[offset:], p.NDN[:])
	offset += len(p.NDN)

	copy(b[offset:], p.VDN[:])
	offset += len(p.VDN)

	copy(b[offset:], p.MSN[:])

	return nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a participation header.
func (p *ParticipationHeader) UnmarshalBinary(b []byte) error {
	err := p.Header.UnmarshalBinary(b)
	if err != nil {
		return err
	}

	offset := p.Header.MarshalLen()
	copy(p.NDN[:], b[offset:offset+len(p.NDN)])
	offset += len(p.NDN)

	copy(p.VDN[:], b[offset:offset+len(p.VDN)])
	offset += len(p.VDN)

	copy(p.MSN[:], b[offset:offset+len(p.MSN)])

	return nil
}

// Trigger is a trigger frame of FA Link frame.
type Trigger struct {
	*ParticipationHeader
}

// NewTrigger creates a new Trigger.
func NewTrigger(sna, dna uint8, vseq, seq uint32, ndn, vdn, msn string) *Trigger {
	t := &Trigger{}
	t.ParticipationHeader = &ParticipationHeader{}
	t.ParticipationHeader.Header = NewFALinkHeader(
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
	t.ParticipationHeader.Header.TFL = uint32(t.MarshalLen())
	t.ParticipationHeader.Header.BSize = uint16(t.MarshalLen())

	t.ParticipationHeader.NDN = [10]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	t.ParticipationHeader.VDN = [10]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	t.ParticipationHeader.MSN = [10]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}

	copy(t.NDN[:len(t.NDN)], []byte(ndn))
	copy(t.VDN[:len(t.VDN)], []byte(vdn))
	copy(t.MSN[:len(t.MSN)], []byte(msn))

	return t
}

// MarshalBinary returns the byte sequence generated from a Trigger.
func (t *Trigger) MarshalBinary() ([]byte, error) {
	b := make([]byte, t.MarshalLen())
	if err := t.ParticipationHeader.MarshalTo(b); err != nil {
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

	err := t.ParticipationHeader.MarshalTo(b)
	if err != nil {
		return err
	}

	return nil
}

// MarshalLen returns the serial length of Trigger
func (t *Trigger) MarshalLen() int {
	return t.ParticipationHeader.MarshalLen()
}

// UnmarshalBinary sets the values retrieved from byte sequence in a trigger frame.
func (t *Trigger) UnmarshalBinary(b []byte) error {
	err := t.ParticipationHeader.UnmarshalBinary(b)
	if err != nil {
		return err
	}
	return nil
}

// ParticipationRequest is a ParticipationRequest frame of FA Link frame.
type ParticipationRequest struct {
	*ParticipationHeader
}

// NewParticipationRequest creates a new ParticipationRequest.
func NewParticipationRequest(sna, dna uint8, vseq, seq uint32, ndn, vdn, msn string) *ParticipationRequest {
	p := &ParticipationRequest{}
	p.ParticipationHeader = &ParticipationHeader{}
	p.ParticipationHeader.Header = NewFALinkHeader(
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
		TCDParticipationRequest, 0, // TCD, VER
		0, 0x04, // C_AD1, C_SZ1
		0, 0x40, // C_AD2, C_SZ2
		0, 3, true, 0x80, 0, // MODE, P_TYPE, PRI
		1, 1, 0x40, // CBN, TBN, BSIZE
		0, 0x32, 0, // LKS, TW, RCT
	)
	p.ParticipationHeader.Header.TFL = uint32(p.MarshalLen())
	p.ParticipationHeader.Header.BSize = uint16(p.MarshalLen())

	p.ParticipationHeader.NDN = [10]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	p.ParticipationHeader.VDN = [10]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}
	p.ParticipationHeader.MSN = [10]byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}

	copy(p.NDN[:len(p.NDN)], []byte(ndn))
	copy(p.VDN[:len(p.VDN)], []byte(vdn))
	copy(p.MSN[:len(p.MSN)], []byte(msn))

	return p
}

// MarshalBinary returns the byte sequence generated from a ParticipationRequest.
func (p *ParticipationRequest) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.MarshalLen())
	if err := p.ParticipationHeader.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

//MarshalTo puts the byte sequence in the byte array given as b.
func (p *ParticipationRequest) MarshalTo(b []byte) error {
	l := len(b)
	if l < p.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	err := p.ParticipationHeader.MarshalTo(b)
	if err != nil {
		return err
	}

	return nil
}

// MarshalLen returns the serial length of ParticipationRequest
func (p *ParticipationRequest) MarshalLen() int {
	return p.ParticipationHeader.MarshalLen()
}

// UnmarshalBinary sets the values retrieved from byte sequence in a ParticipationRequest frame.
func (p *ParticipationRequest) UnmarshalBinary(b []byte) error {
	err := p.ParticipationHeader.UnmarshalBinary(b)
	if err != nil {
		return err
	}
	return nil
}

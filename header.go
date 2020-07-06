// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet

import (
	"encoding/binary"
	"io"

	"github.com/kazukiigeta/go-flnet/utils"
)

// TCD deficitions.
const (
	TcdToken = iota + 65000
	TcdCyclic
	TcdParticipateRequest
	TcdByteBlockReadRequest
	TcdByteBlockWriteRequest
)

// FALinkHeader is a FL-net header.
type FALinkHeader struct {
	HType    [4]uint8
	Tfl      uint32
	Sa       uint32
	Da       uint32
	VSeq     uint32
	Seq      uint32
	MCtl     uint32
	Uls      uint16
	MSz      uint16
	MAdd     uint32
	Mft      uint8
	MRlt     uint8
	Reserved uint16
	Tcd      uint16
	Ver      uint16
	CAd1     uint16
	CSz1     uint16
	CAd2     uint16
	CSz2     uint16
	Mode     uint16
	PType    uint8
	Pri      uint8
	Cbn      uint8
	Tbn      uint8
	BSize    uint16
	Lks      uint8
	Tw       uint8
	Rct      uint16
	Payload  []byte
}

//NewFALinkHeader creates a new FALinkHeader.
func NewFALinkHeader(
	htype [4]uint8,
	tfl uint32,
	sna uint8,
	dna uint8,
	vseq uint32,
	seq uint32,
	bct bool,
	ppt bool,
	rpl bool,
	uls uint16,
	msz uint16,
	madd uint32,
	mft uint8,
	mrlt uint8,
	reserved uint16,
	tcd uint16,
	ver uint16,
	cad1 uint16,
	csz1 uint16,
	cad2 uint16,
	csz2 uint16,
	minver uint,
	majver uint,
	tokmode bool,
	ptype uint8,
	pri uint8,
	cbn uint8,
	tbn uint8,
	bsize uint16,
	lks uint8,
	tw uint8,
	rct uint16,
	payload []byte,
) *FALinkHeader {
	h := &FALinkHeader{}
	h.HType = htype
	h.Tfl = tfl
	h.Sa = 0x00010000 | uint32(sna)
	h.Da = 0x00010000 | uint32(dna)
	h.VSeq = vseq
	h.Seq = seq
	h.MCtl = uint32((utils.BoolToUint(bct) << 31) + (utils.BoolToUint(ppt) << 30))
	h.Uls = uls
	h.MSz = msz
	h.MAdd = madd
	h.Mft = mft
	h.MRlt = mrlt
	h.Reserved = reserved
	h.Tcd = tcd
	h.Ver = ver
	h.CAd1 = cad1
	h.CSz1 = csz1
	h.CAd2 = cad2
	h.CSz2 = csz2
	h.Mode = uint16((minver << 8) + (majver << 4) + utils.BoolToUint(tokmode))
	h.PType = ptype
	h.Cbn = cbn
	h.Tbn = tbn
	h.BSize = bsize
	h.Lks = lks
	h.Tw = tw
	h.Rct = rct
	h.Payload = payload

	return h
}

// MarshalBinary returns the byte sequence generated from a FALinkHeader instance.
func (h *FALinkHeader) MarshalBinary() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

//MarshalTo puts the byte sequence in the byte array given as b.
func (h *FALinkHeader) MarshalTo(b []byte) error {
	copy(b, h.HType[:])
	binary.BigEndian.PutUint32(b[4:], h.Tfl)
	binary.BigEndian.PutUint32(b[8:], h.Sa)
	binary.BigEndian.PutUint32(b[12:], h.Da)
	binary.BigEndian.PutUint32(b[16:], h.VSeq)
	binary.BigEndian.PutUint32(b[20:], h.Seq)
	binary.BigEndian.PutUint32(b[24:], h.MCtl)
	binary.BigEndian.PutUint16(b[28:], h.Uls)
	binary.BigEndian.PutUint16(b[30:], h.MSz)
	binary.BigEndian.PutUint32(b[32:], h.MAdd)
	b[36] = h.Mft
	b[37] = h.MRlt
	binary.BigEndian.PutUint16(b[38:], h.Reserved)
	binary.BigEndian.PutUint16(b[40:], h.Tcd)
	binary.BigEndian.PutUint16(b[42:], h.Ver)
	binary.BigEndian.PutUint16(b[44:], h.CAd1)
	binary.BigEndian.PutUint16(b[46:], h.CSz1)
	binary.BigEndian.PutUint16(b[48:], h.CAd2)
	binary.BigEndian.PutUint16(b[50:], h.CSz2)
	binary.BigEndian.PutUint16(b[52:], h.Mode)
	b[54] = h.PType
	b[55] = h.Pri
	b[56] = h.Cbn
	b[57] = h.Tbn
	binary.BigEndian.PutUint16(b[58:], h.BSize)
	b[60] = h.Lks
	b[61] = h.Tw
	binary.BigEndian.PutUint16(b[62:], h.Rct)
	copy(b[64:h.MarshalLen()], h.Payload)

	return nil
}

// MarshalLen returns the serial length.
func (h *FALinkHeader) MarshalLen() int {
	return 64 + len(h.Payload)
}

// ParseHeader decodes given byte sequence as a FL-net common header.
func ParseHeader(b []byte) (*FALinkHeader, error) {
	h := &FALinkHeader{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a FL-net common header.
func (h *FALinkHeader) UnmarshalBinary(b []byte) error {
	if len(b) < 64 {
		return io.ErrUnexpectedEOF
	}
	copy(h.HType[:], b[:4])
	h.Tfl = binary.BigEndian.Uint32(b[4:8])
	h.Sa = binary.BigEndian.Uint32(b[8:12])
	h.Da = binary.BigEndian.Uint32(b[12:16])
	h.VSeq = binary.BigEndian.Uint32(b[16:20])
	h.Seq = binary.BigEndian.Uint32(b[20:24])
	h.MCtl = binary.BigEndian.Uint32(b[24:28])
	h.Uls = binary.BigEndian.Uint16(b[28:30])
	h.MSz = binary.BigEndian.Uint16(b[30:32])
	h.MAdd = binary.BigEndian.Uint32(b[32:36])
	h.Mft = uint8(b[37])
	h.MRlt = uint8(b[38])
	h.Reserved = binary.BigEndian.Uint16(b[38:40])
	h.Tcd = binary.BigEndian.Uint16(b[40:42])
	h.Ver = binary.BigEndian.Uint16(b[42:44])
	h.CAd1 = binary.BigEndian.Uint16(b[44:46])
	h.CSz1 = binary.BigEndian.Uint16(b[46:48])
	h.CAd2 = binary.BigEndian.Uint16(b[48:50])
	h.CSz2 = binary.BigEndian.Uint16(b[50:52])
	h.Mode = binary.BigEndian.Uint16(b[52:54])
	h.PType = uint8(b[55])
	h.Pri = uint8(b[56])
	h.Cbn = uint8(b[57])
	h.Tbn = uint8(b[58])
	h.BSize = binary.BigEndian.Uint16(b[58:60])
	h.Tbn = uint8(b[61])
	h.Tbn = uint8(b[62])
	h.Rct = binary.BigEndian.Uint16(b[62:64])
	h.Payload = b[64:]
	return nil
}

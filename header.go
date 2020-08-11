// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet

import (
	"encoding/binary"
	"io"

	"github.com/kazukiigeta/go-flnet/utils"
)

// TCD definitions.
const (
	TCDToken uint16 = iota + 65000
	TCDCyclic
	TCDParticipateRequest
	TCDByteBlockReadRequest
	TCDByteBlockWriteRequest
	TCDWordBlockReadRequest
	TCDWordBlockWriteRequest
	TCDNetworkParameterReadRequest
	TCDNetworkParameterWriteRequest
	TCDStopCommandRequest
	TCDOperationCommandRequest
	TCDProfileReadRequest
	TCDTrigger
)

// FALinkHeader is a FL-net header.
type FALinkHeader struct {
	HType    [4]uint8
	TFL      uint32
	SA       uint32
	DA       uint32
	VSeq     uint32
	Seq      uint32
	MCTL     uint32
	ULS      uint16
	MSZ      uint16
	MADD     uint32
	MFT      uint8
	MRLT     uint8
	Reserved uint16
	TCD      uint16
	Ver      uint16
	CAD1     uint16
	CSZ1     uint16
	CAD2     uint16
	CSZ2     uint16
	Mode     uint16
	PType    uint8
	Pri      uint8
	CBN      uint8
	TBN      uint8
	BSize    uint16
	LKS      uint8
	TW       uint8
	RCT      uint16
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
) *FALinkHeader {
	return &FALinkHeader{
		HType:    htype,
		TFL:      tfl,
		SA:       0x00010000 | uint32(sna),
		DA:       0x00010000 | uint32(dna),
		VSeq:     vseq,
		Seq:      seq,
		MCTL:     uint32((utils.BoolToUint(bct) << 31) + (utils.BoolToUint(ppt) << 30)),
		ULS:      uls,
		MSZ:      msz,
		MADD:     madd,
		MFT:      mft,
		MRLT:     mrlt,
		Reserved: reserved,
		TCD:      tcd,
		Ver:      ver,
		CAD1:     cad1,
		CSZ1:     csz1,
		CAD2:     cad2,
		CSZ2:     csz2,
		Mode:     uint16((minver << 8) + (majver << 4) + utils.BoolToUint(tokmode)),
		PType:    ptype,
		CBN:      cbn,
		TBN:      tbn,
		BSize:    bsize,
		LKS:      lks,
		TW:       tw,
		RCT:      rct,
	}
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
	l := len(b)
	if l < 64 {
		return io.ErrUnexpectedEOF
	}

	copy(b, h.HType[:])
	binary.BigEndian.PutUint32(b[4:], h.TFL)
	binary.BigEndian.PutUint32(b[8:], h.SA)
	binary.BigEndian.PutUint32(b[12:], h.DA)
	binary.BigEndian.PutUint32(b[16:], h.VSeq)
	binary.BigEndian.PutUint32(b[20:], h.Seq)
	binary.BigEndian.PutUint32(b[24:], h.MCTL)
	binary.BigEndian.PutUint16(b[28:], h.ULS)
	binary.BigEndian.PutUint16(b[30:], h.MSZ)
	binary.BigEndian.PutUint32(b[32:], h.MADD)
	b[36] = h.MFT
	b[37] = h.MRLT
	binary.BigEndian.PutUint16(b[38:], h.Reserved)
	binary.BigEndian.PutUint16(b[40:], h.TCD)
	binary.BigEndian.PutUint16(b[42:], h.Ver)
	binary.BigEndian.PutUint16(b[44:], h.CAD1)
	binary.BigEndian.PutUint16(b[46:], h.CSZ1)
	binary.BigEndian.PutUint16(b[48:], h.CAD2)
	binary.BigEndian.PutUint16(b[50:], h.CSZ2)
	binary.BigEndian.PutUint16(b[52:], h.Mode)
	b[54] = h.PType
	b[55] = h.Pri
	b[56] = h.CBN
	b[57] = h.TBN
	binary.BigEndian.PutUint16(b[58:], h.BSize)
	b[60] = h.LKS
	b[61] = h.TW
	binary.BigEndian.PutUint16(b[62:], h.RCT)

	return nil
}

// MarshalLen returns the serial length.
func (h *FALinkHeader) MarshalLen() int {
	return 64
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
	h.TFL = binary.BigEndian.Uint32(b[4:8])
	h.SA = binary.BigEndian.Uint32(b[8:12])
	h.DA = binary.BigEndian.Uint32(b[12:16])
	h.VSeq = binary.BigEndian.Uint32(b[16:20])
	h.Seq = binary.BigEndian.Uint32(b[20:24])
	h.MCTL = binary.BigEndian.Uint32(b[24:28])
	h.ULS = binary.BigEndian.Uint16(b[28:30])
	h.MSZ = binary.BigEndian.Uint16(b[30:32])
	h.MADD = binary.BigEndian.Uint32(b[32:36])
	h.MFT = uint8(b[37])
	h.MRLT = uint8(b[38])
	h.Reserved = binary.BigEndian.Uint16(b[38:40])
	h.TCD = binary.BigEndian.Uint16(b[40:42])
	h.Ver = binary.BigEndian.Uint16(b[42:44])
	h.CAD1 = binary.BigEndian.Uint16(b[44:46])
	h.CSZ1 = binary.BigEndian.Uint16(b[46:48])
	h.CAD2 = binary.BigEndian.Uint16(b[48:50])
	h.CSZ2 = binary.BigEndian.Uint16(b[50:52])
	h.Mode = binary.BigEndian.Uint16(b[52:54])
	h.PType = uint8(b[55])
	h.Pri = uint8(b[56])
	h.CBN = uint8(b[57])
	h.TBN = uint8(b[58])
	h.BSize = binary.BigEndian.Uint16(b[58:60])
	h.TBN = uint8(b[61])
	h.TBN = uint8(b[62])
	h.RCT = binary.BigEndian.Uint16(b[62:64])

	return nil
}

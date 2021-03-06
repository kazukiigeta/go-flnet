// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet_test

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kazukiigeta/go-flnet"
)

func TestToken(t *testing.T) {
	var testcases = []testCase{
		{
			description: "Token frame",
			structured:  flnet.NewToken(),
			serialized: []byte{
				0x46, 0x41, 0x43, 0x4e, // H_TYPE
				0x00, 0x00, 0x00, 0x40, // TFL
				0x00, 0x01, 0x00, 0x01, // SA
				0x00, 0x01, 0x00, 0x55, // DA
				0x00, 0x00, 0x00, 0x00, // V_SEQ
				0x00, 0x00, 0x00, 0x00, // SEQ
				0x00, 0x00, 0x00, 0x00, // M_CTL
				0x00, 0x00, 0x00, 0x00, // ULS, M_SZ
				0x00, 0x00, 0x00, 0x00, // M_ADD
				0x00, 0x00, 0x00, 0x00, // MFT, M_RLT, reserved
				0xfd, 0xe8, 0x00, 0x00, // TCD, VER
				0x00, 0x00, 0x00, 0x00, // C_AD1, C_SZ1
				0x00, 0x00, 0x00, 0x00, // C_AD2, C_SZ2
				0x00, 0x31, 0x80, 0x00, // MODE, P_TYPE, PRI
				0x01, 0x01, 0x00, 0x40, // CBN, TBN, BSIZE
				0x00, 0x32, 0x00, 0x00, // LKS, TW, RCT
			},
		},
	}

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := flnet.Parse(c.serialized)
				if err != nil {
					t.Fatal(err)
				}
				got, want := msg, c.structured
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				got, want := b, c.serialized
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
		})
	}
}

func TestTrigger(t *testing.T) {
	var testcases = []testCase{
		{
			description: "Trigger frame",
			structured:  flnet.NewTrigger(1, 255, 0, 0, "NODE", "VENDOR", "MANUF."),
			serialized: []byte{
				0x46, 0x41, 0x43, 0x4e, // H_TYPE
				0x00, 0x00, 0x00, 0x60, // TFL
				0x00, 0x01, 0x00, 0x01, // SA
				0x00, 0x01, 0x00, 0xff, // DA
				0x00, 0x00, 0x00, 0x00, // V_SEQ
				0x00, 0x00, 0x00, 0x00, // SEQ
				0x00, 0x00, 0x00, 0x00, // M_CTL
				0x00, 0x00, 0x00, 0x00, // ULS, M_SZ
				0x00, 0x00, 0x00, 0x00, // M_ADD
				0x0a, 0x00, 0x00, 0x00, // MFT, M_RLT, reserved
				0xfd, 0xf4, 0x00, 0x00, // TCD, VER
				0x00, 0x00, 0x00, 0x04, // C_AD1, C_SZ1
				0x00, 0x00, 0x00, 0x40, // C_AD2, C_SZ2
				0x00, 0x31, 0x80, 0x00, // MODE, P_TYPE, PRI
				0x01, 0x01, 0x00, 0x60, // CBN, TBN, BSIZE
				0x00, 0x32, 0x00, 0x00, // LKS, TW, RCT
				0x4e, 0x4f, 0x44, 0x45, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, // NDN
				0x56, 0x45, 0x4e, 0x44, 0x4f, 0x52, 0x20, 0x20, 0x20, 0x20, // VDN
				0x4d, 0x41, 0x4e, 0x55, 0x46, 0x2e, 0x20, 0x20, 0x20, 0x20, // MSN
				0x00, 0x00, // Reserved
			},
		},
	}

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := flnet.Parse(c.serialized)
				if err != nil {
					t.Fatal(err)
				}
				got, want := msg, c.structured
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				if got, want := b, c.serialized; !reflect.DeepEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
				}
			})
		})
	}
}

func TestCyclic(t *testing.T) {
	data := make([]byte, 136)
	var testcases = []testCase{
		{
			description: "Cyclic frame with data",
			structured: flnet.NewCyclic(
				0x55, 0x01, 0x0c466d, 4, 4, 64, 64, &data,
			),
			serialized: []byte{
				0x46, 0x41, 0x43, 0x4e, // H_TYPE
				0x00, 0x00, 0x00, 0xc8, // TFL
				0x00, 0x01, 0x00, 0x55, // SA
				0x00, 0x01, 0x00, 0x01, // DA
				0x00, 0x0c, 0x46, 0x6d, // V_SEQ
				0x00, 0x00, 0x00, 0x00, // SEQ
				0x00, 0x00, 0x00, 0x00, // M_CTL
				0x00, 0x00, 0x00, 0x00, // ULS, M_SZ
				0x00, 0x00, 0x00, 0x00, // M_ADD
				0x0a, 0x00, 0x00, 0x00, // MFT, M_RLT, reserved
				0xfd, 0xe9, 0x00, 0x00, // TCD, VER
				0x00, 0x04, 0x00, 0x04, // C_AD1, C_SZ1
				0x00, 0x40, 0x00, 0x40, // C_AD2, C_SZ2
				0x00, 0x31, 0x80, 0x00, // MODE, P_TYPE, PRI
				0x01, 0x01, 0x00, 0xc8, // CBN, TBN, BSIZE
				0x00, 0x32, 0x00, 0x00, // LKS, TW, RCT
				0x00, 0x00, 0x00, 0x00, // Data
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			description: "Cyclic frame without data",
			structured: flnet.NewCyclic(
				0x55, 0x01, 0x085d, 4, 4, 64, 64, nil,
			),
			serialized: []byte{
				0x46, 0x41, 0x43, 0x4e, // H_TYPE
				0x00, 0x00, 0x00, 0xc8, // TFL
				0x00, 0x01, 0x00, 0x55, // SA
				0x00, 0x01, 0x00, 0x01, // DA
				0x00, 0x00, 0x08, 0x5d, // V_SEQ
				0x00, 0x00, 0x00, 0x00, // SEQ
				0x00, 0x00, 0x00, 0x00, // M_CTL
				0x00, 0x00, 0x00, 0x00, // ULS, M_SZ
				0x00, 0x00, 0x00, 0x00, // M_ADD
				0x0a, 0x00, 0x00, 0x00, // MFT, M_RLT, reserved
				0xfd, 0xe9, 0x00, 0x00, // TCD, VER
				0x00, 0x04, 0x00, 0x04, // C_AD1, C_SZ1
				0x00, 0x40, 0x00, 0x40, // C_AD2, C_SZ2
				0x00, 0x31, 0x80, 0x00, // MODE, P_TYPE, PRI
				0x01, 0x01, 0x00, 0x40, // CBN, TBN, BSIZE
				0x00, 0x32, 0x00, 0x00, // LKS, TW, RCT
			},
		},
	}

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := flnet.Parse(c.serialized)
				if err != nil {
					t.Fatal(err)
				}
				got, want := msg, c.structured
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				got, want := b, c.serialized
				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("differs: (-want +got)\n%s", diff)
				}
			})
		})
	}
}

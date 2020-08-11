// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet_test

import (
	"reflect"
	"testing"

	"github.com/kazukiigeta/go-flnet"
)

var testcases = []testCase{
	{
		description: "FA link header",
		structured: flnet.NewFALinkHeader(
			[4]uint8{0x46, 0x41, 0x43, 0x4e}, // H_TYPE
			0x40,                             // TFL
			1,                                // SA
			0xff,                             // DA
			0,                                // V_SEQ
			0,                                // SEQ
			true, false, false,               // M_CTL
			0, 0, // ULS, M_SZ
			0,       // M_ADD
			0, 0, 0, // MFT, M_RLT, reserved
			0xfdec, 0, // TCD, VER
			0, 0, // C_AD1, C_SZ1
			0, 0, // C_AD2, C_SZ2
			0, 3, true, 0, 0, // MODE, P_TYPE, PRI
			0, 0, 0x40, // CBN, TBN, BSIZE
			0, 0, 0, // LKS, TW, RCT
		),
		serialized: []byte{
			0x46, 0x41, 0x43, 0x4e, // H_TYPE
			0x00, 0x00, 0x00, 0x40, // TFL
			0x00, 0x01, 0x00, 0x01, // SA
			0x00, 0x01, 0x00, 0xff, // DA
			0x00, 0x00, 0x00, 0x00, // V_SEQ
			0x00, 0x00, 0x00, 0x00, // SEQ
			0x80, 0x00, 0x00, 0x00, // M_CTL
			0x00, 0x00, 0x00, 0x00, // ULS, M_SZ
			0x00, 0x00, 0x00, 0x00, // M_ADD
			0x00, 0x00, 0x00, 0x00, // MFT, M_RLT, reserved
			0xfd, 0xec, 0x00, 0x00, // TCD, VER
			0x00, 0x00, 0x00, 0x00, // C_AD1, C_SZ1
			0x00, 0x00, 0x00, 0x00, // C_AD2, C_SZ2
			0x00, 0x31, 0x00, 0x00, // MODE, P_TYPE, PRI
			0x00, 0x00, 0x00, 0x40, // CBN, TBN, BSIZE
			0x00, 0x00, 0x00, 0x00, // LKS, TW, RCT
		},
	},
}

func TestMessages(t *testing.T) {
	t.Helper()
	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := flnet.ParseHeader(c.serialized)
				if err != nil {
					t.Fatal(err)
				}
				if got, want := msg, c.structured; !reflect.DeepEqual(got, want) {
					t.Errorf("got %v, want %v", got, want)
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

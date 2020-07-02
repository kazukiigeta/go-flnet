// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package utils_test

import (
	"testing"

	"github.com/kazukiigeta/go-flnet/utils"
)

func TestBoolToUint(t *testing.T) {
	cases := []struct {
		description string
		b           bool
		u           uint
	}{
		{
			"case of true",
			true,
			1,
		},
		{
			"case of false",
			false,
			0,
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if got, want := utils.BoolToUint(c.b), c.u; got != want {
				t.Fail()
			}

		})
	}

}

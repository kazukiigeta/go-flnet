// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet_test

type serializeable interface {
	MarshalBinary() ([]byte, error)
	MarshalLen() int
}

type testCase = struct {
	description string
	structured  serializeable
	serialized  []byte
}

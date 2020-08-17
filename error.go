// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package flnet

import "errors"

// Error definitions.
var (
	ErrTooShortToMarshalBinary = errors.New("insufficient buffer to serialize parameter to")
	ErrTooShortToParse         = errors.New("too short to decode as parameter")
	ErrNotImplemented          = errors.New("not implemented")
)

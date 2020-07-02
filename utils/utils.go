// Copyright 2020 go-flnet authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Package utils provides some utilities which might be useful for protcol stack.
package utils

// BoolToUint converts bool to uint.
func BoolToUint(b bool) uint {
	if b == true {
		return 1
	}
	return 0
}

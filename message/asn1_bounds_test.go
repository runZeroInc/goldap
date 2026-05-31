// Copyright 2026 runZero, Inc. All rights reserved.
//
// Regression test: parseTagAndLength read the first tag byte
// (bytes[offset]) before verifying the offset was within the buffer, so an
// exhausted/truncated buffer caused an index-out-of-range panic. After the fix
// it must return a parse error instead of panicking.

package message

import "testing"

// TestParseTagAndLengthExhaustedBuffer calls parseTagAndLength at an offset
// equal to len(bytes) (the buffer is exhausted). Before the fix the initial
// bytes[offset] read panicked.
func TestParseTagAndLengthExhaustedBuffer(t *testing.T) {
	buf := []byte{0x30} // a single byte; offset 1 is past the end

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parseTagAndLength panicked on exhausted buffer: %v", r)
		}
	}()

	if _, _, err := parseTagAndLength(buf, len(buf)); err == nil {
		t.Fatalf("expected a parse error for an exhausted buffer, got nil")
	}
}

// TestParseTagAndLengthEmptyBuffer calls parseTagAndLength on an empty
// buffer at offset 0.
func TestParseTagAndLengthEmptyBuffer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parseTagAndLength panicked on empty buffer: %v", r)
		}
	}()
	if _, _, err := parseTagAndLength([]byte{}, 0); err == nil {
		t.Fatalf("expected a parse error for an empty buffer, got nil")
	}
}

// TestParseTagAndLengthTruncatedLength confirms a tag with no following
// length byte is still rejected (existing later guard) without panicking.
func TestParseTagAndLengthTruncatedLength(t *testing.T) {
	buf := []byte{0x04} // OctetString tag, but no length byte follows
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parseTagAndLength panicked on truncated length: %v", r)
		}
	}()
	if _, _, err := parseTagAndLength(buf, 0); err == nil {
		t.Fatalf("expected a parse error for a truncated length, got nil")
	}
}

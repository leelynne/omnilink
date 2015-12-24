package omni

import "testing"

func TestDeser(t *testing.T) {
	m := genmsg{
		SeqNum: 1,
		Type:   ClientReqNewSession,
	}
	b := m.serialize()
	expected := []byte{0x1, 0x0, 0x1, 0x0}
	if !eq(b, expected) {
		t.Errorf("Bad serialization %v %v", b, expected)
	}
}

func eq(a []byte, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

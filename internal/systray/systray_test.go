package systray

import (
	"bytes"
	"testing"
)

func TestGetIcon(t *testing.T) {
	icon := getIcon()

	// Icon should not be empty
	if len(icon) == 0 {
		t.Error("getIcon() returned empty slice")
	}

	// Icon should be valid PNG (starts with PNG signature)
	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if !bytes.HasPrefix(icon, pngSignature) {
		t.Error("getIcon() did not return valid PNG data")
	}

	// PNG should end with IEND chunk
	iendSignature := []byte{0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82}
	if !bytes.HasSuffix(icon, iendSignature) {
		t.Error("getIcon() PNG data missing IEND chunk")
	}
}

func TestGetIconConsistency(t *testing.T) {
	// Multiple calls should return the same icon
	icon1 := getIcon()
	icon2 := getIcon()

	if !bytes.Equal(icon1, icon2) {
		t.Error("getIcon() should return consistent results")
	}
}

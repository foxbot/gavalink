package gavalink_test

import (
	"testing"

	"github.com/foxbot/gavalink"
)

func TestDecoder(t *testing.T) {
	data := "QAAAkAIALGxvZmkgaGlwIGhvcCByYWRpbyAtIGJlYXRzIHRvIHJlbGF4L3N0dWR5IHRvAApDaGlsbGVkQ293f/////////8AC2hIVzFvWTI2a3hRAQEAK2h0dHBzOi8vd3d3LnlvdXR1YmUuY29tL3dhdGNoP3Y9aEhXMW9ZMjZreFEAB3lvdXR1YmUAAAAAAAAAAA=="
	track, err := gavalink.DecodeString(data)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(track)
}

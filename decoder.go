package gavalink

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"io"
)

const trackInfoVersioned int32 = 1

// DecodeString decodes a base64 Lavaplayer string to a TrackInfo
func DecodeString(data string) (*TrackInfo, error) {
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	return Decode(r)
}

// Decode decodes a reader into a TrackInfo
func Decode(r io.Reader) (*TrackInfo, error) {
	// https://github.com/serenity-rs/lavalink.rs/blob/master/src/decoder.rs

	var value uint8
	if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
		return nil, err
	}

	flags := int32(int64(value) & 0xC00000000)

	// irrelevant data?
	var ignore [2]byte
	if err := binary.Read(r, binary.LittleEndian, &ignore); err != nil {
		return nil, err
	}

	var version uint8
	if flags&trackInfoVersioned == 0 {
		version = 1
	} else {
		if err := binary.Read(r, binary.LittleEndian, &version); err != nil {
			return nil, err
		}
	}

	if err := binary.Read(r, binary.LittleEndian, &ignore); err != nil {
		return nil, err
	}

	title, err := readString(r)
	if err != nil {
		return nil, err
	}

	author, err := readString(r)
	if err != nil {
		return nil, err
	}

	var length uint64
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	identifier, err := readString(r)
	if err != nil {
		return nil, err
	}

	var stream uint8
	if err := binary.Read(r, binary.LittleEndian, &stream); err != nil {
		return nil, err
	}

	var hasURL uint8
	if err := binary.Read(r, binary.LittleEndian, &hasURL); err != nil {
		return nil, err
	}

	var url string
	if hasURL == 1 {
		url, err = readString(r)
		if err != nil {
			return nil, err
		}
	} else {
		var size uint8
		if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
			return nil, err
		}

		ignore := make([]byte, size)
		if err := binary.Read(r, binary.LittleEndian, &ignore); err != nil {
			return nil, err
		}
	}

	/*source, err := readString(r)
	if err != nil {
		return nil, err
	}*/

	track := &TrackInfo{
		Identifier: identifier,
		Title:      title,
		Author:     author,
		URI:        url,
		Stream:     stream == 1,
		Length:     int(length),
	}

	return track, nil
}

func readString(r io.Reader) (string, error) {
	var size uint16
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return "", err
	}
	buf := make([]byte, size)
	if err := binary.Read(r, binary.BigEndian, &buf); err != nil {
		return "", err
	}

	return string(buf), nil
}

// Package core implements the EV1 decoding logic.
//
// An EV1 file is nothing more than a regular FLV video whose first
// HeaderLen bytes have been obfuscated with a XOR 0xFF mask. Decoding
// therefore only needs to restore those leading bytes and rename the
// file so that players recognise it as an FLV.
package core

import (
	"io"
	"os"
)

const (
	// EV1Ext is the file extension that identifies an EV1 video.
	EV1Ext = ".ev1"

	// HeaderLen is the number of leading bytes that are obfuscated in
	// the EV1 format. The rest of the file is already plain FLV data.
	HeaderLen = 100

	// xorKey is the byte mask used to obfuscate the EV1 header.
	// Applying it twice restores the original bytes.
	xorKey = 0xFF
)

// DecryptFile decodes a single EV1 file in place.
//
// It XORs the first HeaderLen bytes with 0xFF, writes them back and then
// renames the file by appending ".flv" (e.g. "video.ev1" -> "video.ev1.flv").
// The new path is returned on success.
func DecryptFile(path string) (string, error) {
	f, err := os.OpenFile(path, os.O_RDWR, 0o644)
	if err != nil {
		return "", err
	}

	// Read up to HeaderLen bytes. Short files simply yield fewer bytes,
	// mirroring the behaviour of Python's f.read(100).
	buf := make([]byte, HeaderLen)
	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		f.Close()
		return "", err
	}

	// Unmask the header.
	for i := 0; i < n; i++ {
		buf[i] ^= xorKey
	}

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		f.Close()
		return "", err
	}
	if _, err := f.Write(buf[:n]); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	newPath := path + ".flv"
	if err := os.Rename(path, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

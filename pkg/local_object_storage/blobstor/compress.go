package blobstor

import (
	"github.com/klauspost/compress/zstd"
)

func noOpCompressor(data []byte) []byte {
	return data
}

func noOpDecompressor(data []byte) ([]byte, error) {
	return data, nil
}

func zstdCompressor() (func([]byte) []byte, error) {
	enc, err := zstd.NewWriter(nil)
	if err != nil {
		return nil, err
	}

	return func(data []byte) []byte {
		return enc.EncodeAll(data, make([]byte, 0, len(data)))
	}, nil
}

func zstdDecompressor() (func([]byte) ([]byte, error), error) {
	dec, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}

	return func(data []byte) ([]byte, error) {
		return dec.DecodeAll(data, nil)
	}, nil
}

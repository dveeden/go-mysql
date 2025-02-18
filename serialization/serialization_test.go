package serialization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrailingOneBitCount(t *testing.T) {
	testcases := []struct {
		input  byte
		result int
	}{
		{0b00000000, 0},
		{0b00000001, 1},
		{0b00000011, 2},
		{0b00000111, 3},
		{0b00001111, 4},
		{0b00011111, 5},
		{0b00111111, 6},
		{0b01111111, 7},
		{0b11111111, 8},
		{0b10000000, 0},
		{0b11111101, 1},
	}

	for _, tc := range testcases {
		actual := trailingOneBitCount(tc.input)
		require.Equal(t, tc.result, actual)
	}
}

func TestDecodeFixed(t *testing.T) {
	testcases := []struct {
		input  []byte
		len    int
		result []byte
		err    string
	}{
		{
			[]byte{0xee, 0x81, 0x02, 0xc1, 0x02, 0x01, 0x03, 0x41, 0x03, 0x81, 0x03, 0xc1, 0x03, 0xc5, 0x03, 0x22,
				0x22, 0xee, 0xfd, 0x03, 0xee, 0xfd, 0x03, 0xee, 0xfd, 0x03},
			16,
			[]byte{0x77, 0xa0, 0xb0, 0xc0, 0xd0, 0xe0, 0xf0, 0xf1, 0x11, 0x11, 0x77, 0xff, 0x77, 0xff, 0x77, 0xff},
			"",
		},
		{
			[]byte{0xee, 0x81},
			16,
			[]byte{},
			"data truncated",
		},
		{
			[]byte{},
			16,
			[]byte{},
			"data truncated",
		},
		{
			[]byte{0xee, 0x81, 0x04, 0xc1, 0x02, 0x01, 0x03, 0x41, 0x03, 0x81, 0x03, 0xc1, 0x03, 0xc5, 0x03, 0x22,
				0x22, 0xee, 0xfd, 0x03, 0xee, 0xfd, 0x03, 0xee, 0xfd, 0x03},
			16,
			[]byte{},
			"unknown decoding for",
		},
	}

	for _, tc := range testcases {
		actual, _, err := decodeFixed(tc.input, 0, tc.len)
		if tc.err == "" {
			require.NoError(t, err)
			require.Equal(t, tc.result, actual)
			require.Equal(t, tc.len, len(actual))
		} else {
			require.ErrorContains(t, err, tc.err)
		}
	}
}

func TestDecodeString(t *testing.T) {
	testcases := []struct {
		input  []byte
		result string
		err    string
	}{
		{
			[]byte{0x18, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c},
			"abcdefghijkl",
			"",
		},
		{
			[]byte{0x18, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67},
			"",
			"string truncated",
		},
		{
			[]byte{},
			"",
			"string truncated, expected at least one byte",
		},
		{
			[]byte{0x18},
			"",
			"string truncated, expected length",
		},
	}

	for _, tc := range testcases {
		s, _, err := decodeString(tc.input, 0)
		if tc.err == "" {
			require.NoError(t, err)
			require.Equal(t, tc.result, s)
		} else {
			require.ErrorContains(t, err, tc.err)
		}
	}
}

func TestDecodeVar(t *testing.T) {
	testcases := []struct {
		input    []byte
		unsigned bool
		result   interface{}
		err      string
	}{
		{
			[]byte{},
			false,
			0,
			"data truncated",
		},
		{
			[]byte{0xd9},
			false,
			0,
			"truncated data",
		},
		{
			[]byte{0x4},
			false,
			int64(1),
			"",
		},
		{
			[]byte{0xd9, 0x03},
			false,
			int64(123),
			"",
		},
		{
			[]byte{0xc3, 02, 0x0b},
			true,
			uint64(90200),
			"",
		},
		{
			// From the example on https://dev.mysql.com/doc/dev/mysql-server/latest/PageLibsMysqlSerialization.html
			// But converted to LE
			// unsigned integer, 65535
			[]byte{0b11111011, 0b11111111, 0b00000111},
			true,
			uint64(65535),
			"",
		},
		{
			// From the example on https://dev.mysql.com/doc/dev/mysql-server/latest/PageLibsMysqlSerialization.html
			// But converted to LE
			// signed integer, 65535
			[]byte{0b11110011, 0b11111111, 0b00001111},
			false,
			int64(65535),
			"",
		},
		{
			// From the example on https://dev.mysql.com/doc/dev/mysql-server/latest/PageLibsMysqlSerialization.html
			// But converted to LE
			// signed integer, -65535
			[]byte{0b11101011, 0b11111111, 0b00001111},
			false,
			int64(-65535),
			"",
		},
		{
			// From the example on https://dev.mysql.com/doc/dev/mysql-server/latest/PageLibsMysqlSerialization.html
			// But converted to LE
			// signed integer, 65536
			[]byte{0b11111011, 0b11111111, 0b00001111},
			false,
			int64(-65536),
			"",
		},
		{
			[]byte{0x5d, 0x03},
			true,
			uint64(215),
			"",
		},
		{
			[]byte{0x7f, 0x39, 0x7d, 0x89, 0x70, 0xdb, 0x2d, 0x06},
			true,
			uint64(1739270369410361),
			"",
		},
	}

	for _, tc := range testcases {
		r, _, err := decodeVar(tc.input, 0, tc.unsigned)
		if tc.err == "" {
			require.NoError(t, err)
			require.Equal(t, tc.result, r, tc.result)
		} else {
			require.ErrorContains(t, err, tc.err)
		}
	}
}

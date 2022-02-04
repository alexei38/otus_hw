package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoneCopy(t *testing.T) {
	tests := []struct {
		offset int64
		limit  int64
	}{
		{offset: 0, limit: 0},
		{offset: 0, limit: 10},
		{offset: 0, limit: 1000},
		{offset: 0, limit: 10000},
		{offset: 100, limit: 1000},
		{offset: 6000, limit: 1000},
	}
	srcFile := "testdata/input.txt"

	for _, tc := range tests {
		tc := tc
		expectFile := fmt.Sprintf("testdata/out_offset%d_limit%d.txt", tc.offset, tc.limit)
		dstFile := fmt.Sprintf("/tmp/dst_offset%d_limit%d.txt", tc.offset, tc.limit)
		t.Run(expectFile, func(t *testing.T) {
			err := Copy(srcFile, dstFile, tc.offset, tc.limit)
			require.NoError(t, err)
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					t.Fatalf("Failed remove tempfile %s: %s", name, err)
				}
			}(dstFile)

			expect, err := os.ReadFile(expectFile)
			if err != nil {
				t.Fatalf("Error opening file %s: %s", expectFile, err)
			}
			got, err := os.ReadFile(dstFile)
			if err != nil {
				t.Fatalf("Error opening file %s: %s", expectFile, err)
			}
			require.Equal(t, string(expect), string(got))
		})
	}
}

func TestErrCopy(t *testing.T) {
	tests := []struct {
		offset int64
		limit  int64
		error  error
		src    string
		dst    string
	}{
		{
			offset: 0,
			limit:  0,
			error:  ErrUnsupportedFile,
			src:    "/dev/zero",
			dst:    "/tmp/dst_err_zero_offset0_limit0.txt",
		},
		{
			offset: 0,
			limit:  0,
			error:  ErrUnsupportedFile,
			src:    "/dev/random",
			dst:    "/tmp/dst_err_random_offset0_limit0.txt",
		},
		{
			offset: 10000,
			limit:  0,
			error:  ErrOffsetExceedsFileSize,
			src:    "testdata/input.txt",
			dst:    "/tmp/dst_err_offset10000_limit0.txt",
		},
		{
			offset: 0,
			limit:  0,
			error:  os.ErrNotExist,
			src:    "testdata/input.txt",
			dst:    "/not/found/dst/path",
		},
		{
			offset: 0,
			limit:  0,
			error:  os.ErrNotExist,
			src:    "/not/found/src/path",
			dst:    "/tmp/err_not_found_src.txt",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.dst, func(t *testing.T) {
			err := Copy(tc.src, tc.dst, tc.offset, tc.limit)
			require.Error(t, err)
			require.Truef(t, errors.Is(err, tc.error), "actual error %q", err)
			_, err = os.Stat(tc.dst)
			require.Truef(t, errors.Is(err, os.ErrNotExist), "actual error %q", err)
		})
	}
}

package main

import (
	"fmt"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/rhysd/go-fakeio"
	"github.com/stretchr/testify/require"
)

func TestRunCmdCodes(t *testing.T) {
	dirEnv, err := ReadDir(testDataPath)
	require.NoError(t, err, "actual err - %v", err)

	tests := []struct {
		name   string
		args   []string
		expect int
	}{
		{name: "nil args", args: nil, expect: 111},
		{name: "empty args", args: []string{}, expect: 111},
		{name: "cmd", args: []string{"true"}, expect: 0},
		{name: "cmd with args", args: []string{"echo", "true"}, expect: 0},
		{name: "unknown cmd path", args: []string{faker.Name()}, expect: 127},
		{name: "pass through exit code", args: []string{"exit", "99"}, expect: 99},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			code := RunCmd(tc.args, dirEnv)
			require.Equal(t, tc.expect, code)
		})
	}
}

func TestRunCmdEnv(t *testing.T) {
	dirEnv, err := ReadDir(testDataPath)
	require.NoError(t, err, "actual err - %v", err)
	tests := map[string]struct {
		value      string
		needRemove bool
	}{
		"BAR":   {value: "bar"},
		"EMPTY": {value: ""},
		"FOO":   {value: "   foo\nwith new line"},
		"HELLO": {value: `"hello"`},
		"UNSET": {needRemove: true},
	}

	for name, tc := range tests {
		name := name
		tc := tc
		t.Run(name, func(t *testing.T) {
			for key := range tests {
				t.Setenv(key, faker.Name())
			}
			fake := fakeio.Stdout()
			defer fake.Restore()
			code := RunCmd([]string{"echo", "-n", fmt.Sprintf(`"$%s"`, name)}, dirEnv)
			require.Equal(t, 0, code)
			out, err := fake.String()
			require.Equal(t, tc.value, out)
			require.NoError(t, err, "actual err - %v", err)
		})
	}
}

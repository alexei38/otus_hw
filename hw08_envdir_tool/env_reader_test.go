package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/require"
)

const testDataPath = "./testdata/env"

func TestCheckName(t *testing.T) {
	tests := []struct {
		name   string
		expect bool
	}{
		{name: "abcd", expect: true},
		{name: faker.Name(), expect: true},
		{name: "ab=cd", expect: false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equalf(t, checkName(tc.name), tc.expect, fmt.Sprintf("name %s not expected %v", tc.name, tc.expect))
		})
	}
}

func TestModifyLine(t *testing.T) {
	bufNewLine := bytes.Buffer{}
	bufNewLine.WriteString("AB")
	bufNewLine.WriteRune(0x00)
	bufNewLine.WriteString("CD")
	tests := []struct {
		name   []byte
		expect string
	}{
		{name: []byte("ABCD"), expect: "ABCD"},
		{name: []byte(" 12345"), expect: " 12345"},
		{name: []byte("ABCD\t"), expect: "ABCD"},
		{name: []byte("ABCD "), expect: "ABCD"},
		{name: bufNewLine.Bytes(), expect: "AB\nCD"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(string(tc.name), func(t *testing.T) {
			require.Equalf(t, tc.expect, modifyLine(tc.name), fmt.Sprintf("name %s not expected %v", tc.name, tc.expect))
		})
	}
}

func TestReadDir(t *testing.T) {
	testData := map[string]struct {
		value      string
		needRemove bool
	}{
		"BAR":   {value: "bar"},
		"EMPTY": {value: ""},
		"FOO":   {value: "   foo\nwith new line"},
		"HELLO": {value: `"hello"`},
		"UNSET": {needRemove: true},
	}

	t.Run("Check testdata without env", func(t *testing.T) {
		dirEnv, err := ReadDir(testDataPath)
		require.NoError(t, err, "actual err - %v", err)

		readDir, err := os.ReadDir(testDataPath)
		require.NoError(t, err, "actual err - %v", err)
		require.Len(t, dirEnv, len(readDir))

		for name, env := range dirEnv {
			v, ok := testData[name]
			require.True(t, ok)
			env.Value = v.value
			env.NeedRemove = v.needRemove
		}
	})

	t.Run("Check testdata with env", func(t *testing.T) {
		t.Setenv("BAR", faker.Name())
		t.Setenv("EMPTY", faker.Name())
		t.Setenv("FOO", faker.Name())
		t.Setenv("HELLO", faker.Name())
		t.Setenv("UNSET", faker.Name())
		customEnv := faker.Name()
		t.Setenv("T_CUSTOM_ENV", customEnv)

		require.Equal(t, os.Getenv("T_CUSTOM_ENV"), customEnv)

		dirEnv, err := ReadDir(testDataPath)
		require.NoError(t, err, "actual err - %v", err)

		readDir, err := os.ReadDir(testDataPath)
		require.NoError(t, err, "actual err - %v", err)
		require.Len(t, dirEnv, len(readDir))

		for name, env := range dirEnv {
			v, ok := testData[name]
			require.True(t, ok)
			env.Value = v.value
			env.NeedRemove = v.needRemove
		}
	})
}

package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetUsers(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()
	require.Equal(b, 1, len(r.File))

	for i := 0; i < b.N; i++ {
		data, _ := r.File[0].Open()
		getUsers(data)
	}
}

func BenchmarkCountDomains(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()
	require.Equal(b, 1, len(r.File))
	data, err := r.File[0].Open()
	require.NoError(b, err)

	users, err := getUsers(data)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		countDomains(users, "biz")
	}
}

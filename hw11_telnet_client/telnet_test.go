package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/rhysd/go-fakeio"
	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("check port missing", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer

		client := NewTelnetClient("127.0.0.1", time.Second*10, in, out)
		err := client.Connect()
		require.Equal(t, "dial tcp: address 127.0.0.1: missing port in address", err.Error())
	})

	t.Run("failed connection", func(t *testing.T) {
		var in io.ReadCloser
		var out io.Writer

		client := NewTelnetClient("127.0.0.1:65432", time.Second*10, in, out)
		err := client.Connect()
		require.Equal(t, "dial tcp 127.0.0.1:65432: connect: connection refused", err.Error())
	})

	t.Run("stderr messages", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			fake := fakeio.Stderr()
			defer fake.Restore()

			client := NewTelnetClient(l.Addr().String(), time.Second*5, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			stderr, err := fake.String()
			require.NoError(t, err)
			require.Contains(t, stderr, fmt.Sprintf("...Connected to %s", l.Addr().String()))
			require.NotContains(t, stderr, "...EOF")
			require.NotContains(t, stderr, "...Connection was closed by peer")

			fake = fakeio.Stderr()

			in.WriteString("Hello\n")
			err = client.Send()
			require.NoError(t, err)
			err = client.Receive()
			require.NoError(t, err)

			stderr, err = fake.String()
			require.NoError(t, err)
			require.Contains(t, strings.Split(stderr, "\n"), "...EOF")
			require.Contains(t, strings.Split(stderr, "\n"), "...Connection was closed by peer")
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() {
				require.NoError(t, conn.Close())
			}()
		}()

		wg.Wait()
	})
}

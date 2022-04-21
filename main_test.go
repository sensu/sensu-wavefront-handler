package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

func TestExecuteHandler(t *testing.T) {
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check = nil
	event.Metrics = corev2.FixtureMetrics()
	event.Metrics.Points[0].Timestamp = 1580922166749062000

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	errC := make(chan error, 1)
	go func() {
		defer close(errC)
		conn, err := ln.Accept()
		if err != nil {
			errC <- err
			return
		}
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, conn); err != nil {
			errC <- err
			return
		}
		want := `"answer" 42 1580922166 source="entity1" "foo"="bar"`
		if got := strings.TrimSpace(buf.String()); got != want {
			errC <- errors.New(cmp.Diff(got, want))
		}
	}()

	host, portS, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	handlerConfig.Host = host
	port, err := strconv.Atoi(portS)
	if err != nil {
		t.Fatal(err)
	}
	handlerConfig.MetricsPort = port
	if err := executeHandler(event); err != nil {
		t.Fatal(err)
	}
	if err := <-errC; err != nil {
		t.Error(err)
	}
}

func TestSecTimestamp(t *testing.T) {
	tests := []struct {
		In   int64
		Want int64
	}{
		{
			In:   1580922166749062000,
			Want: 1580922166,
		},
		{
			In:   1580922166749062,
			Want: 1580922166,
		},
		{
			In:   1580922166749,
			Want: 1580922166,
		},
		{
			In:   1580922166,
			Want: 1580922166,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%d => %d", test.In, test.Want), func(t *testing.T) {
			if got, want := secTimestamp(test.In), test.Want; got != want {
				t.Error(cmp.Diff(got, want))
			}
		})
	}
}

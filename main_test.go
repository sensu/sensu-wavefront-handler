package main

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteHandler(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check = nil
	event.Metrics = corev2.FixtureMetrics()
	event.Metrics.Points[0].Timestamp = 1580922166749062000

	var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gr, err := gzip.NewReader(r.Body)
		assert.NoError(err)
		body, err := ioutil.ReadAll(gr)
		assert.NoError(err)
		expectedBody := `"answer" 42 1580922166 source="entity1" "foo"="bar"`
                assert.Equal(expectedBody, strings.Trim(string(body), "\n"))
		w.WriteHeader(http.StatusOK)
	}))

	url, err := url.ParseRequestURI(test.URL)
	assert.NoError(err)
	handlerConfig.Host = url.Hostname()
	port, err := strconv.Atoi(url.Port())
	require.NoError(t, err)
	handlerConfig.MetricsPort = port
	assert.NoError(executeHandler(event))
}

func TestSecTimestamp(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check = nil
	event.Metrics = corev2.FixtureMetrics()
	event.Metrics.Points[0].Timestamp = 1580922166749062000
        ts := secTimestamp(event.Metrics.Points[0].Timestamp )
        assert.Equal(int64(1580922166),ts) 
	event.Metrics.Points[0].Timestamp = 1580922166749062
        ts = secTimestamp(event.Metrics.Points[0].Timestamp )
        assert.Equal(int64(1580922166),ts) 
	event.Metrics.Points[0].Timestamp = 1580922166749
        ts = secTimestamp(event.Metrics.Points[0].Timestamp )
        assert.Equal(int64(1580922166),ts) 
	event.Metrics.Points[0].Timestamp = 1580922166
        ts = secTimestamp(event.Metrics.Points[0].Timestamp )
        assert.Equal(int64(1580922166),ts) 
}

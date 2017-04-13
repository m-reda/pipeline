package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestAllPipelinesHandler(t *testing.T) {
	var pipelines []Pipeline
	err := json.Unmarshal(call(t, "GET", "/api/pipelines"), &pipelines)

	require.NoError(t, err)

	if len(pipelines) < 1 {
		t.Fatal("pipelines length should be more than one")
	}
}

func TestSinglePipelineHandler(t *testing.T) {
	var pipeline Pipeline
	err := json.Unmarshal(call(t, "GET", "/api/pipelines/"+mainPipeline.ID), &pipeline)

	require.NoError(t, err)
}

func TestSinglePipelineHandlerNotExist(t *testing.T) {
	c := call(t, "GET", "/api/pipelines/-")
	require.Equal(t, []byte(`{"success": false, "message": "Pipeline not exist"}`), c)
}

func TestStorePipelineHandlerEmptyData(t *testing.T) {
	c := call(t, "POST", "/api/pipelines")

	require.Equal(t, []byte(`{"success": false, "message": "pipeline field empty"}`), c)
}

func TestStorePipelineHandlerInvalidJSON(t *testing.T) {
	c := call(t, "POST", "/api/pipelines", strings.NewReader("pipeline=abc"))

	require.Equal(t, []byte(`{"success": false, "message": "invalid character 'a' looking for beginning of value"}`), c)
}

func TestStorePipelineHandlerAddNew(t *testing.T) {
	form := url.Values{}
	form.Add("pipeline", "{}")
	formReader := strings.NewReader(form.Encode())

	var resp struct {
		Success bool
		ID      string
	}
	err := json.Unmarshal(call(t, "POST", "/api/pipelines", formReader), &resp)

	require.NoError(t, err)
	require.NotEmpty(t, resp.ID)

	err = pipelineDelete(resp.ID)
	require.NoError(t, err)
}

func TestDeletePipelineHandler(t *testing.T) {
	// make new pipeline
	id, err := pipelineNew()

	require.NoError(t, err)
	require.Equal(t, []byte(`{"success": true}`), call(t, "DELETE", "/api/pipelines/"+id))
}

func TestBuildPipelineHandler(t *testing.T) {

	c := call(t, "GET", "/api/pipelines/"+mainPipeline.ID+"/build")
	require.Equal(t, []byte(`{"success": true}`), c)

	// wait until the build finish
	time.Sleep(30 * time.Millisecond)
}

func TestBuildPipelineHandlerWebsocket(t *testing.T) {
	u := url.URL{
		Scheme: "ws",
		Host:   strings.Split(server.URL, "//")[1],
		Path:   "/api/pipelines/" + mainPipeline.ID + "/build",
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.NoError(t, err)
	defer c.Close()

	done := make(chan struct{})
	buildLogs := []BuildLog{}

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()

			if string(message) == "done" || string(message) == "fail" || err != nil {
				return // close the connection
			}

			var l BuildLog
			err = json.Unmarshal(message, &l)
			assert.NoError(t, err)

			buildLogs = append(buildLogs, l)
		}
	}()

	<-done
	require.Equal(t, 2, len(buildLogs))
}

func TestBuildDetailsHandler(t *testing.T) {

	var build Build

	// correct request
	err := json.Unmarshal(call(t, "GET", "/api/pipelines/"+mainPipeline.ID+"/build/1"), &build)
	require.NoError(t, err)
	require.Equal(t, 2, len(build.Logs))

	// bad request
	err = json.Unmarshal(call(t, "GET", "/api/pipelines/-/build/-"), nil)
	require.Error(t, err)
}

func TestDeleteBuildHandler(t *testing.T) {

	// correct request
	c := call(t, "DELETE", "/api/pipelines/"+mainPipeline.ID+"/build/1")
	require.Equal(t, []byte(`{"success": true}`), c)

	// bad request
	c = call(t, "DELETE", "/api/pipelines/-/build/-")
	require.Equal(t, []byte(`{"success": false, "message": "open ./data/pipelines/-/pipeline.json: no such file or directory"}`), c)
}

func TestAllUnitsHandler(t *testing.T) {
	var units []Unit
	err := json.Unmarshal(call(t, "GET", "/api/units"), &units)

	require.NoError(t, err)
	if len(units) < 1 {
		t.Fatal("units length should be more than one")
	}
}

func TestHomeHandler(t *testing.T) {
	call(t, "GET", "/")
}

func call(t *testing.T, method, url string, bodyReader ...io.Reader) []byte {

	// form data
	var b io.Reader
	if len(bodyReader) > 0 {
		b = bodyReader[0]
	}

	// make the request
	request, err := http.NewRequest(method, server.URL+url, b)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(request)
	require.NoError(t, err)

	// check the status
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// read the body
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return body
}

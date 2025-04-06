package goargus

import (
	"os"
	"testing"
	"time"

	"github.com/epikur-io/go-argus/pkg/reader/filereader"
	"github.com/epikur-io/go-argus/pkg/watcher/filewatcher"
	"github.com/rs/zerolog"
	assert "github.com/stretchr/testify/assert"
)

type testSample struct {
	Key    string `json:"key" yaml:"key" toml:"key" ini:"key"`
	Num    int    `json:"num" yaml:"num" toml:"num" ini:"num"`
	Nested struct {
		Key string `json:"key" yaml:"key" toml:"key" ini:"key"`
	} `json:"nested" yaml:"nested" toml:"nested" ini:"nested"`
}

func TestDecoders(t *testing.T) {
	expected := testSample{
		Key: "value",
		Num: 13,
	}
	expected.Nested.Key = "value"

	testCases := []struct {
		name    string
		decoder option
		file    string
	}{
		{
			name:    "YAML decoder sould succeed",
			decoder: WithYamlDecoder(),
			file:    "sample.yaml",
		},
		{
			name:    "TOML decoder sould succeed",
			decoder: WithTomlDecoder(),
			file:    "sample.toml",
		},
		{
			name:    "JSON decoder sould succeed",
			decoder: WithJsonDecoder(),
			file:    "sample.json",
		},
		{
			name:    "INI decoder sould succeed",
			decoder: WithIniDecoder(),
			file:    "sample.ini",
		},
	}

	_ = os.Chdir("./test-data")
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			opts := []option{
				WithReader(filereader.New(testCase.file)),
				testCase.decoder,
			}
			argus, err := NewArgus[testSample](opts...)
			if err != nil {
				t.Error(err)
			}
			val := argus.GetValue()
			assert.Equal(t, val, expected)
		})
	}

}

func TestInvalidJsonFormat(t *testing.T) {
	_ = os.Chdir("./test-data")
	opts := []option{
		WithReader(filereader.New("sample_invalid.json")),
		WithJsonDecoder(),
	}
	argus, err := NewArgus[testSample](opts...)
	assert.Empty(t, argus)
	assert.NotEmpty(t, err)
}

func TestFileWatcher(t *testing.T) {
	// !TODO
	_ = os.Chdir("./test-data")
	testFile := "sample_watcher01.yaml"
	contents := `x: 1`
	fh, err := os.Create(testFile)
	assert.Empty(t, err)
	fh.WriteString(contents)
	assert.Empty(t, fh.Close())
	defer os.Remove(testFile)
	stop := make(chan struct{})
	opts := []option{
		WithReader(filereader.New(testFile)),
		WithWatcher(filewatcher.New(testFile)),
		WithYamlDecoder(),
		WithCallback(func(l *zerolog.Logger) {
			close(stop)
		}),
	}
	argus, err := NewArgus[testSample](opts...)
	assert.NotEmpty(t, argus)
	assert.Empty(t, err)
	assert.Empty(t, argus.StartWatcher())
	go func() {
		time.Sleep(time.Second / 2)
		assert.Empty(t, os.WriteFile(testFile, []byte(`x: 2`), os.ModePerm))
	}()
	<-stop
}

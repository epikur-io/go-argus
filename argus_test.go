package goargus

import (
	"os"
	"testing"

	"github.com/epikur-io/go-argus/pkg/reader/filereader"
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

func TestFileWatcher(t *testing.T) {
	// !TODO
}

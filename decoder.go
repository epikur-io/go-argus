package goargus

import (
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/ini.v1"
)

type iniDecoderWrapper struct {
	reader io.Reader
}

func (d *iniDecoderWrapper) Decode(val any) error {
	iniFile, err := ini.Load(d.reader)
	if err != nil {
		return err
	}
	return iniFile.MapTo(val)
}

type tomlDecoderWrapper struct {
	decoder *toml.Decoder
}

func (d *tomlDecoderWrapper) Decode(val any) error {
	_, err := d.decoder.Decode(val)
	return err
}

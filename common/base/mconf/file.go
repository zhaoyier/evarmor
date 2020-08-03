package mconf

import (
	"io/ioutil"

	"git.ezbuy.me/ezbuy/base/misc/errors"
	"github.com/BurntSushi/toml"
)

func ReadFile(fp string, obj interface{}) error {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return errors.Trace(err)
	}
	return Unmarshal(data, obj)
}

func Unmarshal(data []byte, obj interface{}) error {
	// data, err := Render(data)
	// if err != nil {
	// 	return errors.Trace(err)
	// }
	if err := toml.Unmarshal(data, obj); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// func Render(data []byte) ([]byte, error) {
// 	return RenderWithService(data, "")
// }

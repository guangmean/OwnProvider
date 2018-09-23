package jwt

import (
	"github.com/BurntSushi/toml"
)

type Toml struct {
	Header   JwtHeader
	Playload JwtPlayload
	Apple    JwtApple
}

func (t *Toml) getConfig(path string) (bool, error) {

	if _, err := toml.DecodeFile(path, t); nil != err {

		return false, err

	} else {

		return true, err

	}

}

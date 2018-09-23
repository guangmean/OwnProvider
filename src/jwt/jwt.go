package jwt

import (
	"time"
)

/*
 *	JWT Data Struct
 *	xxxxxx.yyyyyyyy.zzzzzzzzz
 *	Header.Playload.Signature
 *	@see:
 *		https://developer.apple.com/documentation/applemusicapi/getting_keys_and_creating_tokens
 *		https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns
 */

func Token(file string) (string, error) {

	// Get params from toml configure file
	var tomlConfig Toml
	ok, err := tomlConfig.getConfig(file)
	if false == ok {
		return "", err
	}

	// Step 1: Build JWT Header
	header := tomlConfig.Header
	headerBase64Content, err := header.Base64Content()
	if nil != err {
		return "", err
	}

	// Step 2: Build JWT Playload
	playload := tomlConfig.Playload
	playload.Iat = time.Now().Unix()
	playloadBase64Content, err := playload.Base64Content()
	if nil != err {
		return "", err
	}

	// Step 3: Signature Header & Playload
	signatureBase64Content, err := Signature(headerBase64Content, playloadBase64Content, tomlConfig.Apple.PrivateKey)
	if nil != err {
		return "", err
	}

	token := headerBase64Content + "." + playloadBase64Content + "." + signatureBase64Content

	return token, nil

}

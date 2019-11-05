package jwt

import (
	"OwnProvider/signature"
)

/*
 *	JWT Data Struct
 *	xxxxxx.yyyyyyyy.zzzzzzzzz
 *	Header.Payload.Signature
 *	@see:
 *		https://developer.apple.com/documentation/applemusicapi/getting_keys_and_creating_tokens
 *		https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/establishing_a_token-based_connection_to_apns
 */

func Token(header Header, payload Payload) (string, error) {

	// Step 1: Build JWT Header
	header64, err := header.Base64Content()
	if nil != err {
		return "", err
	}

	// Step 2: Build JWT Payload
	payload64, err := payload.Base64Content()
	if nil != err {
		return "", err
	}

	// Step 3: Signature Header & Payload
	sign64, err := signature.Sign(header64, payload64)
	if nil != err {
		return "", err
	}

	token := header64 + "." + payload64 + "." + sign64

	return token, nil

}

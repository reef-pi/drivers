package tplink

func autokeyEncrypt(cmd []byte) []byte {
	n := len(cmd)
	key := byte(0xAB)
	payload := make([]byte, n)
	for i := 0; i < n; i++ {
		payload[i] = cmd[i] ^ key
		key = payload[i]
	}
	return payload
}

func autokeyDecrypt(resp []byte) []byte {
	n := len(resp)
	key := byte(0xAB)
	payload := make([]byte, n)
	for i := 0; i < n; i++ {
		payload[i] = resp[i] ^ key
		key = resp[i]
	}
	return payload
}

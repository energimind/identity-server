package sessioncookie

// padSecret pads the secret to a minimum length.
// This is necessary for the AES encryption algorithm.
func padSecret(secret string) string {
	const minSecretLength = 32

	if len(secret) >= minSecretLength {
		return secret
	}

	return secret + string(make([]byte, minSecretLength-len(secret)))
}

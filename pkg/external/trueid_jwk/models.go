package trueid_jwk

type JWK struct {
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Use string `json:"use"`
}

type JWKSet struct {
	Keys []JWK `json:"keys"`
}

type Test struct {
}

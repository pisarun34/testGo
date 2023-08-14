package trueid_jwk

import (
	"encoding/json"
	"net/http"
	"os"
)

func FetchJWKSFromURL() (JWKSet, error) {
	resp, err := http.Get(os.Getenv("TRUEID_JWKS_URL"))
	if err != nil {
		return JWKSet{}, err
	}
	defer resp.Body.Close()

	var jwkSet JWKSet
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&jwkSet)
	if err != nil {
		return JWKSet{}, err
	}

	return jwkSet, nil
}

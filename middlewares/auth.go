package middlewares

import (
	"TESTGO/pkg/api/errors"
	"TESTGO/pkg/external/trueid_jwk"
	"crypto/rsa"
	"encoding/base64"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type CustomClaims struct {
	jwt.Claims
	SUB string `json:"sub"`
}

func ConvertJWKToRSAPublicKey(n, e string) (*rsa.PublicKey, error) {
	decodedN, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, err
	}

	decodedE, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, err
	}

	pubKey := &rsa.PublicKey{
		N: new(big.Int).SetBytes(decodedN),
		E: int(new(big.Int).SetBytes(decodedE).Int64()),
	}

	return pubKey, nil
}

// AuthTrueID is a function that check token from TrueID and extract ssoid from token to context
func AuthTrueID() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		jwks, err := trueid_jwk.FetchJWKSFromURL() // ดึง JWK จาก URL
		if err != nil {
			// if failed to fetch JWKs return error
			c.JSON(http.StatusUnauthorized, "Failed to fetch JWKs")
			c.Abort()
			return
		}

		// ตรวจสอบ token ด้วย JWK
		tk, err := jwt.ParseSigned(token)
		//... ตรวจสอบการถอดแกะข้อมูลจาก token และตรวจสอบกับ JWKs

		if err != nil {
			// if token is not valid return error
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		claimsValidated := false
		for _, key := range jwks.Keys {
			rsaPubKey, err := ConvertJWKToRSAPublicKey(key.N, key.E)
			if err != nil {
				continue // Skip to the next key if there's an error
			}

			jwk := jose.JSONWebKey{
				Key:       rsaPubKey,
				KeyID:     key.Kid,
				Algorithm: "RS256",
				Use:       key.Use,
			}

			//var claims jwt.Claims
			var customClaims CustomClaims
			err = tk.Claims(jwk, &customClaims)
			if err == nil {
				claimsValidated = true
				// set ssoid to context
				c.Set("ssoid", customClaims.SUB)
				break
			}
		}
		// if token is not valid return error
		if !claimsValidated {
			c.Abort()
			return
		}

		c.Next()
	}
}

/*
	func AuthSeekster(client external.SeeksterAPI, redis database.RedisClientInterface, db *gorm.DB) gin.HandlerFunc {
		return func(c *gin.Context) {
			if ssoidValue, exists := c.Get("ssoid"); exists {
				ssoid, ok := ssoidValue.(string)
				if !ok {
					c.Error(errors.ErrValidationInputSSOID)
					c.Abort()
					return
				}
				_, err := redis.GetSeeksterToken(context.Background(), ssoid)
				if err != nil {
					resp, err := handlers.SignInSeeksterAuto(handlers.NewSeeksterClient(), c, redis, db)
					if err != nil {
						if resp != nil {
							resultMap, err := handlers.ResponseBodyToStruct(resp)
							if err != nil {
								c.Error(errors.ErrParseJSON)
								c.Abort()
								return
							}
							c.Error(errors.NewExternalAPIError(resp.StatusCode(), resp.Status(), "", "", resultMap))
							c.Abort()
							return
						} else {
							c.Error(errors.ErrInternalServer)
							c.Abort()
							return
						}

					}
				}
			}

			c.Next()
		}
	}
*/
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" || !strings.HasPrefix(token, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// ตรวจสอบ token ที่ได้รับ
		// ...

		c.Next()
	}
}

func CheckAndExtractSSOID(c *gin.Context) (string, error) {
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			return "", errors.ErrValidationInputSSOID
		}
		return ssoid, nil
	}
	return "", errors.ErrExtractJWTTrueID
}

func CheckAndExtractSeeksterToken(c *gin.Context) (string, error) {
	if ssoidValue, exists := c.Get("seeksterToken"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			return "", errors.ErrValidationInputSSOID
		}
		return ssoid, nil
	}
	return "", errors.ErrExtractJWTTrueID
}

package handlers

import (
	"TESTGO/config"
	"TESTGO/middlewares"
	"TESTGO/pkg/api/errors"
	"TESTGO/pkg/api/requests"
	"TESTGO/pkg/database"
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/database/mysql"
	"fmt"
	"time"

	"TESTGO/pkg/external"
	"TESTGO/pkg/external/seekster"

	"context"
	"encoding/json"
	stderrors "errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type MyError struct {
	Message string
	Code    int
}

func (e *MyError) Error() string {
	return e.Message
}

// NewSeeksterClient is a function that return SeeksterAPI interface
func NewSeeksterClient() external.SeeksterAPI {
	return seekster.NewClient()
}

// ResponseBodyToStruct is a function that convert response body to map[string]interface{}
func ResponseBodyToStruct(resp *resty.Response) (interface{}, error) {
	// response body to map[string]interface{} for map erro response body to struct
	var resultMap map[string]interface{}
	bodyBytes := resp.String()
	// if response body is empty read body from resp.RawResponse.Body (Mock Seekster API)
	if resp.String() == "" {
		bodyBytes2, err3 := io.ReadAll(resp.RawResponse.Body)
		if err3 != nil {
			return nil, err3
		}
		bodyBytes = string(bodyBytes2)
	}
	// convert response body to map[string]interface{}
	err := json.Unmarshal([]byte(bodyBytes), &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap, nil

}

func bindAndValidateInput(c *gin.Context, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		return errors.ErrBadRequest
	}
	// Validate input
	if err := config.Validate.Struct(obj); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			errorMessages := handleValidationErrors(err, &config.Translator)
			return errors.NewAppError(http.StatusBadRequest, errors.ErrValidationInput.Code, strings.Join(errorMessages, " , "))
		} else {
			return errors.ErrValidationInput
		}
	}
	return nil
}

func handleAPIError(resp *resty.Response, err error) error {
	if err != nil {
		// convert response body to map[string]interface{}
		resultMap, err := ResponseBodyToStruct(resp)
		if err != nil {
			return errors.ErrParseJSON
		}
		return errors.NewExternalAPIError(resp.StatusCode(), resp.Status(), "", "", resultMap)
	}
	return nil
}

func signUpAPICall(client external.SeeksterAPI, user *models.User) (*seekster.SignUpResponse, error) {
	signUpUser, resp, err := client.SignUp(*user)
	if errHandle := handleAPIError(resp, err); errHandle != nil {
		return nil, errHandle
	}
	return signUpUser, nil
}

func signInAPICall(client external.SeeksterAPI, user *models.User) (*seekster.SignResponse, error) {
	signInUser, resp, err := client.SignInByPhone(*user)
	if errHandle := handleAPIError(resp, err); errHandle != nil {
		return nil, errHandle
	}
	return signInUser, err
}

func AuthSeekster(client external.SeeksterAPI, redis database.RedisClientInterface, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ssoid, err := middlewares.CheckAndExtractSSOID(c)
		if err != nil {
			c.Error(err)
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		_, err = redis.GetSeeksterToken(ctx, ssoid)
		if err != nil {
			// ทำการ handle error
			c.Error(errors.ErrRedis)
			return
		}
		seeksterUser, err := mysql.GetSeeksterUserBySSOID(ssoid, db)
		if err != nil {
			c.Error(err)
			return
		}
		signInUser, err := signInAPICall(client, seeksterUser)
		if err != nil {
			c.Error(err)
			return
		}
		ctx2, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		err = redis.SetSeeksterToken(ctx2, ssoid, signInUser.AccessToken)
		if err != nil {
			// ทำการ handle error
			c.Error(errors.ErrRedis)
			return
		}
		c.Next()
	}
}

// SeeksterSignin is a function that call Seekster SignIn api
func SeeksterSignin(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	// check ssoid is in context
	ssoid, err := middlewares.CheckAndExtractSSOID(c)
	if err != nil {
		fmt.Println("ssoid is not in context")
		c.Error(err)
		return
	}

	// check seekster token is in redis
	_, err = redis.GetSeeksterToken(context.Background(), ssoid)
	// if seekster token is not in redis call seekster api and set seekster token to redis if token is in redis return success
	if err != nil {
		// if error is redis: nil call seekster api and set seekster token to redis
		if err.Error() == "redis: nil" {
			seeksterUser, err := mysql.GetSeeksterUserBySSOID(ssoid, db)
			if err != nil {
				c.Error(err)
				return
			}
			// call seekster SignIn api
			signInUser, err := signInAPICall(client, seeksterUser)
			if err != nil {
				c.Error(err)
				return
			}

			// set seekster token to redis
			ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
			defer cancel()
			err = redis.SetSeeksterToken(ctx, ssoid, signInUser.AccessToken)
			if err != nil {
				// ทำการ handle error
				c.Error(errors.ErrRedis)
				return
			}
			//c.JSON(resp.StatusCode(), user)
			c.JSON(http.StatusOK, gin.H{"code": 10001, "message": "Success"})
			return
		} else {
			// Redis error
			c.Error(errors.ErrRedis)
			return
		}
	} else {
		// seekster token is in redis can call Seekster API
		c.JSON(http.StatusOK, gin.H{"code": 10001, "message": "Success"})
		return
	}

}

// SeeksterSignup is a function that call Seekster SignUp api
func SeeksterSignup(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	// Bind input
	var input requests.SignUpInput
	if err := bindAndValidateInput(c, &input); err != nil {
		c.Error(err)
		return
	}

	ssoid, err := middlewares.CheckAndExtractSSOID(c)
	if err != nil {
		c.Error(err)
		return
	}

	newUser, err := mysql.CreateAndInsertSeeksterUser(ssoid, &input, db)
	if err != nil {
		c.Error(err)
		return
	}
	// call seekster SignUp api
	signUpUser, err := signUpAPICall(client, newUser)
	if err != nil {
		c.Error(err)
		return
	}
	// return SignUpResponse
	// set seekster token to redis
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err = redis.SetSeeksterToken(ctx, ssoid, signUpUser.AccessToken)
	if err != nil {
		// ทำการ handle error
		c.Error(errors.ErrRedis)
		return
	}
	c.JSON(200, signUpUser)

}

func SignInSeeksterAuto(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) (*resty.Response, error) {
	// check ssoid is in context
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			return nil, stderrors.New("ssoid is not string")
		}
		// check seekster token is in redis
		_, err := redis.GetSeeksterToken(context.Background(), ssoid)
		// if seekster token is not in redis call seekster api and set seekster token to redis if token is in redis return success
		if err != nil {
			// if error is redis: nil call seekster api and set seekster token to redis
			if err.Error() == "redis: nil" {
				var seeksterUser models.User
				// get seekster user from db
				if err := db.Where("SSOID = ?", ssoid).
					Preload("SeeksterUser").
					First(&seeksterUser).Error; err != nil {
					return nil, err
				}
				// call seekster SignIn api
				user, resp, err := client.SignInByPhone(seeksterUser)
				if err != nil {
					return resp, err
				}
				// set seekster token to redis
				redis.SetSeeksterToken(context.Background(), ssoid, user.AccessToken)
				//c.JSON(resp.StatusCode(), user)
				return nil, nil
			} else {
				// Redis error
				return nil, err
			}
		} else {
			// seekster token is in redis can call Seekster API
			return nil, err
		}
	} else {
		// ssoid is not in context
		return nil, stderrors.New("ssoid is not in context")
	}

}
func InsertSeeksterUser(client external.SeeksterAPI, c *gin.Context, db *gorm.DB) {

}

func GetServiceInfo(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	c.JSON(http.StatusOK, gin.H{"message": "test"})
}

package handlers

import (
	"TESTGO/config"
	"TESTGO/pkg/api/errors"
	"TESTGO/pkg/api/requests"
	"TESTGO/pkg/database"
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/database/mysql"
	"TESTGO/pkg/database/redis"
	"TESTGO/pkg/entities"
	"TESTGO/pkg/external"
	"TESTGO/pkg/external/seekster"
	"TESTGO/pkg/utils"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"io"
	"time"

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

// SeeksterSignin is a function that call Seekster SignIn api
func SeeksterSignin(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	// check ssoid is in context
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			fmt.Println("ssoid is not string")
			c.JSON(errors.ErrValidationInputSSOID.Status, errors.ErrValidationInputSSOID)
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
					c.JSON(errors.ErrDatabase.Status, errors.ErrDatabase)
					return
				}
				// call seekster SignIn api
				user, resp, err := client.SignInByPhone(seeksterUser)
				if err != nil {
					// convert response body to map[string]interface{}
					resultMap, err := ResponseBodyToStruct(resp)
					if err != nil {
						c.JSON(errors.ErrParseJSON.Status, errors.ErrParseJSON)
						return
					}
					c.JSON(resp.StatusCode(), resultMap)
					return
				}
				// set seekster token to redis
				redis.SetSeeksterToken(context.Background(), ssoid, user.AccessToken)
				//c.JSON(resp.StatusCode(), user)
				c.JSON(200, gin.H{"code": 10001, "message": seeksterUser})
				return
			} else {
				// Redis error
				c.JSON(errors.ErrRedis.Status, errors.ErrRedis)
				return
			}
		} else {
			// seekster token is in redis can call Seekster API
			c.JSON(200, gin.H{"code": 10001, "message": "Success"})
			return
		}
	} else {
		// ssoid is not in context
		c.JSON(errors.ErrExtractJWTTrueID.Status, errors.ErrExtractJWTTrueID)
		return
	}

}

// SeeksterSignup is a function that call Seekster SignUp api
func SeeksterSignup(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	// Bind input
	var input requests.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(errors.ErrBadRequest.Status, errors.ErrBadRequest)
		return
	}
	// Validate input
	if err := config.Validate.Struct(&input); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			errorMessages := handleValidationErrors(err, &config.Translator)
			c.JSON(400, gin.H{"error": errorMessages})
			return
		} else {
			c.JSON(errors.ErrValidationInput.Status, errors.ErrValidationInput)
			return
		}
	}
	// check ssoid is in context
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			c.JSON(errors.ErrValidationInputSSOID.Status, errors.ErrValidationInputSSOID)
		}
		// create seekster user input
		user := entities.InsertSeeksterUserInput{
			SSOID:       ssoid,
			PhoneNumber: input.PhoneNumber,
			Password:    utils.GenerateRandomPassword(15),
			UUID:        utils.GenerateUUIDv4(),
		}
		// Validate input
		if err := config.Validate.Struct(&user); err != nil {
			c.JSON(errors.ErrValidationModel.Status, errors.ErrValidationModel)
			return
		}
		// Insert DB SeeksterUser and User if exist return error
		newUser, err := mysql.InsertSeeksterUser(db, user)
		if err != nil {
			// if error is not SeeksterUser already exists return error
			if err.Error() == "SeeksterUser already exists" {
				c.JSON(errors.ErrSeeksterUserExist.Status, errors.ErrSeeksterUserExist)
			} else {
				c.JSON(errors.ErrInternalServer.Status, errors.ErrInternalServer)
			}
			return
		}
		// call seekster SignUp api
		signUpUser, resp, err := client.SignUp(*newUser)
		if err != nil {
			// convert response body to map[string]interface{}
			resultMap, err := ResponseBodyToStruct(resp)
			if err != nil {
				c.JSON(errors.ErrParseJSON.Status, errors.ErrParseJSON)
				return
			}
			c.JSON(resp.StatusCode(), resultMap)
			return
		}
		// set seekster token to redis
		err = redis.SetSeeksterToken(context.Background(), ssoid, signUpUser.AccessToken)
		if err != nil {
			// ทำการ handle error
			c.JSON(errors.ErrRedis.Code, errors.ErrRedis)
			return
		}
		// return SignUpResponse
		c.JSON(200, signUpUser)
	} else {
		c.JSON(errors.ErrExtractJWTTrueID.Status, errors.ErrExtractJWTTrueID)
		return
	}

}

func SignInSeeksterAuto(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) (*resty.Response, error) {
	// check ssoid is in context
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			c.JSON(errors.ErrValidationInputSSOID.Status, errors.ErrValidationInputSSOID)
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
	ctx := context.Background()
	pong, err := redis.RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Failed to connect to Redis:", err)
	} else {
		fmt.Println("Connected to Redis. Response:", pong)
	}
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		fmt.Println(ssoid)
		fmt.Println(redis.SeeksterTokenNamespace)
		if !ok {
			c.JSON(500, gin.H{"error": "internal error"})
		}
		err := redis.RedisClient.Set(context.TODO(), fmt.Sprintf("%s:%s", redis.SeeksterTokenNamespace, ssoid), "signUpUser.AccessToken", 24*7*time.Hour).Err()
		if err != nil {
			// ทำการ handle error
			c.JSON(errors.ErrRedis.Code, errors.ErrRedis)
			return
		}
		c.JSON(200, "OK")
	}
}

func GetServiceInfo(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	c.JSON(200, gin.H{"message": "test"})
}

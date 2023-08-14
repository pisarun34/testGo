package handlers

import (
	"TESTGO/config"
	"TESTGO/models"
	"TESTGO/pkg/database"
	"TESTGO/pkg/database/mysql"
	"TESTGO/pkg/database/redis"
	"TESTGO/pkg/external"
	"TESTGO/pkg/external/seekster"
	"TESTGO/utils"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

func NewSeeksterClient() external.SeeksterAPI {
	return seekster.NewClient()
}

func ResponseBodyToStruct(resp *resty.Response) (interface{}, error) {
	var resultMap map[string]interface{}
	bodyBytes := resp.String()
	if resp.String() == "" {
		bodyBytes2, err3 := io.ReadAll(resp.RawResponse.Body)
		if err3 != nil {
			return nil, err3
		}
		bodyBytes = string(bodyBytes2)
	}
	err := json.Unmarshal([]byte(bodyBytes), &resultMap)
	if err != nil {
		return nil, err
	}
	return resultMap, nil

}

func SeeksterSignin(client external.SeeksterAPI, c *gin.Context, redis database.RedisClientInterface, db *gorm.DB) {
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			fmt.Println("ssoid is not string")
			c.JSON(models.ErrValidationInputSSOID.Status, models.ErrValidationInputSSOID)
		}

		_, err := redis.GetSeeksterToken(context.Background(), ssoid)
		if err != nil {
			if err.Error() == "redis: nil" {
				var seeksterUser mysql.User
				if err := db.Where("SSOID = ?", ssoid).
					Preload("SeeksterUser").
					First(&seeksterUser).Error; err != nil {
					c.JSON(models.ErrDatabase.Status, models.ErrDatabase)
					return
				}
				user, resp, err := client.SignInByPhone(seeksterUser)
				if err != nil {
					resultMap, err := ResponseBodyToStruct(resp)
					if err != nil {
						c.JSON(models.ErrParseJSON.Status, models.ErrParseJSON)
						return
					}
					c.JSON(resp.StatusCode(), resultMap)
					return
				}
				redis.SetSeeksterToken(context.Background(), ssoid, user.AccessToken)
				//c.JSON(resp.StatusCode(), user)
				c.JSON(200, gin.H{"code": 10001, "message": "Success"})
				return
			} else {
				c.JSON(models.ErrRedis.Status, models.ErrRedis)
				return
			}
		} else {
			c.JSON(200, gin.H{"code": 10001, "message": "Success"})
			return
		}
	} else {
		c.JSON(models.ErrExtractJWTTrueID.Status, models.ErrExtractJWTTrueID)
		return
	}

}

func SeeksterSignup(client external.SeeksterAPI, c *gin.Context, db *gorm.DB) {
	var input models.SignUpInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(models.ErrBadRequest.Status, models.ErrBadRequest)
		return
	}

	if err := config.Validate.Struct(&input); err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			errorMessages := handleValidationErrors(err, &config.Translator)
			c.JSON(400, gin.H{"error": errorMessages})
			return
		} else {
			c.JSON(models.ErrValidationInput.Status, models.ErrValidationInput)
			return
		}
	}
	if ssoidValue, exists := c.Get("ssoid"); exists {
		ssoid, ok := ssoidValue.(string)
		if !ok {
			c.JSON(models.ErrValidationInputSSOID.Status, models.ErrValidationInputSSOID)
		}
		user := models.InsertSeeksterUserInput{
			SSOID:       ssoid,
			PhoneNumber: input.PhoneNumber,
			Password:    utils.GenerateRandomPassword(15),
			UUID:        utils.GenerateUUIDv4(),
		}
		if err := config.Validate.Struct(&user); err != nil {
			c.JSON(models.ErrValidationModel.Status, models.ErrValidationModel)
			return
		}
		newUser, err := mysql.InsertSeeksterUser(db, user)
		if err != nil {
			if err.Error() == "SeeksterUser already exists" {
				c.JSON(models.ErrSeeksterUserExist.Status, models.ErrSeeksterUserExist)
			} else {
				c.JSON(models.ErrInternalServer.Status, models.ErrInternalServer)
			}
			return
		}
		signUpUser, resp, err := client.SignUp(*newUser)
		if err != nil {
			var resultMap map[string]interface{}
			err := json.Unmarshal([]byte(resp.String()), &resultMap)
			if err != nil {
				c.JSON(models.ErrParseJSON.Status, models.ErrParseJSON)
				return
			}
			c.JSON(resp.StatusCode(), resultMap)
			return
		}
		err = redis.RedisClient.Set(context.TODO(), fmt.Sprintf("%s:%s", redis.SeeksterTokenNamespace, ssoid), signUpUser.AccessToken, 24*7*time.Hour).Err() // ค่า "some_value" คือค่าที่คุณต้องการเก็บ
		if err != nil {
			// ทำการ handle error
			c.JSON(models.ErrRedis.Code, models.ErrRedis)
			return
		}
		c.JSON(200, signUpUser)
	} else {
		c.JSON(models.ErrExtractJWTTrueID.Status, models.ErrExtractJWTTrueID)
		return
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
			c.JSON(models.ErrRedis.Code, models.ErrRedis)
			return
		}
		c.JSON(200, "OK")
	}
}

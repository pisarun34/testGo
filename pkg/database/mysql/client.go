package mysql

import (
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/entities"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func InsertSeeksterUser(db *gorm.DB, input entities.InsertSeeksterUserInput) (*models.User, error) {
	var user models.User
	// Check if User already exists
	if err := db.Where("ssoid = ?", input.SSOID).First(&user).Error; err != nil {
		// If the user does not exist, create a new User and SeeksterUser
		if err == gorm.ErrRecordNotFound {
			newUser := models.User{
				SSOID: input.SSOID,
			}

			// Insert User to the database
			if err := db.Create(&newUser).Error; err != nil {
				return nil, err
			}

			// Create a new SeeksterUser with the given input
			seeksterUser := models.SeeksterUser{
				PhoneNumber: input.PhoneNumber,
				Password:    input.Password,
				UUID:        input.UUID,
				UserID:      newUser.ID, // Set foreign key to the newly created User's ID
			}

			// Insert SeeksterUser to the database
			if err := db.Create(&seeksterUser).Error; err != nil {
				return nil, err
			}
			// Set SeeksterUser to the newly created User
			newUser.SeeksterUser = seeksterUser
			return &newUser, nil
		} else {
			// Handle other errors that might have occurred during the query
			return nil, err
		}
	} else {
		// If the user exists, check if SeeksterUser already exists
		var seeksterUser models.SeeksterUser
		if err := db.Where("user_id = ?", user.ID).First(&seeksterUser).Error; err != nil {
			fmt.Println(err)
			// If SeeksterUser does not exist, create a new SeeksterUser
			if err == gorm.ErrRecordNotFound {
				newSeeksterUser := models.SeeksterUser{
					PhoneNumber: input.PhoneNumber,
					Password:    input.Password,
					UUID:        input.UUID,
					UserID:      user.ID,
				}

				// Insert SeeksterUser to the database
				if err := db.Create(&newSeeksterUser).Error; err != nil {
					return nil, err
				}
				// Set SeeksterUser to the existing User
				user.SeeksterUser = newSeeksterUser
				return &user, nil
			} else {
				// Handle other errors that might have occurred during the query
				return nil, err
			}
		} else {
			return nil, errors.New("SeeksterUser already exists")
		}
	}
}

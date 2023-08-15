package external

import (
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/external/seekster"
	"TESTGO/pkg/external/trueidpartner"

	"github.com/go-resty/resty/v2"
)

// SeeksterAPI สร้าง interface สำหรับเรียกใช้งาน TrueID API
type SeeksterAPI interface {
	SignInByPhone(seeksterUser models.User) (*seekster.SignResponse, *resty.Response, error)
	SignUp(seeksterUser models.User) (*seekster.SignUpResponse, *resty.Response, error)
}

// TrueIDSubscripberAPI สร้าง interface สำหรับเรียกใช้งาน TrueID API
type TrueIDSubscripberAPI interface {
	GetSubscribers(ssoid string) (*trueidpartner.Subscribers, *resty.Response, error)
}

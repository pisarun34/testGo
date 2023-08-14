package external

import (
	"TESTGO/pkg/database/mysql"
	"TESTGO/pkg/external/seekster"
	"TESTGO/pkg/external/trueidpartner"

	"github.com/go-resty/resty/v2"
)

// SeeksterAPI สร้าง interface สำหรับเรียกใช้งาน TrueID API
type SeeksterAPI interface {
	SignInByPhone(seeksterUser mysql.User) (*seekster.SignResponse, *resty.Response, error)
	SignUp(seeksterUser mysql.User) (*seekster.SignUpResponse, *resty.Response, error)
}

// TrueIDSubscripberAPI สร้าง interface สำหรับเรียกใช้งาน TrueID API
type TrueIDSubscripberAPI interface {
	GetSubscribers(ssoid string) (*trueidpartner.Subscribers, *resty.Response, error)
}

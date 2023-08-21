package external

import (
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/external/seekster"
	"TESTGO/pkg/external/trueidpartner"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// SeeksterAPI สร้าง interface สำหรับเรียกใช้งาน TrueID API
type SeeksterAPI interface {
	SignInByPhone(seeksterUser models.User) (*seekster.SignResponse, *resty.Response, error)
	SignUp(seeksterUser models.User) (*seekster.SignUpResponse, *resty.Response, error)
	GetServiceList(c *gin.Context) (*seekster.GetServiceListResponse, *resty.Response, error)
	GetServiceDetails(c *gin.Context) (*seekster.GetServiceDetailsResponse, *resty.Response, error)
	GetSlotsQuantity(c *gin.Context) (*seekster.GetSlotsQuantityResponse, *resty.Response, error)
	BookingServiceBySlot(c *gin.Context, RequestBody *seekster.BookingServiceBySlotRequest) (*seekster.BookingServiceBySlotResponse, *resty.Response, error)
	GetInquiryList(c *gin.Context) (*seekster.GetInquiryListResponse, *resty.Response, error)
}

// TrueIDSubscripberAPI สร้าง interface สำหรับเรียกใช้งาน TrueID API
type TrueIDSubscripberAPI interface {
	GetSubscribers(ssoid string) (*trueidpartner.SubscribersResponse, *resty.Response, error)
}

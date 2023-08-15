package external

import (
	"TESTGO/pkg/database/models"
	"TESTGO/pkg/external/seekster"

	"github.com/go-resty/resty/v2"
)

type SeeksterAPIMock struct {
	// คุณสามารถกำหนด field อื่น ๆ ที่ต้องการสำหรับ Mock Object ตามที่คุณต้องการในการทดสอบ
}

func (m *SeeksterAPIMock) SignInByPhone(seeksterUser models.User) (*seekster.SignResponse, *resty.Response, error) {
	// ในที่นี้คุณสามารถกำหนดการทำงานของ Mock Object สำหรับฟังก์ชันนี้ได้ตามเหตุการณ์ที่คุณต้องการในการทดสอบ
	return nil, nil, nil
}

func (m *SeeksterAPIMock) SignUp(seeksterUser models.User) (*seekster.SignUpResponse, *resty.Response, error) {
	// ในที่นี้คุณสามารถกำหนดการทำงานของ Mock Object สำหรับฟังก์ชันนี้ได้ตามเหตุการณ์ที่คุณต้องการในการทดสอบ
	return nil, nil, nil
}

package seekster

type SignInRequest struct {
	DevicesAttributes struct {
		Brand             string `json:"brand"`
		Carrier           string `json:"carrier"`
		ClientType        string `json:"client_type"`
		ClientVersion     string `json:"client_version"`
		Locale            string `json:"locale"`
		Model             string `json:"model"`
		OsVersion         string `json:"os_version"`
		UUID              string `json:"uuid"`
		RegistrationToken string `json:"registration_token"`
	} `json:"devices_attributes"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type SignUpRequest struct {
	DevicesAttributes struct {
		Brand             string `json:"brand"`
		Carrier           string `json:"carrier"`
		OwnerType         string `json:"owner_type"`
		ClientVersion     string `json:"client_version"`
		Locale            string `json:"locale"`
		Model             string `json:"model"`
		OsVersion         string `json:"os_version"`
		RegistrationToken string `json:"registration_token"`
		UUID              string `json:"uuid"`
	} `json:"devices_attributes"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	Password          string `json:"password"`
	PreferredLanguage string `json:"preferred_language"`
}

type BookingServiceBySlotRequest struct {
	AnswersAttributes []struct {
		QuestionID int    `json:"question_id"`
		ChoiceID   int    `json:"choice_id,omitempty"`
		Value      string `json:"value,omitempty"`
		ChoiceIds  []int  `json:"choice_ids,omitempty"`
	} `json:"answers_attributes"`
	OrdersAttributes []struct {
		PackageID int `json:"package_id"`
		Quantity  int `json:"quantity"`
	} `json:"orders_attributes"`
	ServiceID      int `json:"service_id"`
	JobsAttributes []struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	} `json:"jobs_attributes"`
	LocationAttributes struct {
		Name                  string `json:"name"`
		ContactName           string `json:"contact_name"`
		District              string `json:"district"`
		HouseNumber           string `json:"house_number"`
		Landmark              string `json:"landmark"`
		Latitude              string `json:"latitude"`
		Longitude             string `json:"longitude"`
		PhoneNumber           string `json:"phone_number"`
		ProjectName           string `json:"project_name"`
		Province              string `json:"province"`
		Street                string `json:"street"`
		SubDistrict           string `json:"sub_district"`
		Alley                 string `json:"alley"`
		Floor                 string `json:"floor"`
		ZipCode               string `json:"zip_code"`
		Country               string `json:"country"`
		AdditionalInformation string `json:"additional_information"`
	} `json:"location_attributes"`
}

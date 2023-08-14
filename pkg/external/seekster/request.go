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

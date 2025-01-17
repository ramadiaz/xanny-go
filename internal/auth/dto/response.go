package dto

type UserOutput struct {
	UUID            string `json:"uuid"`
	Username        string `json:"username"`
	BusinessName    string `json:"business_name"`
	BusinessOwner   string `json:"business_owner"`
	BusinessPhone   string `json:"business_phone"`
	BusinessAddress string `json:"business_address"`
	BusinessLogo    string `json:"business_logo"`
	Slogan          string `json:"slogan"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body,omitempty"`
}

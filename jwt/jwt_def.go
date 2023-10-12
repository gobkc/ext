package jwt

type UserClaims struct {
	Uid         int64  `json:"uid"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	CountryCode string `json:"country_code"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	State       int    `json:"state"`
	Type        int    `json:"type"`
	LastIP      string `json:"last_ip"`
	Expire      int64  `json:"expire"`
	Ext         any    `json:"ext"`
}

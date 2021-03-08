package bean

type UserBean struct {
	UserName     string `json:"user_name"`
	UserId       string `json:"user_id"`
	Account      string `json:"account"`
	HeadPortrait string `json:"head_portrait"`
	VipTime      string `json:"vip_time"`
	Pwd          string `json:"-"`
}

func NewUserBean(data map[string]string) *UserBean {
	return &UserBean{
		Account:      data["account"],
		UserName:     data["user_name"],
		UserId:       data["user_id"],
		VipTime:      data["vip_time"],
		HeadPortrait: data["head_portrait"],
		Pwd:          data["pwd"],
	}
}

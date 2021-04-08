package bean

type UserBean struct {
	UserName     string `json:"user_name"`
	UserId       string `json:"user_id"`
	Account      string `json:"account"`
	HeadPortrait string `json:"head_portrait"`
	VipTime      string `json:"vip_time"`
	Permissions  string `json:"permissions"`
	Pwd          string `json:"-"`
}

func NewUserBean(data map[string]string) *UserBean {
	return &UserBean{
		Account:      data["account"],
		UserName:     data["user_name"],
		UserId:       data["user_id"],
		VipTime:      data["vip_time"],
		HeadPortrait: data["head_portrait"],
		Permissions:  data["permissions"],
		Pwd:          data["pwd"],
	}
}
func (u *UserBean) Modify(new UserBean) {
	if len(new.UserName) != 0 {
		u.UserName = new.UserName
	}
	if len(new.UserName) != 0 {
		u.UserName = new.UserName
	}
	if len(new.HeadPortrait) != 0 {
		u.HeadPortrait = new.HeadPortrait
	}
	if len(new.HeadPortrait) != 0 {
		u.HeadPortrait = new.HeadPortrait
	}
	if len(new.Pwd) != 0 {
		u.Pwd = new.Pwd
	}
	if len(new.VipTime) != 0 {
		u.VipTime = new.VipTime
	}
}

type SecurityBean struct {
	Question1 string `json:"question1"`
	Answer1   string `json:"answer1"`
	Question2 string `json:"question2"`
	Answer2   string `json:"answer2"`
}

func NewSecurityBean(data map[string]string) *SecurityBean {
	return &SecurityBean{
		Question1: data["question1"],
		Answer1:   data["answer1"],
		Question2: data["question2"],
		Answer2:   data["answer2"],
	}
}

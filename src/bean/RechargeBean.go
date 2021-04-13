package bean

type RechargeBean struct {
	UserId               string `json:"user_id"`
	UserAccount          string `json:"user_account"`
	UserName             string `json:"user_name"`
	ExecutionUserId      string `json:"execution_user_id"`
	ExecutionUserName    string `json:"execution_user_name"`
	ExecutionUserAccount string `json:"execution_user_account"`
	Money                string `json:"money"`
	Day                  string `json:"day"`
	RechargeType         string `json:"recharge_type"`
	Created              string `json:"created"`
}

func NewRechargeBean(data map[string]string) *RechargeBean {
	return &RechargeBean{
		UserId:               data["user_id"],
		UserName:             data["user_name"],
		UserAccount:          data["user_account"],
		ExecutionUserId:      data["execution_user_id"],
		ExecutionUserName:    data["execution_user_name"],
		ExecutionUserAccount: data["execution_user_account"],
		Money:                data["money"],
		Day:                  data["day"],
		RechargeType:         data["recharge_type"],
		Created:              data["created"],
	}
}

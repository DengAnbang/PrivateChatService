package bean

type ChatGroup struct {
	GroupId       string `json:"group_id"`
	UserId        string `json:"user_id"`
	UserType      string `json:"user_type"`
	ChatPwd       string `json:"chat_pwd"`
	GroupName     string `json:"group_name"`
	GroupPortrait string `json:"group_portrait"`
}

func NewChatGroup(data map[string]string) *ChatGroup {
	return &ChatGroup{
		GroupId:       data["group_id"],
		UserId:        data["user_id"],
		UserType:      data["user_type"],
		ChatPwd:       data["chat_pwd"],
		GroupName:     data["group_name"],
		GroupPortrait: data["group_portrait"],
	}
}

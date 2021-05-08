package bean

type PriceBean struct {
	Id        string `json:"id"`
	Money     string `json:"money"`
	Day       string `json:"day"`
	GivingDay string `json:"giving_day"`
	PayImage  string `json:"pay_image"`
}

func NewPriceBean(data map[string]string) *PriceBean {
	return &PriceBean{
		Id:        data["id"],
		Money:     data["money"],
		Day:       data["day"],
		GivingDay: data["giving_day"],
		PayImage:  data["pay_image"],
	}
}

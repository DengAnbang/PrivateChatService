package bean

type PriceBean struct {
	Id        string `json:"id"`
	Money     string `json:"money"`
	Day       string `json:"day"`
	GivingDay string `json:"giving_day"`
}

func NewPriceBean(data map[string]string) *PriceBean {
	return &PriceBean{
		Id:        data["id"],
		Money:     data["money"],
		Day:       data["day"],
		GivingDay: data["giving_day"],
	}
}

package bean

type FileBean struct {
	Id       string `json:"-"`
	FileId   string `json:"file_id"`
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
}

func NewFileBean(data map[string]string) FileBean {
	return FileBean{
		Id:       data["id"],
		FileName: data["file_name"],
		FilePath: data["file_path"],
		FileId:   data["file_id"],
	}
}

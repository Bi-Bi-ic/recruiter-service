package factory

type Fileable struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type FileInfoFactoty struct{}

func (factory FileInfoFactoty) Create(fileName string, id string) Fileable {
	var file = Fileable{id, "/api/file/" + fileName}
	return file
}

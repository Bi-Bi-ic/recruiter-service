package repository

//LocationRepository is an interface can be implemented
type LocationRepository interface {
	GetProvinceList() (RepoResponse, Status)
	GetDistrictList(provinceID int) (RepoResponse, Status)
	GetWardList(cityID, districtID int) (RepoResponse, Status)
	FindAddress(string) (RepoResponse, Status)
}

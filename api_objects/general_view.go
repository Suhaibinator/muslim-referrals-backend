package api_objects

import "github.com/Suhaibinator/muslim-referrals-backend/database"

// GeneralView represents the fields that the general user will be able to see

type GeneralViewCompany struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func ConvertDbCompanyToGeneralViewCompany(dbCompany *database.Company) *GeneralViewCompany {
	return &GeneralViewCompany{
		Id:   dbCompany.Id,
		Name: dbCompany.Name,
	}
}

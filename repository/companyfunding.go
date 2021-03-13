package repository

import (
	"company-funding/dto"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type CompanyFunding struct {
	gorm.Model
	Name                 string
	Description          string
	FundingAmount        string
	FundingRound         string
	Investors            string
	LocationCity         string
	LocationMunicipality string
	Industry             string
	EmployeeCount        int
}

var (
	pattern = regexp.MustCompile(`(\d+)`)
)

func ConvertCompanyDto(dto *dto.FlCompanyDTO) *CompanyFunding {
	ecSlice := pattern.FindAllString(dto.EmployeeCount, 1)
	var ec string
	e := -1
	if len(ecSlice) > 0 {
		ec = ecSlice[0]
		e, _ = strconv.Atoi(ec)
	}

	cf := &CompanyFunding{
		Name:          dto.Name,
		Description:   dto.Description,
		Investors:     dto.Investors,
		Industry:      dto.Industry,
		EmployeeCount: e,
	}

	funding := strings.Split(dto.Funding, ",")
	// slice funding is length 1 if there is no comma
	if len(funding) > 1 {
		cf.FundingAmount = funding[0]
		cf.FundingRound = funding[1]
	} else {
		cf.FundingRound = funding[0]
	}

	location := strings.Split(dto.Location, ",")
	if len(location) > 0 {
		cf.LocationCity = location[0]
		cf.LocationMunicipality = location[1]
	} else {
		cf.LocationMunicipality = location[0]
	}

	return cf
}

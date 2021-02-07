package repository

import (
	"company-funding/dto"
	"fmt"
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
	ec := pattern.FindAllString(dto.EmployeeCount, 1)[0]
	e, err := strconv.Atoi(ec)
	if err != nil {
		fmt.Printf("Warn: could not convert ec for company %s\n", dto.Name)
		e = -1
	}

	funding := strings.Split(dto.Funding, ",")
	location := strings.Split(dto.Location, ",")

	cf := &CompanyFunding{
		Name:                 dto.Name,
		Description:          dto.Description,
		FundingAmount:        funding[0],
		FundingRound:         funding[1],
		Investors:            dto.Investors,
		LocationCity:         location[0],
		LocationMunicipality: location[1],
		Industry:             dto.Industry,
		EmployeeCount:        e,
	}

	return cf
}

package dto

import (
	"company-funding/repository"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FlCompanyDTO struct {
	Name          string
	Description   string
	Funding       string
	Investors     string
	Location      string
	Industry      string
	EmployeeCount string
}

var (
	pattern = regexp.MustCompile(`(\d+)`)
)

func (dto *FlCompanyDTO) PrintDto() {
	fmt.Printf("Company: %s\n", dto.Name)
	fmt.Printf("Description: %s\n", dto.Description)
	fmt.Printf("Funding: %s\n", dto.Funding)
	fmt.Printf("Investors: %s\n", dto.Investors)
	fmt.Printf("Location: %s\n", dto.Location)
	fmt.Printf("Industry: %s\n", dto.Industry)
	fmt.Printf("EmployeeCount: %s\n", dto.EmployeeCount)
}

func (dto *FlCompanyDTO) IsSaveAble() bool {
	if dto.Funding != "" || dto.Industry != "" || dto.Investors != "" || dto.EmployeeCount != "" || dto.Location != "" {
		return true
	}
	return false
}

func (dto *FlCompanyDTO) IsEmpty() bool {
	if dto.Name == "" && dto.Description == "" && dto.Funding == "" && dto.Industry == "" && dto.Investors == "" && dto.EmployeeCount == "" && dto.Location == "" {
		return true
	}
	return false
}

func (dto *FlCompanyDTO) ConvertFlCompanyDto(dao *repository.CompanyFunding) {
	ecSlice := pattern.FindAllString(dto.EmployeeCount, 1)
	var ec string
	e := -1
	if len(ecSlice) > 0 {
		ec = ecSlice[0]
		e, _ = strconv.Atoi(ec)
	}

	dao.Name = dto.Name
	dao.Description = dto.Description
	dao.Industry = dto.Industry
	dao.Investors = dto.Investors
	dao.EmployeeCount = e

	funding := strings.Split(dto.Funding, ",")
	// slice funding is length 1 if there is no comma
	if len(funding) > 1 {
		dao.FundingAmount = funding[0]
		dao.FundingRound = funding[1]
	} else {
		dao.FundingRound = funding[0]
	}

	location := strings.Split(dto.Location, ",")
	if len(location) > 0 {
		dao.LocationCity = location[0]
		dao.LocationMunicipality = location[1]
	} else {
		dao.LocationMunicipality = location[0]
	}

}

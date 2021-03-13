package dto

import "fmt"

type FlCompanyDTO struct {
	Name          string
	Description   string
	Funding       string
	Investors     string
	Location      string
	Industry      string
	EmployeeCount string
}

func (dto FlCompanyDTO) PrintDto() {
	fmt.Printf("Company: %s\n", dto.Name)
	fmt.Printf("Description: %s\n", dto.Description)
	fmt.Printf("Funding: %s\n", dto.Funding)
	fmt.Printf("Investors: %s\n", dto.Investors)
	fmt.Printf("Location: %s\n", dto.Location)
	fmt.Printf("Industry: %s\n", dto.Industry)
	fmt.Printf("EmployeeCount: %s\n", dto.EmployeeCount)
}

func (dto FlCompanyDTO) IsSaveAble() bool {
	if dto.Funding != "" || dto.Industry != "" || dto.Investors != "" || dto.EmployeeCount != "" || dto.Location != "" {
		return true
	}
	return false
}

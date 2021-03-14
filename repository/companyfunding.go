package repository

import (
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

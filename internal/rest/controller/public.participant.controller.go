package controller

import (
	cr "github.com/labovector/vecsys-api/internal/rest/repository/category"
	ir "github.com/labovector/vecsys-api/internal/rest/repository/institution"
	pr "github.com/labovector/vecsys-api/internal/rest/repository/payment"
	vr "github.com/labovector/vecsys-api/internal/rest/repository/referal"
	rr "github.com/labovector/vecsys-api/internal/rest/repository/region"
	ur "github.com/labovector/vecsys-api/internal/rest/repository/user"
)

type ParticipantController struct {
	ParticipantRepository ur.UserRepository
	CategoryRepository    cr.CategoryRepository
	RegionRepository      rr.RegionRepository
	PaymentRepository     pr.PaymentRepository
	ReferalRepository     vr.ReferalRepository
	InstitutionRepository ir.InstitutionRepository
}

func NewParticipantController(
	participantRepository ur.UserRepository,
	categoryRepository cr.CategoryRepository,
	regionRepository rr.RegionRepository,
	paymentRepository pr.PaymentRepository,
	referalRepository vr.ReferalRepository,
	institutionRepository ir.InstitutionRepository,
) *ParticipantController {
	return &ParticipantController{
		ParticipantRepository: participantRepository,
		CategoryRepository:    categoryRepository,
		RegionRepository:      regionRepository,
		PaymentRepository:     paymentRepository,
		ReferalRepository:     referalRepository,
		InstitutionRepository: institutionRepository,
	}
}

// TODO: Get All Participant Data
// TODO: Get Participant State
// TODO: Get All Event Category
// TODO: Get All Event Region
// TODO: Pick Category and Region
// TODO: Get All Payment Option
// TODO: Validate Referal
// TODO: Payment
// TODO: Get All Institution
// TODO: Add Institution
// TODO: Pick Institution
// TODO: Add Members (Receive an Array of Biodata)
// TODO: Lock Data

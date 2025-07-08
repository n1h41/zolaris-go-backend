package services

import (
	"context"
	"log"
	"github.com/afreedicp/zolaris-backend-app/internal/repositories"
)

// PolicyService handles business logic for IoT policy operations
type PolicyService struct {
	policyRepo *repositories.PolicyRepository
	policyName string
}

// NewPolicyService creates a new policy service instance
func NewPolicyService(policyRepo *repositories.PolicyRepository, policyName string) *PolicyService {
	return &PolicyService{
		policyRepo: policyRepo,
		policyName: policyName,
	}
}

// AttachIoTPolicy attaches an IoT policy to an identity
func (s *PolicyService) AttachIoTPolicy(ctx context.Context, identityID string) error {
	log.Printf("Attaching policy %s to identity %s", s.policyName, identityID)
	return s.policyRepo.AttachPolicy(ctx, s.policyName, identityID)
}


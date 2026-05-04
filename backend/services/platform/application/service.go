package application

import (
	"errors"
	"strings"

	"travel-api/services/platform/domain"
	"travel-api/services/platform/infrastructure"
)

var ErrInvalidPlatformRequest = errors.New("invalid_platform_request")

type PlatformService struct {
	repo *infrastructure.SQLitePlatformRepo
}

func NewPlatformService() *PlatformService {
	return &PlatformService{repo: infrastructure.NewSQLitePlatformRepo()}
}

func (s *PlatformService) Snapshot(userID, email string) (domain.Snapshot, error) {
	metrics, err := s.repo.Metrics()
	if err != nil {
		return domain.Snapshot{}, err
	}
	merchants, err := s.repo.ListMerchants()
	if err != nil {
		return domain.Snapshot{}, err
	}
	inventory, err := s.repo.ListInventory()
	if err != nil {
		return domain.Snapshot{}, err
	}
	orders, err := s.repo.ListOrders()
	if err != nil {
		return domain.Snapshot{}, err
	}
	refunds, err := s.repo.ListRefunds()
	if err != nil {
		return domain.Snapshot{}, err
	}
	cms, err := s.repo.ListCMS()
	if err != nil {
		return domain.Snapshot{}, err
	}
	profile := domain.UserProfile{}
	if strings.TrimSpace(userID) != "" {
		profile, _ = s.repo.GetProfile(userID, email)
	}
	return domain.Snapshot{Metrics: metrics, Merchants: merchants, Inventory: inventory, Orders: orders, Refunds: refunds, CMS: cms, Profile: profile}, nil
}

func (s *PlatformService) ListMerchants() ([]domain.Merchant, error) { return s.repo.ListMerchants() }
func (s *PlatformService) ListInventory() ([]domain.InventoryItem, error) {
	return s.repo.ListInventory()
}
func (s *PlatformService) UpdateInventory(req domain.InventoryUpdateRequest) (domain.InventoryItem, error) {
	if req.PackageID <= 0 || strings.TrimSpace(req.Date) == "" || req.Stock < 0 {
		return domain.InventoryItem{}, ErrInvalidPlatformRequest
	}
	return s.repo.UpdateInventory(req)
}
func (s *PlatformService) ListOrders() ([]domain.PlatformOrder, error) { return s.repo.ListOrders() }
func (s *PlatformService) CreateRefund(req domain.CreateRefundRequest) (domain.RefundRequest, error) {
	if strings.TrimSpace(req.UserID) == "" || req.OrderID <= 0 {
		return domain.RefundRequest{}, ErrInvalidPlatformRequest
	}
	if strings.TrimSpace(req.Reason) == "" {
		req.Reason = "Customer requested refund"
	}
	return s.repo.CreateRefund(req)
}
func (s *PlatformService) Profile(userID, email string) (domain.UserProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.UserProfile{}, ErrInvalidPlatformRequest
	}
	return s.repo.GetProfile(userID, email)
}
func (s *PlatformService) UpdateProfile(userID string, req domain.UserProfile) (domain.UserProfile, error) {
	if strings.TrimSpace(userID) == "" {
		return domain.UserProfile{}, ErrInvalidPlatformRequest
	}
	return s.repo.UpsertProfile(userID, req)
}
func (s *PlatformService) ListCMS() ([]domain.CMSArticle, error) { return s.repo.ListCMS() }
func (s *PlatformService) CreateCMS(req domain.CMSArticle) (domain.CMSArticle, error) {
	if strings.TrimSpace(req.Slug) == "" || strings.TrimSpace(req.Title) == "" {
		return domain.CMSArticle{}, ErrInvalidPlatformRequest
	}
	return s.repo.CreateCMS(req)
}
func (s *PlatformService) Metrics() (domain.DashboardMetrics, error) { return s.repo.Metrics() }

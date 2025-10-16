package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alireza0/s-ui/database"
	"github.com/alireza0/s-ui/database/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TelegramService struct{}

var (
	telegramListenersMu sync.RWMutex
	telegramListeners   []func()
	sharedService       = &TelegramService{}
)

func SharedTelegramService() *TelegramService {
	return sharedService
}

type TelegramConfigPayload struct {
	Enabled           bool              `json:"enabled"`
	BotToken          string            `json:"botToken"`
	WebhookDomain     string            `json:"webhookDomain"`
	WebhookSecret     string            `json:"webhookSecret"`
	YooKassaShopID    string            `json:"yooKassaShopId"`
	YooKassaSecretKey string            `json:"yooKassaSecretKey"`
	DownloadLinks     map[string]string `json:"downloadLinks"`
}

type TelegramTariffPayload struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PriceMinor   int64  `json:"priceMinor"`
	Currency     string `json:"currency"`
	DurationDays int    `json:"durationDays"`
	SortOrder    int    `json:"sortOrder"`
	Active       bool   `json:"active"`
}

type TelegramButtonPayload struct {
	ID        uint   `json:"id"`
	TariffID  uint   `json:"tariffId"`
	Label     string `json:"label"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
	SortOrder int    `json:"sortOrder"`
}

type TelegramConfigDTO struct {
	Enabled        bool              `json:"enabled"`
	BotTokenMasked string            `json:"botTokenMasked"`
	WebhookDomain  string            `json:"webhookDomain"`
	WebhookSecret  string            `json:"webhookSecret"`
	YooKassaShopID string            `json:"yooKassaShopId"`
	DownloadLinks  map[string]string `json:"downloadLinks"`
}

type TelegramButtonDTO struct {
	ID        uint   `json:"id"`
	TariffID  uint   `json:"tariffId"`
	Label     string `json:"label"`
	Action    string `json:"action"`
	Payload   string `json:"payload"`
	SortOrder int    `json:"sortOrder"`
}

type TelegramTariffDTO struct {
	ID           uint                `json:"id"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	PriceMinor   int64               `json:"priceMinor"`
	Currency     string              `json:"currency"`
	DurationDays int                 `json:"durationDays"`
	SortOrder    int                 `json:"sortOrder"`
	Active       bool                `json:"active"`
	Buttons      []TelegramButtonDTO `json:"buttons"`
}

type TelegramAdminState struct {
	Config        *TelegramConfigDTO            `json:"config"`
	Tariffs       []TelegramTariffDTO           `json:"tariffs"`
	Broadcasts    []TelegramBroadcastDTO        `json:"broadcasts"`
	PromoCodes    []TelegramPromoCodeDTO        `json:"promoCodes"`
	Conversations []TelegramConversationSummary `json:"conversations"`
}

type TelegramBroadcastAudience struct {
	AllUsers               bool   `json:"allUsers"`
	TariffIDs              []uint `json:"tariffIds"`
	IncludeNeverSubscribed bool   `json:"includeNeverSubscribed"`
	IncludeExpired         bool   `json:"includeExpired"`
}

type TelegramBroadcastPayload struct {
	ID       uint                      `json:"id"`
	Title    string                    `json:"title"`
	Body     string                    `json:"body"`
	Editable bool                      `json:"editable"`
	Audience TelegramBroadcastAudience `json:"audience"`
}

type TelegramBroadcastDTO struct {
	ID         uint                      `json:"id"`
	Title      string                    `json:"title"`
	Body       string                    `json:"body"`
	Editable   bool                      `json:"editable"`
	Status     string                    `json:"status"`
	Audience   TelegramBroadcastAudience `json:"audience"`
	SentAt     *time.Time                `json:"sentAt"`
	Deliveries int                       `json:"deliveries"`
	Success    int                       `json:"success"`
	Failed     int                       `json:"failed"`
}

type TelegramBroadcastEditPayload struct {
	BroadcastID uint   `json:"broadcastId"`
	Body        string `json:"body"`
}

type TelegramPromoCodePayload struct {
	ID              uint       `json:"id"`
	Code            string     `json:"code"`
	Description     string     `json:"description"`
	DiscountPercent int        `json:"discountPercent"`
	FreeDays        int        `json:"freeDays"`
	MaxUses         int        `json:"maxUses"`
	Active          bool       `json:"active"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	NoExpiry        bool       `json:"noExpiry"`
}

type TelegramPromoCodeDTO struct {
	ID              uint       `json:"id"`
	Code            string     `json:"code"`
	Description     string     `json:"description"`
	DiscountPercent int        `json:"discountPercent"`
	FreeDays        int        `json:"freeDays"`
	MaxUses         int        `json:"maxUses"`
	UsedCount       int        `json:"usedCount"`
	Active          bool       `json:"active"`
	ExpiresAt       *time.Time `json:"expiresAt"`
}

type TelegramUserDTO struct {
	ID                 uint       `json:"id"`
	TelegramID         int64      `json:"telegramId"`
	Username           string     `json:"username"`
	FirstName          string     `json:"firstName"`
	LastName           string     `json:"lastName"`
	Language           string     `json:"language"`
	EverPaid           bool       `json:"everPaid"`
	ActiveSubscription bool       `json:"activeSubscription"`
	SubscriptionEnds   *time.Time `json:"subscriptionEnds"`
	LastTariffID       *uint      `json:"lastTariffId"`
}

type TelegramConversationMessage struct {
	ID                uint      `json:"id"`
	Direction         string    `json:"direction"`
	Body              string    `json:"body"`
	TelegramMessageID string    `json:"telegramMessageId"`
	Seen              bool      `json:"seen"`
	CreatedAt         time.Time `json:"createdAt"`
}

type TelegramConversationDTO struct {
	User     TelegramUserDTO               `json:"user"`
	Messages []TelegramConversationMessage `json:"messages"`
}

type TelegramConversationSummary struct {
	User        TelegramUserDTO              `json:"user"`
	LastMessage *TelegramConversationMessage `json:"lastMessage"`
	UnreadCount int                          `json:"unreadCount"`
}

type TelegramInboundMessage struct {
	TelegramID            int64      `json:"telegramId"`
	Username              string     `json:"username"`
	FirstName             string     `json:"firstName"`
	LastName              string     `json:"lastName"`
	Language              string     `json:"language"`
	Message               string     `json:"message"`
	MessageID             string     `json:"messageId"`
	TariffID              *uint      `json:"tariffId"`
	EverPaid              bool       `json:"everPaid"`
	ActiveSubscription    bool       `json:"activeSubscription"`
	SubscriptionExpiresAt *time.Time `json:"subscriptionExpiresAt"`
}

func (p *TelegramConfigPayload) Validate() error {
	if p.Enabled {
		if strings.TrimSpace(p.BotToken) == "" {
			return errors.New("bot token is required when bot is enabled")
		}
		if strings.TrimSpace(p.YooKassaShopID) == "" {
			return errors.New("yookassa shop id is required when bot is enabled")
		}
		if strings.TrimSpace(p.YooKassaSecretKey) == "" {
			return errors.New("yookassa secret key is required when bot is enabled")
		}
	}
	return nil
}

func (p *TelegramTariffPayload) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("tariff title is required")
	}
	if strings.TrimSpace(p.Currency) == "" {
		return errors.New("currency is required")
	}
	if p.PriceMinor <= 0 {
		return errors.New("price must be greater than zero")
	}
	if p.DurationDays < 0 {
		return errors.New("duration must be zero or positive")
	}
	return nil
}

func (p *TelegramButtonPayload) Validate() error {
	if strings.TrimSpace(p.Label) == "" {
		return errors.New("button label is required")
	}
	if strings.TrimSpace(p.Action) == "" {
		return errors.New("button action is required")
	}
	return nil
}

func (a TelegramBroadcastAudience) Validate() error {
	if a.AllUsers {
		return nil
	}
	if len(a.TariffIDs) == 0 && !a.IncludeNeverSubscribed && !a.IncludeExpired {
		return errors.New("select at least one audience filter")
	}
	for _, id := range a.TariffIDs {
		if id == 0 {
			return errors.New("invalid tariff id in audience")
		}
	}
	return nil
}

func (p *TelegramBroadcastPayload) Validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("broadcast title is required")
	}
	if strings.TrimSpace(p.Body) == "" {
		return errors.New("broadcast body is required")
	}
	if err := p.Audience.Validate(); err != nil {
		return err
	}
	return nil
}

func (p *TelegramBroadcastEditPayload) Validate() error {
	if p.BroadcastID == 0 {
		return errors.New("broadcast id is required")
	}
	if strings.TrimSpace(p.Body) == "" {
		return errors.New("updated text is required")
	}
	return nil
}

func (p *TelegramPromoCodePayload) Validate() error {
	if strings.TrimSpace(p.Code) == "" {
		return errors.New("promo code is required")
	}
	if p.DiscountPercent < 0 || p.DiscountPercent > 100 {
		return errors.New("discount percent must be between 0 and 100")
	}
	if p.FreeDays < 0 {
		return errors.New("free days must be zero or greater")
	}
	if p.MaxUses < 0 {
		return errors.New("max uses must be zero or greater")
	}
	if !p.NoExpiry {
		if p.ExpiresAt == nil || p.ExpiresAt.IsZero() {
			return errors.New("expiration date is required or explicitly mark as no expiry")
		}
	}
	return nil
}

func (s *TelegramService) ensureConfig(tx *gorm.DB) (*model.TelegramBotConfig, error) {
	var cfg model.TelegramBotConfig
	err := tx.FirstOrCreate(&cfg, &model.TelegramBotConfig{ID: 1}).Error
	if err != nil {
		return nil, err
	}
	if cfg.DownloadLinks == "" {
		_ = cfg.SetDownloadLinks(map[string]string{})
		if err := tx.Save(&cfg).Error; err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}

func (s *TelegramService) GetConfig() (*model.TelegramBotConfig, error) {
	db := database.GetDB()
	return s.ensureConfig(db)
}

func (s *TelegramService) UpdateConfig(payload *TelegramConfigPayload) (*TelegramConfigDTO, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	db := database.GetDB()
	var result *model.TelegramBotConfig
	err := db.Transaction(func(tx *gorm.DB) error {
		cfg, err := s.ensureConfig(tx)
		if err != nil {
			return err
		}
		cfg.Enabled = payload.Enabled
		cfg.BotToken = strings.TrimSpace(payload.BotToken)
		cfg.WebhookDomain = strings.TrimSpace(payload.WebhookDomain)
		cfg.WebhookSecret = strings.TrimSpace(payload.WebhookSecret)
		cfg.YooKassaShopID = strings.TrimSpace(payload.YooKassaShopID)
		cfg.YooKassaSecretKey = strings.TrimSpace(payload.YooKassaSecretKey)
		if err := cfg.SetDownloadLinks(payload.DownloadLinks); err != nil {
			return err
		}
		if err := tx.Save(cfg).Error; err != nil {
			return err
		}
		result = cfg
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramConfigDTO(result), nil
}

func newTelegramConfigDTO(cfg *model.TelegramBotConfig) *TelegramConfigDTO {
	if cfg == nil {
		return nil
	}
	masked := ""
	if cfg.BotToken != "" {
		if len(cfg.BotToken) <= 8 {
			masked = strings.Repeat("*", len(cfg.BotToken))
		} else {
			masked = fmt.Sprintf("%s***%s", cfg.BotToken[:4], cfg.BotToken[len(cfg.BotToken)-4:])
		}
	}
	return &TelegramConfigDTO{
		Enabled:        cfg.Enabled,
		BotTokenMasked: masked,
		WebhookDomain:  cfg.WebhookDomain,
		WebhookSecret:  cfg.WebhookSecret,
		YooKassaShopID: cfg.YooKassaShopID,
		DownloadLinks:  cfg.GetDownloadLinks(),
	}
}

func (s *TelegramService) ListTariffs() ([]model.TelegramTariff, error) {
	db := database.GetDB()
	var tariffs []model.TelegramTariff
	err := db.Preload("Buttons", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc, id asc")
	}).Order("sort_order asc, id asc").Find(&tariffs).Error
	return tariffs, err
}

func (s *TelegramService) UpsertTariff(payload *TelegramTariffPayload) (*TelegramTariffDTO, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	db := database.GetDB()
	var tariff model.TelegramTariff
	if payload.ID != 0 {
		if err := db.First(&tariff, payload.ID).Error; err != nil {
			return nil, err
		}
	}
	tariff.Title = strings.TrimSpace(payload.Title)
	tariff.Description = strings.TrimSpace(payload.Description)
	tariff.PriceMinor = payload.PriceMinor
	tariff.Currency = strings.ToUpper(strings.TrimSpace(payload.Currency))
	tariff.DurationDays = payload.DurationDays
	tariff.SortOrder = payload.SortOrder
	tariff.Active = payload.Active
	if err := db.Save(&tariff).Error; err != nil {
		return nil, err
	}
	if err := db.Preload("Buttons", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order asc, id asc")
	}).First(&tariff, tariff.ID).Error; err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramTariffDTO(&tariff), nil
}

func (s *TelegramService) DeleteTariff(id uint) error {
	if id == 0 {
		return errors.New("tariff id is required")
	}
	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TelegramTariffButton{}, "tariff_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&model.TelegramTariff{}, id).Error; err != nil {
			return err
		}
		return nil
	})
	if err == nil {
		s.notifyChange()
	}
	return err
}

func (s *TelegramService) UpsertButton(payload *TelegramButtonPayload) (*TelegramButtonDTO, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	db := database.GetDB()
	if payload.ID == 0 {
		if payload.TariffID == 0 {
			return nil, errors.New("tariff id is required")
		}
		var exists bool
		err := db.Model(&model.TelegramTariff{}).Select("count(1) > 0").Where("id = ?", payload.TariffID).Find(&exists).Error
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("tariff %d not found", payload.TariffID)
		}
	}
	var button model.TelegramTariffButton
	if payload.ID != 0 {
		if err := db.First(&button, payload.ID).Error; err != nil {
			return nil, err
		}
	}
	button.TariffID = payload.TariffID
	button.Label = strings.TrimSpace(payload.Label)
	button.Action = strings.TrimSpace(payload.Action)
	button.Payload = strings.TrimSpace(payload.Payload)
	button.SortOrder = payload.SortOrder
	if err := db.Save(&button).Error; err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramButtonDTO(&button), nil
}

func (s *TelegramService) DeleteButton(id uint) error {
	if id == 0 {
		return errors.New("button id is required")
	}
	db := database.GetDB()
	if err := db.Delete(&model.TelegramTariffButton{}, id).Error; err != nil {
		return err
	}
	s.notifyChange()
	return nil
}

func (s *TelegramService) ListBroadcasts() ([]model.TelegramBroadcast, error) {
	db := database.GetDB()
	var broadcasts []model.TelegramBroadcast
	err := db.Preload("Deliveries").Order("created_at desc").Find(&broadcasts).Error
	return broadcasts, err
}

func (s *TelegramService) UpsertBroadcast(payload *TelegramBroadcastPayload) (*TelegramBroadcastDTO, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	db := database.GetDB()
	var broadcast model.TelegramBroadcast
	if payload.ID != 0 {
		if err := db.Preload("Deliveries").First(&broadcast, payload.ID).Error; err != nil {
			return nil, err
		}
		if broadcast.Status != "" && broadcast.Status != "draft" {
			return nil, errors.New("only draft broadcasts can be edited")
		}
	}
	audienceRaw, err := encodeAudience(payload.Audience)
	if err != nil {
		return nil, err
	}
	broadcast.Title = strings.TrimSpace(payload.Title)
	broadcast.Body = strings.TrimSpace(payload.Body)
	broadcast.Editable = payload.Editable
	broadcast.Audience = audienceRaw
	if strings.TrimSpace(broadcast.Status) == "" {
		broadcast.Status = "draft"
	}
	if err := db.Save(&broadcast).Error; err != nil {
		return nil, err
	}
	if err := db.Preload("Deliveries").First(&broadcast, broadcast.ID).Error; err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramBroadcastDTO(&broadcast, broadcast.Deliveries), nil
}

func (s *TelegramService) DeleteBroadcast(id uint) error {
	if id == 0 {
		return errors.New("broadcast id is required")
	}
	db := database.GetDB()
	var broadcast model.TelegramBroadcast
	if err := db.First(&broadcast, id).Error; err != nil {
		return err
	}
	if broadcast.Status != "draft" {
		return errors.New("only draft broadcasts can be deleted")
	}
	if err := db.Delete(&model.TelegramBroadcast{}, id).Error; err != nil {
		return err
	}
	s.notifyChange()
	return nil
}

func (s *TelegramService) SendBroadcast(id uint) (*TelegramBroadcastDTO, error) {
	if id == 0 {
		return nil, errors.New("broadcast id is required")
	}
	db := database.GetDB()
	var broadcast model.TelegramBroadcast
	if err := db.Preload("Deliveries").First(&broadcast, id).Error; err != nil {
		return nil, err
	}
	if broadcast.Status == "sent" {
		return nil, errors.New("broadcast already sent")
	}
	audience := decodeAudience(broadcast.Audience)
	users, err := s.selectAudienceUsers(audience)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("no users match broadcast audience")
	}
	now := time.Now()
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&broadcast).Updates(map[string]interface{}{
			"status":  "sent",
			"sent_at": &now,
		}).Error; err != nil {
			return err
		}
		broadcast.Status = "sent"
		broadcast.SentAt = &now
		for _, user := range users {
			messageID, sendErr := s.dispatchBroadcastMessage(tx, &broadcast, &user)
			status := "sent"
			errorMessage := ""
			if sendErr != nil {
				status = "failed"
				errorMessage = sendErr.Error()
			}
			delivery := model.TelegramBroadcastDelivery{
				BroadcastID:       broadcast.ID,
				UserID:            user.ID,
				TelegramMessageID: messageID,
				Status:            status,
				ErrorMessage:      errorMessage,
				SentAt:            &now,
			}
			if err := tx.Where("broadcast_id = ? AND user_id = ?", broadcast.ID, user.ID).
				Assign(delivery).FirstOrCreate(&delivery).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := db.Preload("Deliveries").First(&broadcast, broadcast.ID).Error; err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramBroadcastDTO(&broadcast, broadcast.Deliveries), nil
}

func (s *TelegramService) EditBroadcast(payload *TelegramBroadcastEditPayload) (*TelegramBroadcastDTO, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	db := database.GetDB()
	var broadcast model.TelegramBroadcast
	if err := db.Preload("Deliveries").First(&broadcast, payload.BroadcastID).Error; err != nil {
		return nil, err
	}
	if !broadcast.Editable {
		return nil, errors.New("broadcast is not marked as editable")
	}
	broadcast.Body = strings.TrimSpace(payload.Body)
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&broadcast).Error; err != nil {
			return err
		}
		for _, delivery := range broadcast.Deliveries {
			if delivery.Status != "sent" {
				continue
			}
			if _, err := s.editBroadcastMessage(tx, &broadcast, &delivery); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := db.Preload("Deliveries").First(&broadcast, broadcast.ID).Error; err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramBroadcastDTO(&broadcast, broadcast.Deliveries), nil
}

func (s *TelegramService) GetBroadcastDeliveries(broadcastID uint) ([]model.TelegramBroadcastDelivery, error) {
	if broadcastID == 0 {
		return nil, errors.New("broadcast id is required")
	}
	db := database.GetDB()
	var deliveries []model.TelegramBroadcastDelivery
	err := db.Where("broadcast_id = ?", broadcastID).Order("created_at asc").Find(&deliveries).Error
	return deliveries, err
}

func (s *TelegramService) ListPromoCodes() ([]model.TelegramPromoCode, error) {
	db := database.GetDB()
	var promos []model.TelegramPromoCode
	err := db.Order("created_at desc").Find(&promos).Error
	return promos, err
}

func (s *TelegramService) UpsertPromoCode(payload *TelegramPromoCodePayload) (*TelegramPromoCodeDTO, error) {
	if payload == nil {
		return nil, errors.New("payload is required")
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	db := database.GetDB()
	code := strings.ToUpper(strings.TrimSpace(payload.Code))
	var promo model.TelegramPromoCode
	if payload.ID != 0 {
		if err := db.First(&promo, payload.ID).Error; err != nil {
			return nil, err
		}
	}
	// ensure unique code
	var existing model.TelegramPromoCode
	err := db.Where("code = ? AND id <> ?", code, promo.ID).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err == nil && existing.ID != promo.ID && existing.ID != 0 {
		return nil, fmt.Errorf("promo code %s already exists", code)
	}
	promo.Code = code
	promo.Description = strings.TrimSpace(payload.Description)
	promo.DiscountPercent = payload.DiscountPercent
	promo.FreeDays = payload.FreeDays
	promo.MaxUses = payload.MaxUses
	promo.Active = payload.Active
	if payload.NoExpiry || payload.ExpiresAt == nil {
		promo.ExpiresAt = nil
	} else {
		exp := *payload.ExpiresAt
		promo.ExpiresAt = &exp
	}
	if err := db.Save(&promo).Error; err != nil {
		return nil, err
	}
	s.notifyChange()
	return newTelegramPromoCodeDTO(&promo), nil
}

func (s *TelegramService) DeletePromoCode(id uint) error {
	if id == 0 {
		return errors.New("promo code id is required")
	}
	db := database.GetDB()
	if err := db.Delete(&model.TelegramPromoCode{}, id).Error; err != nil {
		return err
	}
	s.notifyChange()
	return nil
}

func (s *TelegramService) ListConversations(limit int) ([]TelegramConversationSummary, error) {
	db := database.GetDB()
	var users []model.TelegramUserProfile
	query := db.Order("updated_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	summaries := make([]TelegramConversationSummary, 0, len(users))
	for _, user := range users {
		var last model.TelegramUserMessage
		err := db.Where("user_id = ?", user.ID).Order("created_at desc").First(&last).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		var unreadCount int64
		if err := db.Model(&model.TelegramUserMessage{}).
			Where("user_id = ? AND direction = ? AND seen = ?", user.ID, "inbound", false).
			Count(&unreadCount).Error; err != nil {
			return nil, err
		}
		summary := TelegramConversationSummary{
			User:        newTelegramUserDTO(&user),
			UnreadCount: int(unreadCount),
		}
		if err == nil {
			summary.LastMessage = newTelegramConversationMessageDTO(&last)
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (s *TelegramService) GetConversation(userID uint) (*TelegramConversationDTO, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}
	db := database.GetDB()
	var user model.TelegramUserProfile
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	var messages []model.TelegramUserMessage
	if err := db.Where("user_id = ?", user.ID).Order("created_at asc").Find(&messages).Error; err != nil {
		return nil, err
	}
	dto := &TelegramConversationDTO{
		User:     newTelegramUserDTO(&user),
		Messages: make([]TelegramConversationMessage, 0, len(messages)),
	}
	for i := range messages {
		dto.Messages = append(dto.Messages, *newTelegramConversationMessageDTO(&messages[i]))
	}
	_ = db.Model(&model.TelegramUserMessage{}).
		Where("user_id = ? AND direction = ?", user.ID, "inbound").
		Update("seen", true).Error
	return dto, nil
}

func (s *TelegramService) ReplyToConversation(userID uint, text string) (*TelegramConversationDTO, error) {
	if userID == 0 {
		return nil, errors.New("user id is required")
	}
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("message body is required")
	}
	db := database.GetDB()
	var user model.TelegramUserProfile
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	msg := model.TelegramUserMessage{
		UserID:            user.ID,
		Direction:         "outbound",
		Body:              strings.TrimSpace(text),
		TelegramMessageID: fmt.Sprintf("admin-%d-%d", user.TelegramID, now.UnixNano()),
		Seen:              true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		return tx.Model(&model.TelegramUserProfile{}).
			Where("id = ?", user.ID).
			Updates(map[string]interface{}{
				"updated_at":          now,
				"last_interaction_at": now,
			}).Error
	})
	if err != nil {
		return nil, err
	}
	return s.GetConversation(userID)
}

func (s *TelegramService) RecordInboundMessage(input *TelegramInboundMessage) (*TelegramConversationDTO, error) {
	if input == nil {
		return nil, errors.New("input is required")
	}
	if input.TelegramID == 0 {
		return nil, errors.New("telegram id is required")
	}
	if strings.TrimSpace(input.Message) == "" {
		return nil, errors.New("message body is required")
	}
	db := database.GetDB()
	now := time.Now()
	var profile model.TelegramUserProfile
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(model.TelegramUserProfile{TelegramID: input.TelegramID}).FirstOrCreate(&profile).Error; err != nil {
			return err
		}
		profile.Username = input.Username
		profile.FirstName = input.FirstName
		profile.LastName = input.LastName
		profile.Language = input.Language
		profile.EverPaid = profile.EverPaid || input.EverPaid
		profile.ActiveSubscription = input.ActiveSubscription
		profile.SubscriptionExpiresAt = input.SubscriptionExpiresAt
		profile.LastTariffID = input.TariffID
		profile.LastInteractionAt = now
		if err := tx.Save(&profile).Error; err != nil {
			return err
		}
		msg := model.TelegramUserMessage{
			UserID:            profile.ID,
			Direction:         "inbound",
			Body:              strings.TrimSpace(input.Message),
			TelegramMessageID: input.MessageID,
			Seen:              false,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
		if err := tx.Create(&msg).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetConversation(profile.ID)
}

func (s *TelegramService) selectAudienceUsers(a TelegramBroadcastAudience) ([]model.TelegramUserProfile, error) {
	db := database.GetDB()
	var users []model.TelegramUserProfile
	if err := db.Order("updated_at desc").Find(&users).Error; err != nil {
		return nil, err
	}
	if a.AllUsers {
		return users, nil
	}
	now := time.Now()
	selected := make([]model.TelegramUserProfile, 0, len(users))
	for _, user := range users {
		include := false
		if len(a.TariffIDs) > 0 && user.LastTariffID != nil {
			for _, id := range a.TariffIDs {
				if *user.LastTariffID == id {
					include = true
					break
				}
			}
		}
		if !include && a.IncludeNeverSubscribed && !user.EverPaid {
			include = true
		}
		if !include && a.IncludeExpired {
			if !user.ActiveSubscription && user.SubscriptionExpiresAt != nil && user.SubscriptionExpiresAt.Before(now) {
				include = true
			}
		}
		if include {
			selected = append(selected, user)
		}
	}
	return selected, nil
}

func (s *TelegramService) dispatchBroadcastMessage(tx *gorm.DB, broadcast *model.TelegramBroadcast, user *model.TelegramUserProfile) (string, error) {
	if user == nil {
		return "", errors.New("user is required")
	}
	if broadcast == nil {
		return "", errors.New("broadcast is required")
	}
	messageID := fmt.Sprintf("broadcast-%d-%d", broadcast.ID, user.TelegramID)
	now := time.Now()
	message := model.TelegramUserMessage{
		UserID:            user.ID,
		Direction:         "outbound",
		Body:              broadcast.Body,
		TelegramMessageID: messageID,
		Seen:              true,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := tx.Create(&message).Error; err != nil {
		return "", err
	}
	if err := tx.Model(&model.TelegramUserProfile{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"updated_at":          now,
		"last_interaction_at": now,
	}).Error; err != nil {
		return "", err
	}
	return messageID, nil
}

func (s *TelegramService) editBroadcastMessage(tx *gorm.DB, broadcast *model.TelegramBroadcast, delivery *model.TelegramBroadcastDelivery) (string, error) {
	if delivery == nil {
		return "", errors.New("delivery is required")
	}
	messageID := delivery.TelegramMessageID
	if messageID == "" {
		messageID = fmt.Sprintf("broadcast-%d-%d", broadcast.ID, delivery.UserID)
	}
	updates := map[string]interface{}{
		"body":       broadcast.Body,
		"updated_at": time.Now(),
	}
	result := tx.Model(&model.TelegramUserMessage{}).
		Where("user_id = ? AND telegram_message_id = ?", delivery.UserID, delivery.TelegramMessageID).
		Updates(updates)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		msg := model.TelegramUserMessage{
			UserID:            delivery.UserID,
			Direction:         "outbound",
			Body:              broadcast.Body,
			TelegramMessageID: messageID,
			Seen:              true,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}
		if err := tx.Create(&msg).Error; err != nil {
			return "", err
		}
	}
	return messageID, nil
}

func (s *TelegramService) GetAdminState() (*TelegramAdminState, error) {
	cfg, err := s.GetConfig()
	if err != nil {
		return nil, err
	}
	tariffs, err := s.ListTariffs()
	if err != nil {
		return nil, err
	}
	broadcasts, err := s.ListBroadcasts()
	if err != nil {
		return nil, err
	}
	promos, err := s.ListPromoCodes()
	if err != nil {
		return nil, err
	}
	conversations, err := s.ListConversations(0)
	if err != nil {
		return nil, err
	}
	state := &TelegramAdminState{
		Config:        newTelegramConfigDTO(cfg),
		Tariffs:       make([]TelegramTariffDTO, 0, len(tariffs)),
		Broadcasts:    make([]TelegramBroadcastDTO, 0, len(broadcasts)),
		PromoCodes:    make([]TelegramPromoCodeDTO, 0, len(promos)),
		Conversations: conversations,
	}
	for _, tariff := range tariffs {
		state.Tariffs = append(state.Tariffs, *newTelegramTariffDTO(&tariff))
	}
	for _, broadcast := range broadcasts {
		state.Broadcasts = append(state.Broadcasts, *newTelegramBroadcastDTO(&broadcast, broadcast.Deliveries))
	}
	for _, promo := range promos {
		state.PromoCodes = append(state.PromoCodes, *newTelegramPromoCodeDTO(&promo))
	}
	return state, nil
}

func newTelegramTariffDTO(t *model.TelegramTariff) *TelegramTariffDTO {
	if t == nil {
		return nil
	}
	dto := &TelegramTariffDTO{
		ID:           t.ID,
		Title:        t.Title,
		Description:  t.Description,
		PriceMinor:   t.PriceMinor,
		Currency:     t.Currency,
		DurationDays: t.DurationDays,
		SortOrder:    t.SortOrder,
		Active:       t.Active,
		Buttons:      make([]TelegramButtonDTO, 0, len(t.Buttons)),
	}
	for _, btn := range t.Buttons {
		dto.Buttons = append(dto.Buttons, *newTelegramButtonDTO(&btn))
	}
	return dto
}

func newTelegramButtonDTO(b *model.TelegramTariffButton) *TelegramButtonDTO {
	if b == nil {
		return nil
	}
	return &TelegramButtonDTO{
		ID:        b.ID,
		TariffID:  b.TariffID,
		Label:     b.Label,
		Action:    b.Action,
		Payload:   b.Payload,
		SortOrder: b.SortOrder,
	}
}

func normalizeAudience(a TelegramBroadcastAudience) TelegramBroadcastAudience {
	result := a
	if len(result.TariffIDs) > 0 {
		ids := make([]uint, 0, len(result.TariffIDs))
		seen := make(map[uint]struct{})
		for _, id := range result.TariffIDs {
			if id == 0 {
				continue
			}
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
		sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
		result.TariffIDs = ids
	}
	return result
}

func encodeAudience(a TelegramBroadcastAudience) (string, error) {
	normalized := normalizeAudience(a)
	raw, err := json.Marshal(normalized)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func decodeAudience(raw string) TelegramBroadcastAudience {
	if strings.TrimSpace(raw) == "" {
		return TelegramBroadcastAudience{AllUsers: true}
	}
	var audience TelegramBroadcastAudience
	if err := json.Unmarshal([]byte(raw), &audience); err != nil {
		return TelegramBroadcastAudience{AllUsers: true}
	}
	return normalizeAudience(audience)
}

func newTelegramBroadcastDTO(b *model.TelegramBroadcast, deliveries []model.TelegramBroadcastDelivery) *TelegramBroadcastDTO {
	if b == nil {
		return nil
	}
	dto := &TelegramBroadcastDTO{
		ID:       b.ID,
		Title:    b.Title,
		Body:     b.Body,
		Editable: b.Editable,
		Status:   b.Status,
		Audience: decodeAudience(b.Audience),
		SentAt:   b.SentAt,
	}
	successes := 0
	failures := 0
	for _, d := range deliveries {
		if d.Status == "sent" {
			successes++
		}
		if d.Status == "failed" {
			failures++
		}
	}
	dto.Deliveries = len(deliveries)
	dto.Success = successes
	dto.Failed = failures
	return dto
}

func newTelegramPromoCodeDTO(m *model.TelegramPromoCode) *TelegramPromoCodeDTO {
	if m == nil {
		return nil
	}
	return &TelegramPromoCodeDTO{
		ID:              m.ID,
		Code:            m.Code,
		Description:     m.Description,
		DiscountPercent: m.DiscountPercent,
		FreeDays:        m.FreeDays,
		MaxUses:         m.MaxUses,
		UsedCount:       m.UsedCount,
		Active:          m.Active,
		ExpiresAt:       m.ExpiresAt,
	}
}

func newTelegramUserDTO(u *model.TelegramUserProfile) TelegramUserDTO {
	if u == nil {
		return TelegramUserDTO{}
	}
	return TelegramUserDTO{
		ID:                 u.ID,
		TelegramID:         u.TelegramID,
		Username:           u.Username,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		Language:           u.Language,
		EverPaid:           u.EverPaid,
		ActiveSubscription: u.ActiveSubscription,
		SubscriptionEnds:   u.SubscriptionExpiresAt,
		LastTariffID:       u.LastTariffID,
	}
}

func newTelegramConversationMessageDTO(m *model.TelegramUserMessage) *TelegramConversationMessage {
	if m == nil {
		return nil
	}
	return &TelegramConversationMessage{
		ID:                m.ID,
		Direction:         m.Direction,
		Body:              m.Body,
		TelegramMessageID: m.TelegramMessageID,
		Seen:              m.Seen,
		CreatedAt:         m.CreatedAt,
	}
}

func (s *TelegramService) ExportBotState() (json.RawMessage, error) {
	state, err := s.GetAdminState()
	if err != nil {
		return nil, err
	}
	return json.Marshal(state)
}

func (s *TelegramService) LoadBotState(data json.RawMessage) error {
	if len(data) == 0 {
		return errors.New("empty payload")
	}
	var state TelegramAdminState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}
	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		cfg, err := s.ensureConfig(tx)
		if err != nil {
			return err
		}
		if state.Config != nil {
			cfg.Enabled = state.Config.Enabled
			cfg.WebhookDomain = state.Config.WebhookDomain
			cfg.WebhookSecret = state.Config.WebhookSecret
			cfg.YooKassaShopID = state.Config.YooKassaShopID
			if err := cfg.SetDownloadLinks(state.Config.DownloadLinks); err != nil {
				return err
			}
			if err := tx.Save(cfg).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec("DELETE FROM telegram_tariff_buttons").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM telegram_tariffs").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM telegram_broadcast_deliveries").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM telegram_broadcasts").Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM telegram_promo_codes").Error; err != nil {
			return err
		}
		for _, tariffDTO := range state.Tariffs {
			tariff := model.TelegramTariff{
				Title:        tariffDTO.Title,
				Description:  tariffDTO.Description,
				PriceMinor:   tariffDTO.PriceMinor,
				Currency:     tariffDTO.Currency,
				DurationDays: tariffDTO.DurationDays,
				SortOrder:    tariffDTO.SortOrder,
				Active:       tariffDTO.Active,
			}
			if err := tx.Omit(clause.Associations).Save(&tariff).Error; err != nil {
				return err
			}
			for _, btnDTO := range tariffDTO.Buttons {
				button := model.TelegramTariffButton{
					TariffID:  tariff.ID,
					Label:     btnDTO.Label,
					Action:    btnDTO.Action,
					Payload:   btnDTO.Payload,
					SortOrder: btnDTO.SortOrder,
				}
				if err := tx.Save(&button).Error; err != nil {
					return err
				}
			}
		}
		for _, broadcastDTO := range state.Broadcasts {
			audienceRaw, err := encodeAudience(broadcastDTO.Audience)
			if err != nil {
				return err
			}
			broadcast := model.TelegramBroadcast{
				Title:    broadcastDTO.Title,
				Body:     broadcastDTO.Body,
				Editable: broadcastDTO.Editable,
				Status:   broadcastDTO.Status,
				Audience: audienceRaw,
				SentAt:   broadcastDTO.SentAt,
			}
			if err := tx.Omit(clause.Associations).Save(&broadcast).Error; err != nil {
				return err
			}
		}
		for _, promoDTO := range state.PromoCodes {
			promo := model.TelegramPromoCode{
				Code:            promoDTO.Code,
				Description:     promoDTO.Description,
				DiscountPercent: promoDTO.DiscountPercent,
				FreeDays:        promoDTO.FreeDays,
				MaxUses:         promoDTO.MaxUses,
				UsedCount:       promoDTO.UsedCount,
				Active:          promoDTO.Active,
				ExpiresAt:       promoDTO.ExpiresAt,
			}
			if err := tx.Save(&promo).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err == nil {
		s.notifyChange()
	}
	return err
}

func (s *TelegramService) RegisterChangeListener(fn func()) {
	if fn == nil {
		return
	}
	telegramListenersMu.Lock()
	defer telegramListenersMu.Unlock()
	telegramListeners = append(telegramListeners, fn)
}

func (s *TelegramService) notifyChange() {
	telegramListenersMu.RLock()
	listeners := make([]func(), len(telegramListeners))
	copy(listeners, telegramListeners)
	telegramListenersMu.RUnlock()
	for _, fn := range listeners {
		fn()
	}
}

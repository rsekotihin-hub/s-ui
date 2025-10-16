package model

import (
	"encoding/json"
	"time"
)

type TelegramBotConfig struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Enabled           bool      `json:"enabled"`
	BotToken          string    `json:"botToken"`
	WebhookDomain     string    `json:"webhookDomain"`
	WebhookSecret     string    `json:"webhookSecret"`
	YooKassaShopID    string    `json:"yooKassaShopId"`
	YooKassaSecretKey string    `json:"yooKassaSecretKey"`
	DownloadLinks     string    `json:"downloadLinks" gorm:"type:text"`
	UpdatedAt         time.Time `json:"updatedAt"`
	CreatedAt         time.Time `json:"createdAt"`
}

type TelegramTariff struct {
	ID           uint                   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description" gorm:"type:text"`
	PriceMinor   int64                  `json:"priceMinor"`
	Currency     string                 `json:"currency"`
	DurationDays int                    `json:"durationDays"`
	SortOrder    int                    `json:"sortOrder"`
	Active       bool                   `json:"active"`
	Buttons      []TelegramTariffButton `json:"buttons" gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type TelegramTariffButton struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TariffID  uint      `json:"tariffId" gorm:"index"`
	Label     string    `json:"label"`
	Action    string    `json:"action"`
	Payload   string    `json:"payload" gorm:"type:text"`
	SortOrder int       `json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TelegramUserProfile struct {
	ID                    uint                        `json:"id" gorm:"primaryKey;autoIncrement"`
	TelegramID            int64                       `json:"telegramId" gorm:"uniqueIndex"`
	Username              string                      `json:"username"`
	FirstName             string                      `json:"firstName"`
	LastName              string                      `json:"lastName"`
	Language              string                      `json:"language"`
	Notes                 string                      `json:"notes" gorm:"type:text"`
	EverPaid              bool                        `json:"everPaid"`
	ActiveSubscription    bool                        `json:"activeSubscription"`
	SubscriptionExpiresAt *time.Time                  `json:"subscriptionExpiresAt"`
	LastTariffID          *uint                       `json:"lastTariffId"`
	LastInteractionAt     time.Time                   `json:"lastInteractionAt"`
	CreatedAt             time.Time                   `json:"createdAt"`
	UpdatedAt             time.Time                   `json:"updatedAt"`
	Messages              []TelegramUserMessage       `json:"messages" gorm:"constraint:OnDelete:CASCADE;"`
	BroadcastDeliveries   []TelegramBroadcastDelivery `json:"deliveries" gorm:"constraint:OnDelete:CASCADE;"`
}

type TelegramUserMessage struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID            uint      `json:"userId" gorm:"index"`
	Direction         string    `json:"direction" gorm:"index"`
	Body              string    `json:"body" gorm:"type:text"`
	TelegramMessageID string    `json:"telegramMessageId"`
	Seen              bool      `json:"seen"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type TelegramBroadcast struct {
	ID         uint                        `json:"id" gorm:"primaryKey;autoIncrement"`
	Title      string                      `json:"title"`
	Body       string                      `json:"body" gorm:"type:text"`
	Editable   bool                        `json:"editable"`
	Status     string                      `json:"status"`
	Audience   string                      `json:"audience" gorm:"type:text"`
	SentAt     *time.Time                  `json:"sentAt"`
	CreatedAt  time.Time                   `json:"createdAt"`
	UpdatedAt  time.Time                   `json:"updatedAt"`
	Deliveries []TelegramBroadcastDelivery `json:"deliveries" gorm:"constraint:OnDelete:CASCADE;"`
}

type TelegramBroadcastDelivery struct {
	ID                uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	BroadcastID       uint       `json:"broadcastId" gorm:"index"`
	UserID            uint       `json:"userId" gorm:"index"`
	TelegramMessageID string     `json:"telegramMessageId"`
	Status            string     `json:"status"`
	ErrorMessage      string     `json:"errorMessage" gorm:"type:text"`
	SentAt            *time.Time `json:"sentAt"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

type TelegramPromoCode struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code            string     `json:"code" gorm:"uniqueIndex"`
	Description     string     `json:"description" gorm:"type:text"`
	DiscountPercent int        `json:"discountPercent"`
	FreeDays        int        `json:"freeDays"`
	MaxUses         int        `json:"maxUses"`
	UsedCount       int        `json:"usedCount"`
	Active          bool       `json:"active"`
	ExpiresAt       *time.Time `json:"expiresAt"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (c *TelegramBotConfig) GetDownloadLinks() map[string]string {
	links := map[string]string{}
	if c.DownloadLinks == "" {
		return links
	}
	_ = json.Unmarshal([]byte(c.DownloadLinks), &links)
	return links
}

func (c *TelegramBotConfig) SetDownloadLinks(links map[string]string) error {
	if links == nil {
		links = map[string]string{}
	}
	raw, err := json.Marshal(links)
	if err != nil {
		return err
	}
	c.DownloadLinks = string(raw)
	return nil
}

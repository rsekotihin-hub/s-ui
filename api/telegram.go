package api

import (
	"net/http"

	"github.com/alireza0/s-ui/service"

	"github.com/gin-gonic/gin"
)

func (a *ApiService) GetTelegramState(c *gin.Context) {
	if a.TelegramService == nil {
		jsonObj(c, service.TelegramAdminState{}, nil)
		return
	}
	state, err := a.TelegramService.GetAdminState()
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, state, nil)
}

func (a *ApiService) SaveTelegramConfig(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload service.TelegramConfigPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	cfg, err := a.TelegramService.UpdateConfig(&payload)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, cfg, nil)
}

func (a *ApiService) SaveTelegramTariff(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload service.TelegramTariffPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	tariff, err := a.TelegramService.UpsertTariff(&payload)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, tariff, nil)
}

func (a *ApiService) DeleteTelegramTariff(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	if err := a.TelegramService.DeleteTariff(payload.ID); err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonMsg(c, "", nil)
}

func (a *ApiService) SaveTelegramButton(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload service.TelegramButtonPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	button, err := a.TelegramService.UpsertButton(&payload)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, button, nil)
}

func (a *ApiService) DeleteTelegramButton(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	if err := a.TelegramService.DeleteButton(payload.ID); err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonMsg(c, "", nil)
}

func (a *ApiService) SaveTelegramBroadcast(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload service.TelegramBroadcastPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	broadcast, err := a.TelegramService.UpsertBroadcast(&payload)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, broadcast, nil)
}

func (a *ApiService) DeleteTelegramBroadcast(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	if err := a.TelegramService.DeleteBroadcast(payload.ID); err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonMsg(c, "", nil)
}

func (a *ApiService) SendTelegramBroadcast(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	broadcast, err := a.TelegramService.SendBroadcast(payload.ID)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, broadcast, nil)
}

func (a *ApiService) EditTelegramBroadcast(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload service.TelegramBroadcastEditPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	broadcast, err := a.TelegramService.EditBroadcast(&payload)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, broadcast, nil)
}

func (a *ApiService) GetTelegramBroadcastDeliveries(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `form:"id" json:"id"`
	}
	if err := c.ShouldBindQuery(&payload); err != nil {
		if err := c.ShouldBindJSON(&payload); err != nil {
			jsonMsg(c, "", err)
			return
		}
	}
	deliveries, err := a.TelegramService.GetBroadcastDeliveries(payload.ID)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, deliveries, nil)
}

func (a *ApiService) SaveTelegramPromoCode(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload service.TelegramPromoCodePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	promo, err := a.TelegramService.UpsertPromoCode(&payload)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, promo, nil)
}

func (a *ApiService) DeleteTelegramPromoCode(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `json:"id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	if err := a.TelegramService.DeletePromoCode(payload.ID); err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonMsg(c, "", nil)
}

func (a *ApiService) GetTelegramConversation(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID uint `form:"id" json:"id"`
	}
	if err := c.ShouldBindQuery(&payload); err != nil {
		if err := c.ShouldBindJSON(&payload); err != nil {
			jsonMsg(c, "", err)
			return
		}
	}
	convo, err := a.TelegramService.GetConversation(payload.ID)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, convo, nil)
}

func (a *ApiService) ReplyTelegramConversation(c *gin.Context) {
	if a.TelegramService == nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}
	var payload struct {
		ID   uint   `json:"id"`
		Text string `json:"text"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		jsonMsg(c, "", err)
		return
	}
	convo, err := a.TelegramService.ReplyToConversation(payload.ID, payload.Text)
	if err != nil {
		jsonMsg(c, "", err)
		return
	}
	jsonObj(c, convo, nil)
}

package models

import "time"

type Game struct {
	GameID string `json:"game_id"`

	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type App struct {
	AppID  string `json:"app_id"`
	GameID string `json:"game_id"`

	Platform           string `json:"platform"`
	Market             string `json:"market"`
	DeviceType         string `json:"device_type"`
	Version            string `json:"version"`
	ForceUpdateEnabled bool   `json:"force_update_enabled"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	GooglePlayMarket = "GOOGLE_PLAY"
	TapTapMarket     = "TAP_TAP"
	HuaweiMarket     = "HUAWEI"
	AppStoreMarket   = "APP_STORE"
)

const (
	AndroidPlatform = "ANDROID"
	ApplePlatform   = "APPLE"
	WebPlatform     = "WEB"
)

const (
	PadDeviceType     = "PAD"
	TVDeviceType      = "TV"
	DesktopDeviceType = "DESKTOP"
	MobileDeviceType  = "MOBILE"
	BrowserDeviceType = "BROWSER"
)

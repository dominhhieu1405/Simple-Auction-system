package models

import "time"

const (
	AuctionStatusActive = "active"
	AuctionStatusClosed = "closed"
)

type Auction struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// Người tạo phiên đấu giá
	OwnerID uint `json:"owner_id" gorm:"not null;index"`
	Owner   User `json:"-" gorm:"foreignKey:OwnerID"`

	Title       string `json:"title" gorm:"size:200;not null"`
	Description string `json:"description" gorm:"type:text"`

	// Lưu tiền dưới dạng số nguyên, không dùng float64
	StartingPrice int64 `json:"starting_price" gorm:"not null"`
	MinIncrement  int64 `json:"min_increment" gorm:"not null"`

	EndTime time.Time `json:"end_time" gorm:"not null;index"`
	Status  string    `json:"status" gorm:"size:20;not null;default:active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

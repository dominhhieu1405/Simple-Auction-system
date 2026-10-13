package models

import "time"

type Bid struct {
	ID uint `json:"id" gorm:"primaryKey"`

	AuctionID uint    `json:"auction_id" gorm:"not null;index"`
	Auction   Auction `json:"-" gorm:"foreignKey:AuctionID"`

	UserID uint `json:"user_id" gorm:"not null;index"`
	User   User `json:"-" gorm:"foreignKey:UserID"`

	Amount int64 `json:"amount" gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
}

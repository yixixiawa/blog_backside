package Model

import "time"

type goods struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GameID            int       `gorm:"column:game_id" json:"game_id"`
	GameName          string    `gorm:"column:game_name" json:"game_name"`
	CommodityName     string    `gorm:"column:commodity_name;not null" json:"commodity_name"`
	CommodityHashName string    `gorm:"column:commodity_hash_name" json:"commodity_hash_name"`
	IconURL           string    `gorm:"column:icon_url" json:"icon_url"`
	OnSaleCount       int       `gorm:"column:on_sale_count" json:"on_sale_count"`
	Price             *float64  `gorm:"column:price" json:"price"`
	SteamPrice        *float64  `gorm:"column:steam_price" json:"steam_price"`
	SteamUsdPrice     *float64  `gorm:"column:steam_usd_price" json:"steam_usd_price"`
	TypeName          string    `gorm:"column:type_name" json:"type_name"`
	Exterior          string    `gorm:"column:exterior" json:"exterior"`
	ExteriorColor     string    `gorm:"column:exterior_color" json:"exterior_color"`
	Rarity            string    `gorm:"column:rarity" json:"rarity"`
	RarityColor       string    `gorm:"column:rarity_color" json:"rarity_color"`
	Quality           string    `gorm:"column:quality" json:"quality"`
	QualityColor      string    `gorm:"column:quality_color" json:"quality_color"`
	HaveLease         int       `gorm:"column:have_lease" json:"have_lease"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (goods) TableName() string {
	return "items"
}

// ItemQueryRequest 查询请求结构体
type ItemQueryRequest struct {
	Name     string   `form:"name" json:"name"`
	MinPrice *float64 `form:"min_price" json:"min_price"`
	MaxPrice *float64 `form:"max_price" json:"max_price"`
	GameID   *int     `form:"game_id" json:"game_id"`
	Exterior string   `form:"exterior" json:"exterior"`
	Rarity   string   `form:"rarity" json:"rarity"`
	Quality  string   `form:"quality" json:"quality"`
	Page     int      `form:"page" json:"page"`
	PageSize int      `form:"page_size" json:"page_size"`
	OrderBy  string   `form:"order_by" json:"order_by"` // price_asc, price_desc, name_asc, created_at_desc
}

// ItemStats 统计信息结构体
type ItemStats struct {
	Total       int64    `json:"total"`
	AvgPrice    *float64 `json:"avg_price"`
	MinPrice    *float64 `json:"min_price"`
	MaxPrice    *float64 `json:"max_price"`
	TotalOnSale int64    `json:"total_on_sale"`
}

// ItemResponse 响应结构体
type ItemResponse struct {
	Items []goods   `json:"items"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Size  int       `json:"size"`
	Stats ItemStats `json:"stats,omitempty"`
}

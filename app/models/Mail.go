package models

type Mail struct {
	Model
	ShopifyDomain string `json:"shopify_domain,omitempty" gorm:"type:nvarchar(100);not null"`
	Shop          *Shop  `json:"shop,omitempty" gorm:"foreignkey:ShopifyDomain"`
}

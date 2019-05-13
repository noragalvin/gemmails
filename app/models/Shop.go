package models

type Shop struct {
	Model
	ShopifyDomain string  `json:"shopify_domain;omitempty" gorm:"type:nvarchar(100);not null"`
	Mails         []*Mail `json:"mails;omitempty" gorm:"foreignkey:ShopifyDomain;association_foreignkey:ShopifyDomain"`
}

package models

// Shop ..
type Shop struct {
	Model
	ShopifyDomain string        `json:"shopify_domain;omitempty" gorm:"type:varchar(100);not null"`
	Subscribers   []*Subscriber `json:"subscribers;omitempty" gorm:"foreignkey:ShopifyDomain;association_foreignkey:ShopifyDomain"`
}

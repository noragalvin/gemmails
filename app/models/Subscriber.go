package models

// Subscriber ..
type Subscriber struct {
	Model
	ListID        string `json:"list_id,omitempty" gorm:"type:varchar(50)"`
	Email         string `json:"email,omitempty" gorm:"type:varchar(100);not null"`
	FName         string `json:"fname,omitempty" gorm:"type:varchar(50)"`
	LName         string `json:"lname,omitempty" gorm:"type:varchar(50)"`
	Phone         string `json:"phone,omitempty" gorm:"type:varchar(50)"`
	Address       string `json:"address,omitempty" gorm:"type:varchar(200)"`
	ShopifyDomain string `json:"shopify_domain,omitempty" gorm:"type:varchar(100);not null"`
	Shop          *Shop  `json:"shop,omitempty" gorm:"foreignkey:ShopifyDomain"`
}

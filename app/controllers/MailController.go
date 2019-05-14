package controllers

import (
	"encoding/json"
	"gemmails/app/helpers"
	"gemmails/app/models"
	v "gemmails/app/utils/view"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hanzoai/gochimp3"
)

// FormInfo ..
type FormInfo struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	FName   string `json:"fname"`
	LName   string `json:"lname"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// MailSend ..
func MailSend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	destination := vars["destination"]

	switch destination {
	case "mailchimp":
		err := sendToMailChimp(r)
		if err != nil {
			v.RespondBadRequest(w, v.Message(false, err.Error()))
			return
		}
	case "omnisend":
		err := sendToOmnisend(r)
		if err != nil {
			v.RespondBadRequest(w, v.Message(false, err.Error()))
			return
		}
	default:
		v.RespondBadRequest(w, v.Message(false, "Wrong destination"))
		return
	}

	v.RespondSuccess(w, v.Message(true, "success"))
}

func sendToMailChimp(r *http.Request) error {
	params := r.URL.Query()
	apiKey := params.Get("apiKey")
	listID := params.Get("listId")
	shopifyDomain := params.Get("shopifyDomain")
	version, err := strconv.Atoi(params.Get("version"))
	if err != nil {
		return err
	}

	formData := params.Get("formData")
	fname := ""
	lname := ""
	// formData: {"name":"Nora Galvin","fname":"nora", "lname": "galvin", "address": "Vietnam", "email": "clonenora01@gmail.com", "phone": "123456"}

	if version >= 2 {
		fname = ""
		lname = ""
	} else {
		fname = "John"
		lname = "Doe"
	}
	formInfo := FormInfo{}

	err = json.Unmarshal([]byte(formData), &formInfo)
	if err != nil {
		return err
	}
	fname = formInfo.FName
	lname = formInfo.LName
	if formInfo.Name != "" {
		arrName := strings.Split(formInfo.Name, " ")
		if len(arrName) > 1 {
			fname = arrName[0]
			lname = arrName[1]
		} else if len(arrName) == 1 {
			fname = arrName[0]
		}
	}

	data := make(map[string]interface{})
	data["FNAME"] = fname
	data["LNAME"] = lname
	data["PHONE"] = formInfo.Phone
	data["ADDRESS"] = formInfo.Address

	// Store subscriber to database
	// begin a transaction
	db := models.OpenDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	// http://localhost:8000/api/send-mail/mailchimp?version=2&listId=38eb0a9a92&apiKey=549867f2dd69ed0aaf9efad306eb5f71-us20&shopifyDomain=clonenora02.myshopify.com&formData={"fname":"nora2","lname":"galvin2","address":"Vietnam","email":"clonenora02@gmail.com","phone":"123456"}

	shop := models.Shop{}
	tx.Where("shopify_domain = ?", shopifyDomain).First(&shop)
	if shop.ID == 0 {
		shop.ShopifyDomain = shopifyDomain
		if err := tx.Create(&shop).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		//TODO: update shop info
	}

	// Find subscriber.  Update if exists. Create if not found
	subscriber := models.Subscriber{}
	tx.Where("email = ? AND shopify_domain = ? AND list_id = ?", formInfo.Email, shopifyDomain, listID).First(&subscriber)
	if subscriber.ID == 0 {
		s := models.Subscriber{
			ListID:        listID,
			FName:         fname,
			LName:         lname,
			Phone:         formInfo.Phone,
			Address:       formInfo.Address,
			ShopifyDomain: shopifyDomain,
			Email:         formInfo.Email,
		}
		if err := tx.Create(&s).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		subscriber.FName = fname
		subscriber.LName = lname
		subscriber.Phone = formInfo.Phone
		subscriber.Address = formInfo.Address
		if err := tx.Save(&subscriber).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Send to Mailchimp via Mailchimp api v3
	client := gochimp3.New(apiKey)

	// Add subscriber
	req := &gochimp3.MemberRequest{
		EmailAddress: formInfo.Email,
		Status:       "subscribed",
		MergeFields:  data,
	}

	// Fetch list
	list, err := client.GetList(listID, nil)
	if err != nil {
		return err
	}
	emailMD5 := helpers.GetMD5Hash(formInfo.Email)

	if _, err := list.AddOrUpdateMember(emailMD5, req); err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func sendToOmnisend(r *http.Request) error {
	// Api key: 5cda78668653ed3e50c96af9-zq91qjo93tzNn3BVY4Yr2Njl95HjCLFK2HUPFokznrrvjkwrK9
	// List ID: 5cda877f8653ed591c6056ec

	return nil
}

package controllers

import (
	"encoding/json"
	"gemmails/app/helpers"
	"gemmails/app/models"
	v "gemmails/app/utils/view"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hanzoai/gochimp3"
	"github.com/noragalvin/goklaviyo"
	"github.com/noragalvin/goomnisend"
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
		err := sendToMailChimp(r, destination)
		if err != nil {
			v.RespondBadRequest(w, v.Message(false, err.Error()))
			return
		}
	case "omnisend":
		err := sendToOmnisend(r, destination)
		if err != nil {
			v.RespondBadRequest(w, v.Message(false, err.Error()))
			return
		}
	case "klaviyo":
		err := sendToKlaviyo(r, destination)
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

func sendToMailChimp(r *http.Request, source string) error {
	params := r.URL.Query()
	apiKey := params.Get("apiKey")
	listID := params.Get("listId")
	shopifyDomain := params.Get("shopifyDomain")
	version, err := strconv.Atoi(params.Get("version"))
	if err != nil {
		return err
	}

	formInfo := FormInfo{}

	err = json.NewDecoder(r.Body).Decode(&formInfo)
	if err != nil {
		return err
	}

	if formInfo.Name != "" {
		arrName := strings.Split(formInfo.Name, " ")
		if len(arrName) > 1 {
			formInfo.FName = arrName[0]
			formInfo.LName = arrName[1]
		} else if len(arrName) == 1 {
			formInfo.FName = arrName[0]
		}
	}

	if version >= 2 {
	} else {
		if formInfo.FName == "" {
			formInfo.FName = "John"
		}
		if formInfo.LName == "" {
			formInfo.LName = "Doe"
		}
	}

	data := make(map[string]interface{})
	data["FNAME"] = formInfo.FName
	data["LNAME"] = formInfo.LName
	data["PHONE"] = formInfo.Phone
	data["ADDRESS"] = formInfo.Address

	go storeShopAndSubscriber(shopifyDomain, source, listID, formInfo)

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

	return nil
}

func sendToOmnisend(r *http.Request, source string) error {
	// Api key: 5cda78668653ed3e50c96af9-zq91qjo93tzNn3BVY4Yr2Njl95HjCLFK2HUPFokznrrvjkwrK9
	// List ID: 5cda877f8653ed591c6056ec
	params := r.URL.Query()
	apiKey := params.Get("apiKey")
	listID := params.Get("listId")
	shopifyDomain := params.Get("shopifyDomain")
	version, err := strconv.Atoi(params.Get("version"))
	if err != nil {
		return err
	}

	// formData: {"name":"Nora Galvin","fname":"nora", "lname": "galvin", "address": "Vietnam", "email": "clonenora01@gmail.com", "phone": "123456"}

	formInfo := FormInfo{}

	err = json.NewDecoder(r.Body).Decode(&formInfo)
	if err != nil {
		return err
	}

	if formInfo.Name != "" {
		arrName := strings.Split(formInfo.Name, " ")
		if len(arrName) > 1 {
			formInfo.FName = arrName[0]
			formInfo.LName = arrName[1]
		} else if len(arrName) == 1 {
			formInfo.FName = arrName[0]
		}
	}

	if version >= 2 {
	} else {
		if formInfo.FName == "" {
			formInfo.FName = "John"
		}
		if formInfo.LName == "" {
			formInfo.LName = "Doe"
		}
	}

	data := make(map[string]interface{})
	data["FNAME"] = formInfo.FName
	data["LNAME"] = formInfo.LName
	data["PHONE"] = formInfo.Phone
	data["ADDRESS"] = formInfo.Address

	// TODO: go routines store to database
	go storeShopAndSubscriber(shopifyDomain, source, listID, formInfo)

	// Send to Omnisend via omnisend api
	client := goomnisend.New(apiKey)

	reqParams := goomnisend.MemberRequest{}
	reqParams.Email = formInfo.Email
	reqParams.FirstName = formInfo.FName
	reqParams.LastName = formInfo.LName
	reqParams.Status = "subscribed"
	reqParams.StatusDate = time.Now()
	reqParams.Phone = formInfo.Phone
	reqParams.CustomerProperties = data
	list := &goomnisend.ListResponse{}
	list.ListID = listID
	lists := []*goomnisend.ListResponse{}
	lists = append(lists, list)
	reqParams.Lists = lists

	if _, err := client.CreateMember(&reqParams); err != nil {
		return err
	}

	return nil
}

func sendToKlaviyo(r *http.Request, source string) error {
	// 	apiKey := "pk_fd63d1b7a65f150c2d70994539db5077bb"
	// 	listID := "QGnDEg"
	params := r.URL.Query()
	apiKey := params.Get("apiKey")
	listID := params.Get("listId")
	shopifyDomain := params.Get("shopifyDomain")
	version, err := strconv.Atoi(params.Get("version"))
	if err != nil {
		return err
	}

	// formData: {"name":"Nora Galvin","fname":"nora", "lname": "galvin", "address": "Vietnam", "email": "clonenora01@gmail.com", "phone": "123456"}

	formInfo := FormInfo{}

	err = json.NewDecoder(r.Body).Decode(&formInfo)
	if err != nil {
		return err
	}

	if formInfo.Name != "" {
		arrName := strings.Split(formInfo.Name, " ")
		if len(arrName) > 1 {
			formInfo.FName = arrName[0]
			formInfo.LName = arrName[1]
		} else if len(arrName) == 1 {
			formInfo.FName = arrName[0]
		}
	}

	if version >= 2 {
	} else {
		if formInfo.FName == "" {
			formInfo.FName = "John"
		}
		if formInfo.LName == "" {
			formInfo.LName = "Doe"
		}
	}

	data := make(map[string]interface{})
	data["FNAME"] = formInfo.FName
	data["LNAME"] = formInfo.LName
	data["PHONE"] = formInfo.Phone
	data["ADDRESS"] = formInfo.Address

	// TODO: go routines store to database
	go storeShopAndSubscriber(shopifyDomain, source, listID, formInfo)

	// Send to Klaviyo via klaviyo api
	client := goklaviyo.New(apiKey)

	member := goklaviyo.MemberRequest{}
	member.Email = formInfo.Email
	member.FirstName = formInfo.FName
	member.LastName = formInfo.LName
	member.Phone = formInfo.Phone
	member.Address = formInfo.Address

	lists := []goklaviyo.MemberRequest{}
	lists = append(lists, member)

	memberSubscriber := goklaviyo.MemberSubscribeRequest{}
	memberSubscriber.Profiles = lists

	if _, err := client.CreateMember(listID, &memberSubscriber); err != nil {
		return err
	}

	return nil
}

func storeShopAndSubscriber(shopifyDomain, source, listID string, formInfo FormInfo) {
	// Store subscriber to database
	// begin a transaction
	db := models.OpenDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err == nil {
		// http://localhost:8000/api/send-mail/omnisend?version=2&listId=5cda877f8653ed591c6056ec&apiKey=5cda78668653ed3e50c96af9-zq91qjo93tzNn3BVY4Yr2Njl95HjCLFK2HUPFokznrrvjkwrK9&shopifyDomain=clonenora02.myshopify.com&formData={"fname":"nora2","lname":"galvin2","address":"Vietnam","email":"clonenora02@gmail.com","phone":"123456"}

		shop := models.Shop{}
		tx.Where("shopify_domain = ?", shopifyDomain).First(&shop)
		if shop.ID == 0 {
			shop.ShopifyDomain = shopifyDomain
			if err := tx.Create(&shop).Error; err != nil {
				tx.Rollback()
			}
		} else {
			//TODO: update shop info
		}

		// Find subscriber.  Update if exists. Create if not found
		subscriber := models.Subscriber{}
		tx.Where("email = ? AND shopify_domain = ? AND list_id = ? AND source = ?", formInfo.Email, shopifyDomain, listID, source).First(&subscriber)
		if subscriber.ID == 0 {
			s := models.Subscriber{
				ListID:        listID,
				FName:         formInfo.FName,
				LName:         formInfo.LName,
				Phone:         formInfo.Phone,
				Address:       formInfo.Address,
				ShopifyDomain: shopifyDomain,
				Email:         formInfo.Email,
				Source:        source,
			}
			if err := tx.Create(&s).Error; err != nil {
				tx.Rollback()
			}
		} else {
			subscriber.FName = formInfo.FName
			subscriber.LName = formInfo.LName
			subscriber.Phone = formInfo.Phone
			subscriber.Address = formInfo.Address
			if err := tx.Save(&subscriber).Error; err != nil {
				tx.Rollback()
			}
		}

		tx.Commit()
	}

}

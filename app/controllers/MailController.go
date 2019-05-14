package controllers

import (
	"encoding/json"
	"gemmails/app/helpers"
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

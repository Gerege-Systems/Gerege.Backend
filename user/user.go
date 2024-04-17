package user

import (
	"encoding/json"
	"fmt"

	"github.com/gerege-systems/gerege-backend/httpclient"
)

type ReqFind struct {
	SearchText     string `json:"search_text"`
	CountryCode    string `json:"country_code"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	BirthDate      string `json:"birth_date"`
	Gender         uint   `json:"gender"`
	PassportNumber string `json:"passport_number"`
	DocumentTypeId uint   `json:"document_type_id"`
	DateOfIssue    string `json:"date_of_issue"`
	DateOfExpire   string `json:"date_of_expire"`
	Mrz            string `json:"mrz"`
}

type ResFind struct {
	Id            uint   `json:"id"`
	CivilId       uint   `json:"civil_id"`
	RegNo         string `json:"reg_no"`
	FamilyName    string `json:"family_name"`
	LastName      string `json:"last_name"`
	FirstName     string `json:"first_name"`
	CountryCode   string `json:"country_code"`
	BirthDate     string `json:"birth_date"`
	Gender        uint   `json:"gender"`
	AddressDetail string `json:"address_detail"`
	Hash          string `json:"hash"`
	Email         string `json:"email"`
	PhoneNo       string `json:"phone_no"`
	ProfileImgUrl string `json:"profile_img_url"`
	AimagCode     string `json:"aimag_code"`
	AimagName     string `json:"aimag_name"`
	SumCode       string `json:"sum_code"`
	SumName       string `json:"sum_name"`
	BagCode       string `json:"bag_code"`
	BagName       string `json:"bag_name"`
}

func Find(req *ReqFind) (*ResFind, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpResponse := httpclient.Send(&httpclient.Config{
		Url:     "/user/find",
		Method:  "POST",
		Headers: &map[string]string{
			// "Authorization": "Bearer " + shared.Config.Api.Token,
		},
		Body: body,
	})

	if httpResponse.IsSuccess {
		user := ResFind{}
		if err := json.Unmarshal(httpResponse.Body, &user); err != nil {
			return nil, fmt.Errorf("cant parse user: %s", err.Error())
		}
		return &user, nil
	}

	return nil, fmt.Errorf("%s", httpResponse.Message)
}
package golkdr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hugmouse/golkdr/consts"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Receipt struct {
	Receipts []struct {
		Buyer                string      `json:"buyer"`
		BuyerType            string      `json:"buyerType"`
		CreatedDate          string      `json:"createdDate"`
		FiscalDocumentNumber string      `json:"fiscalDocumentNumber"`
		FiscalDriveNumber    string      `json:"fiscalDriveNumber"`
		FiscalSign           interface{} `json:"fiscalSign"`
		TotalSum             string      `json:"totalSum"`
		KktOwner             string      `json:"kktOwner"`
		KktOwnerInn          string      `json:"kktOwnerInn"`
		Key                  string      `json:"key"`
	} `json:"receipts"`
	HasMore bool `json:"hasMore"`
}

type Challenge struct {
	ChallengeToken             string    `json:"challengeToken"`
	ChallengeTokenExpiresIn    time.Time `json:"challengeTokenExpiresIn"`
	ChallengeTokenExpiresInSec int       `json:"challengeTokenExpiresInSec"`
}

type Count struct {
	NumberUnacknowledgedNotifications int `json:"numberUnacknowledgedNotifications"`
}

type AuthInfo struct {
	RefreshToken          string      `json:"refreshToken"`
	RefreshTokenExpiresIn interface{} `json:"refreshTokenExpiresIn"`
	Token                 string      `json:"token"`
	TokenExpireIn         time.Time   `json:"tokenExpireIn"`
	Profile               struct {
		TaxpayerPerson struct {
			Email         interface{} `json:"email"`
			Phone         string      `json:"phone"`
			Inn           interface{} `json:"inn"`
			FullName      interface{} `json:"fullName"`
			ShortName     interface{} `json:"shortName"`
			Status        string      `json:"status"`
			Address       interface{} `json:"address"`
			Oktmo         interface{} `json:"oktmo"`
			AuthorityCode interface{} `json:"authorityCode"`
			FirstName     interface{} `json:"firstName"`
			LastName      interface{} `json:"lastName"`
			MiddleName    interface{} `json:"middleName"`
		} `json:"taxpayerPerson"`
		AuthType string `json:"authType"`
	} `json:"profile"`
}

type Error struct {
	Code           string      `json:"code"`
	Message        string      `json:"message"`
	AdditionalInfo interface{} `json:"additionalInfo"`
	Err            error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

var ErrorFromAPI = errors.New("received an error in the response from the API")

type User struct {
	Number        uint
	Code          uint
	Authorization string
	AuthInfo      AuthInfo
	challenge     Challenge
}

func NewUser(num uint) (u *User) {
	return &User{
		Number:        num,
		Code:          0,
		Authorization: "Undefined",
		AuthInfo:      AuthInfo{},
		challenge:     Challenge{},
	}
}

func (u *User) RequestSMS() error {

	req, err := http.NewRequest("POST", consts.StartRoute, strings.NewReader(`{"phone":"`+strconv.Itoa(int(u.Number))+`"}`))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", u.Authorization)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(d))

	e := Error{}
	err = json.Unmarshal(d, &e)
	if err != nil {
		return err
	}

	if e.Code != "" {
		return &Error{
			Code:           e.Code,
			Message:        e.Message,
			AdditionalInfo: e.AdditionalInfo,
			Err:            ErrorFromAPI,
		}
	}

	err = json.Unmarshal(d, &u.challenge)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) SetCodeFromSMS(code int) error {
	jsond := fmt.Sprintf(consts.CodeFromSMS, u.challenge.ChallengeToken, u.Number, code)
	req, err := http.NewRequest("POST", consts.VerifyRoute, strings.NewReader(jsond))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", u.Authorization)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(d))

	e := Error{}
	err = json.Unmarshal(d, &e)
	if err != nil {
		return err
	}

	log.Println(string(d))

	if e.Code != "" {
		return &Error{
			Code:           e.Code,
			Message:        e.Message,
			AdditionalInfo: e.AdditionalInfo,
			Err:            ErrorFromAPI,
		}
	}

	err = json.Unmarshal(d, &u.AuthInfo)
	if err != nil {
		return err
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

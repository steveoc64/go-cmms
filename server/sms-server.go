package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"itrak-cmms/shared"
)

func GetSMSBalance() (int, error) {

	// if !Config.SMSOn {
	// 	return 0, nil
	// }

	if Config.SMSServer == "" {
		println("No Local SMS Server Defined")
		return 0, nil
	}

	resp, err := http.PostForm(
		Config.SMSServer,
		url.Values{
			"username": {Config.SMSUser},
			"password": {Config.SMSPasswd},
			"action":   {"balance"},
		})

	if err != nil {
		log.Println("HTTP Post Error", err.Error())
		return 0, err
	}

	/*	log.Println(resp)
		log.Println("status", resp.Status)
		log.Println("status code", resp.StatusCode)
		log.Println("proto", resp.Proto)
		log.Println("major", resp.ProtoMajor)
		log.Println("minor", resp.ProtoMinor)
		log.Println("header", resp.Header)
		log.Println("content length", resp.ContentLength)
		log.Println("transfer", resp.TransferEncoding)
		log.Println("trailer", resp.Trailer)
		log.Println("close", resp.Close)
		log.Println("req", resp.Request)
		log.Println("tls", resp.TLS)
		log.Println("body", resp.Body)
	*/

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	isok := string(body[:3])
	if isok != "OK:" {
		return 0, errors.New("Status Code is" + isok)
	}
	s := string(body[3:])

	return strconv.Atoi(s)
}

func GetIntlBalance() (int, error) {

	// if !Config.SMSOn {
	// 	return 0, nil
	// }

	if Config.SMSIntlServer == "" {
		println("No Intl SMS Server Defined")
		return 0, nil
	}

	// getbalURL := fmt.Sprintf("%s/user/get_credits/1/1.1?username=%s&password=%s",
	// 	Config.SMSIntlServer,
	// 	Config.SMSIntlUser,
	// 	Config.SMSIntlPasswd)
	// println("getbalURL = ", getbalURL)

	resp, err := http.Get(fmt.Sprintf("%s/user/get_credits/1/1.1?username=%s&password=%s",
		Config.SMSIntlServer,
		Config.SMSIntlUser,
		Config.SMSIntlPasswd))

	if err != nil {
		log.Println("HTTP Get Error", err.Error())
		return 0, err
	}

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	b := strings.TrimSpace(string(body))
	f := strings.Split(b, "|")
	// println("got response of", b, f[1])
	v, err := strconv.ParseFloat(f[1], 64)
	// println("converts to", v)
	return int(v), err
}

func SendSMS(number string, message string, ref string, user_id int) error {

	if !Config.SMSOn {
		return nil
	}

	// Get the user use_mobile flag, to see if they have Rx turned off for now

	theUser := shared.User{}
	userErr := DB.SQL("select username,use_mobile from users where id=$1", user_id).QueryStruct(&theUser)
	if userErr != nil {
		println("Cannot read user record for ID", user_id, userErr.Error())
		return nil
	}
	if !theUser.UseMobile {
		println("User ID", user_id, theUser.Username, "has requested no SMS transmissions")
		return nil
	}

	log.Println("Sending SMS to", number, ":", message)

	resp, err := http.PostForm(
		Config.SMSServer,
		url.Values{
			"username": {Config.SMSUser},
			"password": {Config.SMSPasswd},
			"to":       {number},
			"from":     {"SBS Intl"},
			"ref":      {ref},
			"message":  {message},
		})

	if err != nil {
		log.Println("HTTP Post Error", err.Error())
		return err
	}

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	lines := strings.Split(string(body), "\n")

	smsTrans := shared.SMSTrans{
		NumberTo: number,
		UserID:   user_id,
		Message:  message,
	}

	transID := 0

	for _, v := range lines {
		p := strings.Split(v, ":")
		smsTrans.Status = p[0]
		switch p[0] {
		case "OK":
			log.Println("SMS OK", p[1], "ref", p[2])
			smsTrans.NumberUsed = p[1]
			smsTrans.Ref = p[2]
			DB.InsertInto("sms_trans").
				Whitelist("number_to", "number_used", "user_id", "message", "ref", "status", "error").
				Record(smsTrans).
				Returning("id").
				QueryScalar(&transID)
			Connections.BroadcastAll("sms", "new", transID)
			return nil
		case "BAD":
			log.Println("SMS BAD", p[1], "reason", p[2])
			smsTrans.NumberUsed = p[1]
			smsTrans.Error = p[2]
			DB.InsertInto("sms_trans").
				Whitelist("number_to", "number_used", "user_id", "message", "ref", "status", "error").
				Record(smsTrans).
				Returning("id").
				QueryScalar(&transID)
			Connections.BroadcastAll("sms", "bad", transID)
			return errors.New(p[2])
		case "ERROR":
			log.Println("SMS ERROR", p[1])
			smsTrans.Error = p[1]
			DB.InsertInto("sms_trans").
				Whitelist("number_to", "number_used", "user_id", "message", "ref", "status", "error").
				Record(smsTrans).
				Returning("id").
				QueryScalar(&transID)
			Connections.BroadcastAll("sms", "error", transID)
			return errors.New(p[1])
			// default:
			// 	log.Println("Unknown SMS Error", p[0])
			// 	smsTrans.Error = p[1]
			// 	DB.InsertInto("sms_trans").
			// 		Whitelist("number_to", "number_used", "user_id", "message", "ref", "status", "error").
			// 		Record(smsTrans).
			// 		Exec()
			// 	return errors.New(p[1])
		}
	}
	return nil
}

type SMSRPC struct{}

func (s *SMSRPC) List(channel int, smsTrans *[]shared.SMSTrans) error {
	start := time.Now()

	conn := Connections.Get(channel)

	// Read the sites that this user has access to
	err := DB.SQL(`select * from sms_trans order by date_sent desc`).QueryStructs(smsTrans)

	if err != nil {
		log.Println(err.Error())
	}

	logger(start, "SMS.List",
		fmt.Sprintf("Channel %d, User %d %s %s",
			channel, conn.UserID, conn.Username, conn.UserRole),
		fmt.Sprintf("%d Messages", len(*smsTrans)),
		channel, conn.UserID, "sms_trans", 0, false)

	return nil
}

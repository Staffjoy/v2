package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/ttacon/libphonenumber"

	pb "v2.staffjoy.com/sms"
)

const (
	twilioURL     = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"
	defaultRegion = "US"
)

func (s *smsServer) Sender() {
	s.wg.Add(1)
	defer s.wg.Done()

	for task := range s.queue {
		// todo - can we requeue any of these?
		s.send(task)
		// Due to rate limit, put a second between sends
		time.Sleep(time.Second)
	}
}

func (s *smsServer) send(msg *pb.SmsRequest) {
	v := url.Values{}
	v.Set("To", msg.To)

	// Determine country code to send from
	toNum, err := libphonenumber.Parse(msg.To, defaultRegion)
	if err != nil {
		s.internalError(err, "unable to parse 'to' number for sms")
	}
	from, ok := s.sendingConfig.Numbers[*toNum.CountryCode]
	if !ok {
		s.logger.Warningf("Unsupported country code %d", toNum.CountryCode)

	}

	v.Set("From", from)

	v.Set("Body", msg.Body)
	rb := *strings.NewReader(v.Encode())

	// Create client
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(twilioURL, s.twilioSid), &rb)
	if err != nil {
		s.internalError(err, "failed to form twilio request")
	}
	req.SetBasicAuth(s.twilioSid, s.twilioToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		s.internalError(err, "failed to make twilio request")
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			s.logger.WithFields(logrus.Fields{"to": msg.To, "from": from, "body": msg.Body}).Infof("SMS sent - %v", data["sid"])
		}
	} else {
		s.internalError(fmt.Errorf("bad twilio response %v", resp.Status), "failed to send")
	}
}

package suite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// generated with https://mholt.github.io/json-to-go/
type usersResp struct {
	Data []struct {
		Active      bool        `json:"active"`
		Confirmed   bool        `json:"confirmed"`
		Email       string      `json:"email"`
		ID          int         `json:"id"`
		LastSeen    string      `json:"last_seen"`
		MemberSince string      `json:"member_since"`
		Name        string      `json:"name"`
		PhoneNumber interface{} `json:"phone_number"`
		Sudo        bool        `json:"sudo"`
		Username    string      `json:"username"`
	} `json:"data"`
	Filters struct {
		FilterByEmail    string      `json:"filterByEmail"`
		FilterByUsername interface{} `json:"filterByUsername"`
	} `json:"filters"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// AccountExists checks whether a suite user exists for the
// given email
func AccountExists(email string) (exists bool, err error) {
	u, ok := SuiteConfigs[config.Name]
	if !ok {
		err = fmt.Errorf("Unable to determine suite location")
		return
	}
	u.Path = "/api/v2/users/"
	q := u.Query()
	q.Set("filterByEmail", email)
	u.RawQuery = q.Encode()

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	req.SetBasicAuth(apiKey, "")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Unexpected suite status when querying users api - %s", resp.Status)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	data := new(usersResp)
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	logger.Debugf("suite response - %s %v", resp.Status, data)
	// user exists, is not sudo (to avoid annoyign redirects, has confirmed account, and active
	exists = len(data.Data) > 0 && !data.Data[0].Sudo && data.Data[0].Active && data.Data[0].Confirmed
	return
}

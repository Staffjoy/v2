package suite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// PerWeek holds a map of dates to worker scheduled count for said date
type PerWeek map[string]int

// OldData is the struct referencing old data found on the Suite python
// API - it holds people scheduled/week, people on shifts, etc.
type OldData struct {
	Data struct {
		PeopleClockedIn       int     `json:"people_clocked_in"`
		ScheduledPerWeek      PerWeek `json:"people_scheduled_per_week"`
		PeopleOnlineInLastDay int     `json:"people_online_in_last_day"`
		PeopleOnShifts        int     `json:"people_on_shifts"`
	} `json:"data"`
}

// GetOldData grabs v1 'suite' kpi data
func GetOldData() (*OldData, error) {
	cfg := config.Name
	if config.Name == "development" {
		cfg = "staging"
	}

	u, ok := SuiteConfigs[cfg]
	if !ok {
		return nil, fmt.Errorf("Unable to determine suite location")
	}
	u.Path = "/api/v2/internal/kpis/"
	q := u.Query()
	u.RawQuery = q.Encode()

	timeout := time.Duration(1 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", u.String(), nil)

	req.SetBasicAuth(apiKey, os.Getenv("SUITE_API_KEY"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected suite status when querying users api - %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	var od OldData
	if err = json.Unmarshal(body, &od); err != nil {
		return nil, err
	}

	return &od, nil
}

package main

import (
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/context"

	pb "v2.staffjoy.com/company"

	"io/ioutil"
	"strings"
)

// zoneList returns a list of available timezones on this system
func zoneList() ([]string, error) {
	b, err := ioutil.ReadFile("/usr/share/zoneinfo/zone.tab")
	if err != nil {
		return nil, err
	}

	s := string(b)

	var zones []string

	lines := strings.Split(s, "\n")
	for i := 0; i < len(lines); i++ {
		// comments in tab file
		if strings.Contains(lines[i], "#") {
			continue
		}
		parts := strings.Fields(lines[i])
		if len(parts) > 1 {
			zones = append(zones, parts[2])
		}
	}

	return zones, nil
}

func (s *companyServer) ListTimeZones(ctx context.Context, req *pb.TimeZoneListRequest) (*pb.TimeZoneList, error) {
	// no auth

	list, err := zoneList()
	if err != nil {
		return nil, s.internalError(err, "unable to get timezones")
	}

	res := &pb.TimeZoneList{}

	for _, tz := range list {
		res.Timezones = append(res.Timezones, tz)
	}
	return res, nil

}

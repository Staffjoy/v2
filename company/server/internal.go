package main

import (
	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"v2.staffjoy.com/auth"
	pb "v2.staffjoy.com/company"

	"sort"
	"strconv"
	"strings"

	"v2.staffjoy.com/suite"
)

type byDate []*pb.ScheduledPerWeek

func (a byDate) Len() int      { return len(a) }
func (a byDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byDate) Less(i, j int) bool {
	iwk, err := strconv.Atoi(strings.Replace(a[i].Week, "-", "", -1))
	if err != nil {
		logger.Debugf("%v", err)
	}

	jwk, err := strconv.Atoi(strings.Replace(a[j].Week, "-", "", -1))
	if err != nil {
		logger.Debugf("%v", err)
	}

	return iwk < jwk
}

// ScheduledPerWeek returns the weekly number of shifts/week
// it is the union of old data && new data
//
// it tries to query old data and merge that with current date - if
// fails then just returns local data
func (s *companyServer) ScheduledPerWeek() ([]*pb.ScheduledPerWeek, error) {
	var spwlist []*pb.ScheduledPerWeek
	var weeks []*pb.ScheduledPerWeek

	old, err := suite.GetOldData()
	if err != nil {
		logger.Errorf("failed to get old api data - %v", err)

		q := `select cast(a.weekname as char) as week, greatest(a.count, coalesce(b.count,0)) as count from (select 0 as count, str_to_date(concat(year(start), week(start), ' Monday'), '%X%V %W') as weekname from shift where start < NOW() group by weekname) as a left join (select count(distinct(user_uuid)) as count, str_to_date(concat(year(start), week(start), ' Monday'), '%X%V %W') as weekname from shift where start < NOW() and user_uuid != '' and published is true group by weekname) as b on a.weekname = b.weekname;`

		if _, err := s.dbMap.Select(&weeks, q); err != nil {
			return nil, s.internalError(err, "unable to query database")
		}

	} else {

		for k, v := range old.Data.ScheduledPerWeek {
			spw := &pb.ScheduledPerWeek{
				Week:  k,
				Count: int32(v),
			}
			spwlist = append(spwlist, spw)
		}

		// have unsorted data - so let's fix that first
		sort.Sort(byDate(spwlist))

		// take last..
		last := spwlist[len(spwlist)-1]

		// join with the rest..
		q := `select cast(a.weekname as char) as week, greatest(a.count, coalesce(b.count,0)) as count from (select 0 as count, str_to_date(concat(year(start), week(start), ' Monday'), '%X%V %W') as weekname from shift where start > ? and start < NOW() group by weekname) as a left join (select count(distinct(user_uuid)) as count, str_to_date(concat(year(start), week(start), ' Monday'), '%X%V %W') as weekname from shift where start > ? and start < NOW() and user_uuid != '' and published is true group by weekname) as b on a.weekname = b.weekname;`

		if _, err := s.dbMap.Select(&weeks, q, last.Week, last.Week); err != nil {
			return nil, s.internalError(err, "unable to query database")
		}

	}

	for _, wk := range weeks {
		spwlist = append(spwlist, wk)
	}

	return spwlist, nil
}

// PeopleOnShifts returns the count of people working right now
func (s *companyServer) PeopleOnShifts() (int32, error) {

	q := `select count(distinct(user_uuid)) from shift where shift.start <= now() and shift.stop > now() and user_uuid != "" and shift.published = true;`

	cnt, err := s.dbMap.SelectInt(q)
	if err != nil {
		return 0, err
	}

	return int32(cnt), nil
}

func (s *companyServer) GrowthGraph(ctx context.Context, req *pb.GrowthGraphRequest) (*pb.GrowthGraphResponse, error) {
	// Prep
	_, authz, err := getAuth(ctx)
	if err != nil {
		return nil, s.internalError(err, "Failed to authorize")
	}

	switch authz {
	case auth.AuthorizationSupportUser:
	default:
		return nil, grpc.Errorf(codes.PermissionDenied, "you do not have access to this service")
	}

	onShifts, err := s.PeopleOnShifts()
	if err != nil {
		return nil, s.internalError(err, "failed to query database")
	}

	perWeek, err := s.ScheduledPerWeek()
	if err != nil {
		return nil, s.internalError(err, "failed to query database")
	}

	// old api is un-ordered map - not sure if we want to continue that
	stuff := map[string]int32{}
	for i := 0; i < len(perWeek); i++ {
		stuff[perWeek[i].Week] = perWeek[i].Count
	}

	res := &pb.GrowthGraphResponse{
		PeopleScheduledPerWeek: stuff,
		PeopleOnShifts:         onShifts,
	}

	return res, nil

}

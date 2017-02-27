package main

import (
	"strings"
	"time"
	"v2.staffjoy.com/company"
)

// should eventually incorp. linting
// http://severinghaus.org/projects/icv/
//

func calDateFormat(t time.Time) string {
	return t.Format("20060102T150405Z")
}

// Cal contains the necessary information to implement a Staffjoy ical
// stream
type Cal struct {
	Shifts  []company.Shift
	Company string
}

func (cal *Cal) header() string {
	return `BEGIN:VCALENDAR
METHOD:PUBLISH
VERSION:2.0
PRODID:-//Staffjoy//Staffjoy Ical Service//EN
`
}

func (cal *Cal) body() string {
	body := ""

	for i := 0; i < len(cal.Shifts); i++ {
		body += `BEGIN:VEVENT
ORGANIZER;CN=Engineering:MAILTO:support@staffjoy.com
SUMMARY: Work at ` + cal.Company + `
UID:` + cal.Shifts[i].Uuid + `
STATUS:CONFIRMED
DTSTART:` + calDateFormat(cal.Shifts[i].Start) + `
DTEND:` + calDateFormat(cal.Shifts[i].Stop) + `
DTSTAMP:` + calDateFormat(time.Now()) + `
LAST-MODIFIED:` + calDateFormat(time.Now()) + `
LOCATION:  ` + cal.Company + `
END:VEVENT
`
	}

	return body
}

func (cal *Cal) footer() string {
	return `END:VCALENDAR`
}

// Build concats an ical header/body/footer together
func (cal *Cal) Build() string {
	o := cal.header() + cal.body() + cal.footer()
	return strings.Replace(o, "\n", "\r\n", -1)
}

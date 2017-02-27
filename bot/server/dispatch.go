package main

type dispatch int

const (
	dispatchSms = iota
	dispatchEmail
	dispatchUnavailable
)

func (u *user) PreferredDispatch() dispatch {
	// todo - check user notification preferences
	switch {
	case u.Phonenumber != "":
		return dispatchSms
	case u.Email != "":
		return dispatchEmail
	default:
		return dispatchUnavailable
	}
}

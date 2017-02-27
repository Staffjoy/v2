package main

import (
	"net/http"

	"v2.staffjoy.com/auth"
)

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	auth.Logout(res)
	http.Redirect(res, req, "/", http.StatusFound)
}

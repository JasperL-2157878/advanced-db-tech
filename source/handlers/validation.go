package handlers

import (
	"net/url"
	"regexp"
	"strconv"
)

func validateRouteParams(params url.Values) bool {
	var valid bool = true

	valid = valid && params.Has("from")
	valid = valid && params.Has("to")

	if _, err := strconv.Atoi(params.Get("from")); err == nil {
		valid = false
	}

	if _, err := strconv.Atoi(params.Get("to")); err == nil {
		valid = false
	}

	return valid
}

func validateGeocodeParams(params url.Values) bool {
	if !params.Has("address") {
		return false
	}

	re := regexp.MustCompile(`^([A-Za-zÀ-ÿ]+)(\s*)((\d*)?)(\s*)(,?)(\s*)((\d{4})?)(\s*)(([A-Za-zÀ-ÿ]+)?)(\s*)$`)

	return re.MatchString(params.Get("address"))
}

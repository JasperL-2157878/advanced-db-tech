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

	if _, err := strconv.Atoi(params.Get("from")); err != nil {
		valid = false
	}

	if _, err := strconv.Atoi(params.Get("to")); err != nil {
		valid = false
	}

	return valid
}

func validateGeocodeParams(params url.Values) bool {
	re := regexp.MustCompile(`^(?P<street_name>[^0-9,]+)\s*(?P<street_number>\d+)?[,\s]*(?P<postal_code>\d{4})?\s*(?P<city_name>\D+)?$`)
	return re.MatchString(params.Get("address"))
}

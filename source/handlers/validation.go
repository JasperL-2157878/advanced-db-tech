package handlers

import (
	"regexp"
	"strconv"
)

func (c Context) validateRouteParams() bool {
	var valid bool = true

	params := c.req.URL.Query()

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

func (c Context) validateGeocodeParams() bool {
	params := c.req.URL.Query()

	re := regexp.MustCompile(`^(?P<street_name>[^0-9,]+)\s*(?P<street_number>\d+)?[,\s]*(?P<postal_code>\d{4})?\s*(?P<city_name>\D+)?$`)
	return re.MatchString(params.Get("address"))
}

package main

import "regexp"

var PORT int = 8080

var (
	retailerPattern        = regexp.MustCompile(`^[\w\s\-&]+$`)
	totalPattern           = regexp.MustCompile(`^\d+\.\d{2}$`)
	itemDescriptionPattern = regexp.MustCompile(`^[\w\s\-]+$`)
	pricePattern           = regexp.MustCompile(`^\d+\.\d{2}$`)
	datePattern            = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
	timePattern            = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)
)

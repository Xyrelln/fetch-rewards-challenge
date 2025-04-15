package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func newReceiptManager() *receiptManager {
	return &receiptManager{
		wallet: make(map[string]int),
	}
}

// given a parsed receipt, returns its assigned receipt id, and error if receipt does not comply with OpenAPI docs
func (rm *receiptManager) calculateAndSavePoints(rcpt receipt) (string, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	points := 0

	// One point for every alphanumeric character in the retailer name
	if !retailerPattern.MatchString(rcpt.Retailer) {
		return "", fmt.Errorf("retailer field %q must comply '^[\\w\\s\\-&]+$'", rcpt.Retailer)
	}
	for _, r := range rcpt.Retailer {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			points++
		}
	}

	// 50 points if the total is a round dollar amount with no cents
	if !totalPattern.MatchString(rcpt.Total) {
		return "", fmt.Errorf("total field %q must comply '^\\d+\\.\\d{2}$'", rcpt.Total)
	}
	cents, _ := strconv.Atoi(strings.Split(rcpt.Total, ".")[1]) // no need to check again err
	if cents == 0 {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if cents%25 == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(rcpt.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned
	for _, item := range rcpt.Items {
		// check item short description pattern
		if !itemDescriptionPattern.MatchString(item.ShortDescription) { // check item short description matches regex
			return "", fmt.Errorf("item short description field %q must comply '^[\\w\\s\\-]+$'", item.ShortDescription)
		}
		// check item price format
		if !pricePattern.MatchString(item.Price) {
			return "", fmt.Errorf("item %q has invalid price %q, which doesn't comply '^\\d+\\.\\d{2}$'", item.ShortDescription, item.Price)
		}
		trimmed := strings.TrimSpace(item.ShortDescription)

		if len(trimmed)%3 == 0 {
			parts := strings.Split(item.Price, ".")
			dollars, _ := strconv.Atoi(parts[0])
			cents, _ := strconv.Atoi(parts[1])
			itemPriceCents := dollars*100 + cents
			points += int(math.Ceil(0.2 * float64(itemPriceCents) / 100))
		}
	}

	// 6 points if the day in the purchase date is odd
	if !datePattern.MatchString(rcpt.PurchaseDate) {
		return "", fmt.Errorf("date format should be yyyy-mm-dd")
	}
	day := strings.Split(rcpt.PurchaseDate, "-")[2]
	d, _ := strconv.Atoi(day) // we checked err with regex
	if d%2 == 1 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	if !timePattern.MatchString(rcpt.PurchaseTime) {
		return "", fmt.Errorf("time format should be hh:mm")
	}
	h, _ := strconv.Atoi(strings.Split(rcpt.PurchaseTime, ":")[0])
	if h >= 14 && h < 16 {
		points += 10
	}

	var rid string = generateReceiptId()
	rm.wallet[rid] = points
	return rid, nil
}

// get Points given an id. If not found, point will be -1 and error won't be nil
func (rm *receiptManager) getPoint(rid string) (int, error) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	point, exist := rm.wallet[rid]
	if !exist {
		return -1, fmt.Errorf("no receipt found for that ID")
	}
	return point, nil
}

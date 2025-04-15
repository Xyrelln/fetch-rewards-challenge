package main

import (
	"sync"
)

type Cents int

type receiptManager struct {
	wallet map[string]int
	mu     sync.RWMutex
}

type receipt struct {
	Retailer     string
	PurchaseDate string
	PurchaseTime string
	Items        []item
	Total        string // it's a receipt, we parse the string later like a human
}

type item struct {
	Price            string // string
	ShortDescription string
}

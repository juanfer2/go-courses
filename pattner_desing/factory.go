package main

import (
	"fmt"
)

type IProduct interface {
	SetStock(stock int)
	GetStock() int
	SetName(name string)
	GetName() string
}

type Computer struct {
	name  string
	stock int
}

func (c *Computer) SetStock(stock int) {
	c.stock = stock
}

func (c *Computer) GetStock() int {
	return c.stock
}

func (c *Computer) SetName(name string) {
	c.name = name
}

func (c *Computer) GetName() string {
	return c.name
}

type Lapto struct {
	Computer
}

func NewLapto(name string, stock int) IProduct {
	return &Lapto{
		Computer: Computer{
			name:  name,
			stock: stock,
		},
	}
}

func main() {
	mac := NewLapto("MacBook", 30)
	hp := NewLapto("LaptoHP", 25)

	list := []IProduct{mac, hp}

	for _, computer := range list {
		fmt.Printf("Product name: %s, with stock %d\n", computer.GetName(), computer.GetStock())
	}
}

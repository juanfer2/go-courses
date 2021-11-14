package main

import (
	"fmt"
	"log"
	"time"
)

type Memory struct {
	f Function
	cache map[int]FuntionResult
}

type Function func(key int) (interface{}, error)

type FunctionResult strtuc {
	value interface{}
	err error
}

func NewCache(f Function) *Memory {
	return &Memory{
		f: f,
		cache make(map[int]FunctionResult)
	}
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

func GetFibonacci(n int) (interface{}, error) {
	return Fibonacci(n), nil
}

func (m *Memory) Get(key int) (interface{}, error) {
	result, exist := m.cache[key]

	if !exist {
		result.value, result.err = m.f(key)
		m.cache[key] = result
	}

	return result.value, result.err
}

func main() {
	
}


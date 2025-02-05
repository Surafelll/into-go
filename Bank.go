package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Account struct {
	Balance float64 `json:"balance"`
}

const dataFile = "account.json"

func (a *Account) Save() {
	data, _ := json.MarshalIndent(a, "", "  ")
	ioutil.WriteFile(dataFile, data, 0644)
}

func LoadAccount() *Account {
	var account Account
	if data, err := ioutil.ReadFile(dataFile); err == nil {
		json.Unmarshal(data, &account)
	}
	return &account
}

func (a *Account) Deposit(amount float64) {
	if amount > 0 {
		a.Balance += amount
		a.Save()
		fmt.Printf("Deposited: $%.2f\n", amount)
	} else {
		fmt.Println("Invalid deposit amount")
	}
}

func (a *Account) Withdraw(amount float64) {
	if amount > 0 && amount <= a.Balance {
		a.Balance -= amount
		a.Save()
		fmt.Printf("Withdrawn: $%.2f\n", amount)
	} else {
		fmt.Println("Invalid or insufficient funds")
	}
}

func (a *Account) CheckBalance() {
	fmt.Printf("Current Balance: $%.2f\n", a.Balance)
}

func main() {
	account := LoadAccount()
	account.Deposit(1000)
	account.Withdraw(200)
	account.CheckBalance()
}
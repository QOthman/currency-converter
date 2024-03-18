package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type RateResponse struct {
	Conversion_rates map[string]float64 `json:"conversion_rates"`
}

func toUpper(s string)string{
	runes := []rune(s)
	for i := 0; i < len(s); i++ {
		if s[i]>='a' && s[i]<='z' {
			runes[i] = rune(runes[i] - 32)
		}
	}
	return string(runes)
}

func check(base , conv string , list map[string]float64) bool {
	t := 0
	for index := range list  {
		if index == base || index == conv {
			t++
		}
	}
	if t == 2 {
		return true
	}else {
		return false
	}
}

func main() {

	args := os.Args[1:]

	if (len(args) == 0) || (len(args) == 1 && (args[0] == "-h" || args[0] == "--help")) {

		short := `Usage: currency_converter <amount> <base> <target>`
		long := `Currency Converter Help:

	Usage: currency_converter <amount> <base> <target>
		
	Description:
	Converts an amount from one currency to another using the latest exchange rates.
		
	Arguments:
	- amount: Amount to convert (positive number).
	- base: Base currency code (e.g., USD, EUR).
	- target: Target currency code.
		
	Example:
	1.00 usd = 10.05 MAD
		
	Note:
	- Requires internet connection.
		`

		if len(args) == 1 && args[0] == "-h" {
			fmt.Println(short)
		}else {
			fmt.Println(long)
		}

	} else if len(args) == 2 || len(args) == 3 {
		amount := 0.0
		var base, conv string
		if len(args) == 3 {
			t, _ := strconv.Atoi(args[0])
			amount = float64(t)
			base = args[1]
			conv = args[2]
		} else {
			amount = 1
			base = args[0]
			conv = args[1]
		}

		link := "https://v6.exchangerate-api.com/v6/a507912f33b9279524a1d977/latest/" + base
		resp, err := http.Get(link)
		if err != nil {
			fmt.Println("Error fetching data:", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var Response RateResponse
		err = json.Unmarshal(body, &Response)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}
		base = toUpper(base)
		conv = toUpper(conv)
		rate := Response.Conversion_rates[conv]

		if check(base , conv , Response.Conversion_rates) {
			fmt.Printf("%.2f %s = %.2f %s\n", amount, base, amount*rate, conv)
		} else {
			fmt.Printf("Conversion rate for %s not available\n", conv)
		}
	}
}

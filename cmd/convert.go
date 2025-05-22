package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

var (
	fromCurrency string
	toCurrency   string
	amount       float64
)

var convertCommand = &cobra.Command{
	Use:   "convert",
	Short: "Convert currency from one format to another",
	Run: func(cmd *cobra.Command, args []string) {
		if fromCurrency == "" || toCurrency == "" || amount <= 0 {
			fmt.Println("Please provide --from, --to and --amount")
			return
		}
		rate, err := getExchangeRate(fromCurrency, toCurrency)
		if err != nil {
			fmt.Println("error", err)
		}

		totalAmount := amount * rate
		fmt.Printf("Total amount is %.4f \n", totalAmount)
	},
}

func init() {
	convertCommand.Flags().StringVar(&fromCurrency, "from", "", "Base currency (e.g. USD)")
	convertCommand.Flags().StringVar(&toCurrency, "to", "", "Target currency (e.g. EUR)")
	convertCommand.Flags().Float64Var(&amount, "amount", 0, "Amount to convert")

	rootCmd.AddCommand(convertCommand)

}

type ExchangeRateResponse struct {
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func getExchangeRate(from, to string) (float64, error) {
	var API_SECRET = os.Getenv("API_SECRET")
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/%s", API_SECRET, from)
	res, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("error in decoding", err)
		return 0, nil
	}

	var data ExchangeRateResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing json ")
		return 0, nil
	}

	rate, ok := data.ConversionRates[to]
	if !ok {
		return 0, fmt.Errorf("conversion rate for %s not found", to)
	}

	return rate, nil
}

package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

type Product struct {
	Sku   string `json:"sku"`
	Stock int    `json:"stock_quantity"`
	Price string `json:"price"`
}

func main() {
	app := &cli.App{
		Name:  "Syncronizador cicclos",
		Usage: "sync [cmd]",
		Commands: []*cli.Command{
			{
				Name:    "get_products",
				Aliases: []string{"gp"},
				Usage:   "sync get_products",
				Action: func(c *cli.Context) error {
					print_csv(get_products())
					return nil
				},
			},
			{
				Name:    "get_products_by_skus",
				Aliases: []string{"gps"},
				Usage:   "sync get_products_by_skus [skus...]",
				Action: func(c *cli.Context) error {
					products := get_products_by_skus(c.Args().Slice())
					print_csv(products)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func get_products() []Product {

	user, ok := os.LookupEnv("CICCLOS_USER")
	if !ok {
		panic("CICCLOS_USER not set")
	}

	pass, ok := os.LookupEnv("CICCLOS_PASS")

	if !ok {
		panic("CICCLOS_PASS not set")
	}

	credentials := user + ":" + pass

	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

	req, err := http.NewRequest("GET", "https://cicclos.com.ar/wp-json/wc/v3/products", nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Basic "+encoded)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	var products []Product

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&products)

	if err != nil {
		panic(err)
	}

	return products
}

func get_products_by_skus(skus []string) []Product {
	products := get_products()

	var res []Product

	for _, sku := range skus {
		for _, product := range products {
			if product.Sku == sku {
				res = append(res, product)
			}
		}
	}

	return res
}

func print_csv(products []Product) {
	var res [][]string

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	for _, product := range products {

		row := []string{product.Sku, strconv.Itoa(product.Stock), product.Price}
		res = append(res, row)

	}

	w.WriteAll(res)
}


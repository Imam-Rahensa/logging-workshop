package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/imam-rahensa/logging-workshop/external"
	"github.com/tokopedia/tdk/go/log"
	"encoding/json"
)

func main() {

	logConfig := &log.Config{
		Level:     "trace",                            // Default will be in info level
		LogFile:   "./log/logging-workshop.error.log", // If none supplied will goes to os.Stderr, for production you must put log file
		DebugFile: "./log/logging-workshop.debug.log", // If none supplied will goes to os.Stderr, for production you must put log file
		AppName:   "logging-workshop",                 //  your app name, the format will be `{service_name}_{function}`
	}

	// Set Standard Log Config
	err := log.SetStdLog(logConfig)

	if err != nil {
		log.StdErrorf(context.Background(), nil, err, "Failed to initialize log configuration, with config: %s.", json.Marshal(logConfig))
	}

	http.HandleFunc("/", HelloHandler)
	http.ListenAndServe(":8080", nil)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Init context for logging. This will inject request_id to your context
	ctx = log.InitLogContext(ctx)

	var (
		productID int
		err       error
	)

	keys, ok := r.URL.Query()["product_id"]
	if !ok {
		log.StdError(ctx, nil, nil, "No product id supplied")
		fmt.Fprint(w, "No product id supplied")
		return
	}

	// parse the product id
	if len(keys) < 1 {
		log.StdError(ctx, nil, nil, "No product is found to be associated with the supplied product ID.")
		fmt.Fprint(w, "No product is found to be associated with the supplied product ID.")
		return
	}

	productID, err = strconv.Atoi(keys[0])
	if err != nil {
		log.StdErrorf(ctx, nil, nil, "Failed to parse the product ID, with key: %s", json.Marshal(keys))
		fmt.Fprintf(w, "Failed to parse the product ID, with key: %s", json.Marshal(keys))
		return
	}

	// Set your context id. In this case you will put your product id
	ctx = log.SetCtxID(ctx, strconv.Itoa(productID))

	product, err := GetProductFromDB(ctx, productID)
	if err != nil {
		log.StdErrorf(ctx, nil, err, "Failed to get product from database, with product ID: %s", productID)
		fmt.Fprintf(w, "Failed to get product from database, with product ID: %s", productID)
		return
	}

	err = CalculateDiscount(ctx, product)
	if err != nil {
		log.StdErrorf(ctx, nil, nil, "Failed to calculate discount for product: %s", json.Marshal(product))
		fmt.Fprintf(w, "Failed to calculate discount for product: %s", json.Marshal(product))
		return
	}

	fmt.Fprintf(w, "%+v", product)
}

func GetProductFromDB(ctx context.Context, id int) (*external.Product, error) {
	var result external.Product
	if id < 1 {
		return nil, errors.New("Product is invalid!")
	}

	result.Name = "product testing"
	result.Stock = rand.Int()
	return &result, nil
}

func CalculateDiscount(ctx context.Context, p *external.Product) error {
	if p.Stock%2 == 0 {
		p.Discount = 20
		log.StdInfo(ctx, p, nil, "User gets 20 discount")
	} else {
		p.Discount = 0
	}
	return nil
}

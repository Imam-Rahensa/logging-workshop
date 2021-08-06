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
)

const (
	InvalidProductIDErrorMessage string = "Product ID is Invalid"
)

func main() {

	// Set Standard Log Config
	err := log.SetStdLog(&log.Config{
		Level:     "trace",                            // Default will be in info level
		LogFile:   "./log/logging-workshop.error.log", // If none supplied will goes to os.Stderr, for production you must put log file
		DebugFile: "./log/logging-workshop.debug.log", // If none supplied will goes to os.Stderr, for production you must put log file
		AppName:   "logging-workshop",                 //  your app name, the format will be `{service_name}_{function}`
	})

	if err != nil {
		log.StdFatal(context.Background(), nil, err, "Failed to start Log")
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
		log.StdWarn(ctx, nil, nil, "[HelloHandler][Query] No Product ID supplied")
		return
	}

	// parse the product id
	if len(keys) < 1 {
		log.StdWarn(ctx, nil, nil, "[HelloHandler] No Product ID found")
		return
	}

	productID, err = strconv.Atoi(keys[0])
	if err != nil {
		log.StdWarnf(ctx, keys[0], err, "[HelloHandler][Atoi] Error when parse Product ID")
		return
	}

	// Set your context id. In this case you will put your product id
	ctx = log.SetCtxID(ctx, strconv.Itoa(productID))

	product, err := GetProductFromDB(ctx, productID)
	if err != nil {
		if err.Error() == InvalidProductIDErrorMessage {
			log.StdWarn(ctx, productID, err, "[HelloHandler][GetProductFromDB] Invalid Product ID")
		} else {
			log.StdError(ctx, productID, err, "[HelloHandler][GetProductFromDB] Error when get Product from DB")
		}

		return
	}

	err = CalculateDiscount(ctx, product)
	if err != nil {
		log.StdErrorf(ctx, product, err, "[HelloHandler][CalculateDiscount] Error on Calculating Discount")
		return
	}

	fmt.Fprintf(w, "%+v", product)
}

func GetProductFromDB(ctx context.Context, id int) (*external.Product, error) {
	var result external.Product
	if id < 1 {
		return nil, errors.New(InvalidProductIDErrorMessage)
	}

	result.Name = "product testing"
	result.Stock = rand.Int()
	return &result, nil
}

func CalculateDiscount(ctx context.Context, p *external.Product) error {
	if p.Stock%2 == 0 {
		p.Discount = 20
		log.StdDebug(ctx, p, nil, "[CalculateDiscount] User get 20 discount")
	} else {
		p.Discount = 0
	}
	return nil
}

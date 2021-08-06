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

func main() {

	// Set Standard Log Config
	err := log.SetStdLog(&log.Config{
		Level:     "trace",                            // Default will be in info level
		LogFile:   "./log/logging-workshop.error.log", // If none supplied will goes to os.Stderr, for production you must put log file
		DebugFile: "./log/logging-workshop.debug.log", // If none supplied will goes to os.Stderr, for production you must put log file
		AppName:   "logging-workshop",                 //  your app name, the format will be `{service_name}_{function}`
	})

	if err != nil {
		log.StdInfo(context.Background(), nil, err, "Failed to start Log")
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
		log.StdFatal(ctx, nil, nil, "No product id supplied")
		fmt.Fprint(w, "No product id supplied")
		return
	}

	// parse the product id
	if len(keys) < 1 {
		// Dont use fatal, because this filter just want to make sure result is not out of index (IndexError)
		log.StdError(ctx, nil, nil, "No product id found")
		fmt.Fprint(w, "No product id supplied")
		return
	}

	productID, err = strconv.Atoi(keys[0])
	if err != nil {
		// Use Warn if cannot convert Atoi because wrong input parameter
		log.StdWarnf(ctx, nil, nil, "Product id not valid %s", keys[0])
		fmt.Fprint(w, "No product id supplied")
		return
	}

	// Set your context id. In this case you will put your product id
	ctx = log.SetCtxID(ctx, strconv.Itoa(productID))

	product, err := GetProductFromDB(ctx, productID)
	if err != nil {
		// Log error trace from GetProductFromDB
		log.StdError(ctx, nil, err, "Invalid Product Id")
		fmt.Fprint(w, "Invalid id")
		return
	}

	_ = CalculateDiscount(ctx, product)
	// No need err because return always nil(?)
	// if err != nil {
	// 	fmt.Fprint(w, "Invalid id")
	// 	return
	// }

	fmt.Fprintf(w, "%+v", product)
}

func GetProductFromDB(ctx context.Context, id int) (*external.Product, error) {
	var result external.Product
	if id < 1 {
		// errors.New() will return warning if string in capitalized
		return nil, errors.New("product id invalid")
	}

	result.Name = "product testing"
	result.Stock = rand.Int()
	return &result, nil
}

func CalculateDiscount(ctx context.Context, p *external.Product) error {
	if p.Stock%2 == 0 {
		p.Discount = 20
		// No need StdError because want to inform user get 20 discount
		log.StdInfo(ctx, p, nil, "User get 20 discount")
	} else {
		p.Discount = 0
	}
	return nil
}

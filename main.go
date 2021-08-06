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
		log.StdError(ctx, nil, errors.New("No product id supplied"), "[HelloHandler] No product id supplied")
		return
	}

	// parse the product id
	if len(keys) < 1 {
		log.StdError(ctx, nil, errors.New("No product id found"), "[HelloHandler] No product id found")
		return
	}

	productID, err = strconv.Atoi(keys[0])
	if err != nil {
		log.StdError(ctx, nil, err, fmt.Sprintf("[HelloHandler] Product id is not valid. value: %s", keys[0]))
		return
	}

	// Set your context id. In this case you will put your product id
	ctx = log.SetCtxID(ctx, strconv.Itoa(productID))

	product, err := GetProductFromDB(ctx, productID)
	if err != nil {
		return
	}

	err = CalculateDiscount(ctx, product)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "%+v", product)
}

func GetProductFromDB(ctx context.Context, id int) (*external.Product, error) {
	var result external.Product
	if id < 1 {
		log.StdError(ctx, nil, errors.New("Product id Invalid"), fmt.Sprintf("[GetProductFromDB] ProductID is Invalid"))
		return nil, errors.New("Product id Invalid")
	}

	result.Name = "product testing"
	result.Stock = rand.Int()
	return &result, nil
}

func CalculateDiscount(ctx context.Context, p *external.Product) error {
	if p.Stock%2 == 0 {
		p.Discount = 20
		log.StdDebug(ctx, p, nil, fmt.Sprintf("[CalculateDiscount] User get %d discount", p.Discount))
	} else {
		p.Discount = 0
	}

	return nil
}

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
		log.StdError(context.Background(), nil, err, "Failed to start Log")
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
		log.StdError(ctx, keys, errors.New("ProductID Not Found in Query Param"), "[Parsing ProductID Error]")
		fmt.Fprint(w, "No product id supplied")
		return
	}

	// parse the product id
	if len(keys) < 1 {
		log.StdError(ctx, keys, errors.New("ProductID Not Found"), "[Parsing ProductID Error]")
		fmt.Fprint(w, "No product id supplied")
		return
	}

	productID, err = strconv.Atoi(keys[0])
	if err != nil {
		log.StdError(ctx, keys[0], err, "[Parsing ProductID Error]")
		fmt.Fprint(w, "No product id supplied")
		return
	}

	// Set your context id. In this case you will put your product id
	ctx = log.SetCtxID(ctx, strconv.Itoa(productID))

	product, err := GetProductFromDB(ctx, productID)
	if err != nil {
		log.StdError(ctx, product, err, "[GetProductFromDB Error]")
		fmt.Fprint(w, "Invalid id")
		return
	}

	err = CalculateDiscount(ctx, product)
	if err != nil {
		log.StdError(ctx, product, err, "[CalculateDiscount Error]")
		fmt.Fprint(w, "Invalid id")
		return
	}

	fmt.Fprintf(w, "%+v", product)
}

func GetProductFromDB(ctx context.Context, id int) (*external.Product, error) {
	var result external.Product
	if id < 1 {
		return nil, errors.New("Product id Invalid")
	}

	result.Name = "product testing"
	result.Stock = rand.Int()
	return &result, nil
}

func CalculateDiscount(ctx context.Context, p *external.Product) error {
	if p.Stock%2 == 0 {
		p.Discount = 20
	} else {
		p.Discount = 0
	}
	return nil
}

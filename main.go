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
		log.StdWarn(ctx, nil, nil, "Keys not found")
		fmt.Fprint(w, "Keys not found")
		return
	}

	// parse the product id
	if len(keys) < 1 {
		log.StdWarn(ctx, nil, nil, "Keys length is less than 1")
		fmt.Fprint(w, "Keys length is less than 1")
		return
	}

	productID, err = strconv.Atoi(keys[0])
	if err != nil {
		log.StdWarnf(ctx, nil, err, "Key is not integer. Key: %s", keys[0])
		fmt.Fprint(w, "Key is not integer")
		return
	}

	// Set your context id. In this case you will put your product id
	ctx = log.SetCtxID(ctx, strconv.Itoa(productID))

	product, err := GetProductFromDB(ctx, productID)
	if err != nil {
		log.StdErrorf(ctx, nil, err, "Error get product from DB. ProductID: %d", productID)
		fmt.Fprint(w, "Error get product from DB")
		return
	}

	err = CalculateDiscount(ctx, product)
	if err != nil {
		log.StdError(ctx, product, err, "Error calculate discount")
		fmt.Fprint(w, "Error calculate discount")
		return
	}

	fmt.Fprintf(w, "%+v", product)
}

func GetProductFromDB(ctx context.Context, id int) (*external.Product, error) {
	var result external.Product
	if id < 1 {
		return nil, errors.New("product id invalid")
	}

	result.Name = "product testing"
	result.Stock = rand.Int()
	return &result, nil
}

func CalculateDiscount(ctx context.Context, p *external.Product) error {
	if p.Stock%2 == 0 {
		p.Discount = 20
		log.StdTrace(ctx, p, nil, "User get 20 discount")
	} else {
		p.Discount = 0
	}
	return nil
}

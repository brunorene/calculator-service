package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/brunorene/calculator-service/operator"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ErrLogger  = errors.New("logger error")
	ErrNotInt  = errors.New("value not int")
	ErrMissing = errors.New("needs 2 values")
)

func main() {
	args := os.Args

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var logLevel string

	flags.StringVar(&logLevel, "log-level", "info",
		"Log level used for logging. Should be one of: debug, info, warn, error.")

	if err := setupLogger(logLevel); err != nil {
		log.Fatalf("setup logger: %v", err)
	}

	if err := runServer(); err != nil {
		zap.S().Errorf("run server: %v", err)
	}
}

func runServer() error {
	http.HandleFunc("/add/", operatorHandler(&operator.Add{}))
	http.HandleFunc("/sub/", operatorHandler(&operator.Subtract{}))
	http.HandleFunc("/mul/", operatorHandler(&operator.Multiply{}))
	http.HandleFunc("/div/", operatorHandler(&operator.Divide{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(404) })

	if err := http.ListenAndServe(":8090", nil); err != nil {
		return fmt.Errorf("listen & serve on 8090: %v", err)
	}

	return nil
}

type response struct {
	Result  *int   `json:"result,omitempty"`
	IsError bool   `json:"error"`
	Message string `json:"message,omitempty"`
}

func values(parts []string) (left, right int, err error) {
	if len(parts) < 4 {
		return 0, 0, fmt.Errorf("%w: %v", ErrMissing, parts)
	}

	left, err = strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: left %s: %v", ErrNotInt, parts[2], err)
	}

	right, err = strconv.Atoi(parts[3])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: right %s: %v", ErrNotInt, parts[3], err)
	}

	return
}

func operatorHandler(op operator.Operator) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, req *http.Request) {
		parts := strings.Split(req.URL.Path, "/")

		left, right, err := values(parts)
		if err != nil {
			content, err := json.Marshal(response{
				IsError: true,
				Message: err.Error(),
			})
			if err != nil {
				writer.WriteHeader(500)
				zap.S().Errorf("%v: 500", err)

				return
			}

			writer.WriteHeader(400)
			zap.S().Errorf("%v: 400", err)

			writer.Header().Add("Content-Type", "application/json")
			if _, err := writer.Write(content); err != nil {
				zap.S().Errorf("write error content: %v", err)
			}

			return
		}

		result := op.Result(left, right)

		content, err := json.Marshal(response{Result: &result})
		if err != nil {
			writer.WriteHeader(500)
			zap.S().Errorf("%v: 500", err)

			return
		}

		writer.Header().Add("Content-Type", "application/json")
		if _, err := writer.Write(content); err != nil {
			zap.S().Errorf("write ok content: %v", err)
		}
	}
}

func setupLogger(levelName string) error {
	config := zap.NewProductionConfig()

	level, err := zap.ParseAtomicLevel(levelName)
	if err != nil {
		return fmt.Errorf("parse level: %w: %v", ErrLogger, err)
	}

	config.Level.SetLevel(level.Level())

	config.EncoderConfig.TimeKey = "@timestamp"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	log, err := config.Build()
	if err != nil {
		return fmt.Errorf("config build: %w: %v", ErrLogger, err)
	}

	// nolint:errcheck // deferred and not important if fails
	defer log.Sync() // flushes buffer, if any

	zap.ReplaceGlobals(log)

	return nil
}

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// statsFunc represents a generic statistical function.
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	result := 0.0

	for _, v := range data {
		result += v
	}

	return result
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func csv2Float(r io.Reader, column int) ([]float64, error) {
	// Create the CSV reader to take data from csv files.
	cr := csv.NewReader(r)

	// Adjust for a 0-based index.
	column--

	// Read in all CSV data.
	allData, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}

	var data []float64
	// Loop through all records, skipping title rows.
	for i, row := range allData {
		if i == 0 {
			continue
		}

		// Check number of columns in the CSV file compared to user request..
		if len(row) <= column {
			return nil, fmt.Errorf("%w: file has only %d columns", ErrInvalidColumn, len(row))
		}

		// Try to convert the data to a float64.
		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}

	return data, nil
}

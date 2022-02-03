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

func min(data []float64) float64 {
	for _, v := range data {
		isMin := true

		for _, v2 := range data {
			if v > v2 {
				isMin = false
				break
			}
		}

		if isMin {
			return v
		}
	}

	return 0.0
}

func max(data []float64) float64 {
	for _, v := range data {
		isMax := true

		for _, v2 := range data {
			if v < v2 {
				isMax = false
				break
			}
		}

		if isMax {
			return v
		}
	}

	return 0.0
}

func csv2Float(r io.Reader, column int) ([]float64, error) {
	// Create the CSV reader to take data from csv files.
	cr := csv.NewReader(r)
	// Reuse the record to save memory, as we're using Read() inside the loop below.
	cr.ReuseRecord = true

	// Adjust for a 0-based index.
	column--

	var data []float64
	// Loop through all records.
	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from file: %w", err)
		}
		// Skip title rows.
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

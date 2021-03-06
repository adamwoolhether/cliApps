// Package todo defines the API for our todo program.
package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item represents a todo item.
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represents a list of todo items.
type List []item

// String implements the fmt.Stringer interface to print out a formatted list.
func (l *List) String() string {
	formatted := ""

	for k, t := range *l {
		prefix := "  "
		if t.Done {
			prefix = "X "
		}

		formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
	}

	return formatted
}

/*func (l *List) Format(state fmt.State, verb rune) {
	typ := reflect.TypeOf(item{})
	itemFields := make([]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		itemFields[i] = typ.Field(i).Name
	}

	switch verb {
	case 's', 'q':
		val := l.String()
		if verb == 'q' {
			val = fmt.Sprintf("%q", val)
		}
		fmt.Fprint(state, val)
	case 'v':
		if state.Flag('#') {
			fmt.Fprint(state, "%T", l)
		}
		fmt.Fprint(state, "{")
		// val := reflect.ValueOf(*l)
		for i, name := range itemFields {
			if state.Flag('#') || state.Flag('+') {
				fmt.Fprintf(state, "%s:", name)
			}
			// fld := val.FieldByName(name)
			// if name == "Task" && fld.Len() > 0 {
			// 	fmt.Fprint(state, keyMask)
			// } else {
			// fmt.Fprint(state, fld)
			// }
			if i < len(itemFields)-1 {
				fmt.Fprint(state, " ")
			}
		}
		fmt.Fprint(state, "}")
	}
}*/

// Add creates a new todo item and appends it to the list.
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*l = append(*l, t)
}

// Complete marks a todo item as completed
// by setting Done = true and  CompletedAt to the current time.
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	// Adjust for a 0-based index.
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete removes a todo item from the list.
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	// Adjust for 0 based index.
	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save encodes the list as JSON and saves
// it using the provided file name.
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

// Get opens the provided file, decodes
// the JSON and parses it into a list.
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, l)
}

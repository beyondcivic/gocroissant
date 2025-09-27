package croissant

import "fmt"

type CroissantError struct {
	// Message to show the user.
	Message string
	// Value to include with message
	Value any
}

func (e CroissantError) Error() string {
	if e.Value != nil {

		return fmt.Sprintf("%s: %v", e.Message, e.Value)
	} else {
		return fmt.Sprintf("%s", e.Message)
	}
}

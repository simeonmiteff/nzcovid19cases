package nzcovid19cases

import "fmt"

type InvalidUsageError struct{
	Problem string
}

func (e InvalidUsageError) Error() string {
	return fmt.Sprintf("invalid usage: %v", e.Problem)
}

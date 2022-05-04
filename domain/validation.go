package domain

import "context"

type ValidationService interface {
	ValidateActivity(ctx context.Context, activity string) bool
	CleanString(str string) (newStr string)
}

package errors

const (
	EDuplicateEntry = 1062
)

func Is(expected int, actual int) bool {
	return expected == actual
}

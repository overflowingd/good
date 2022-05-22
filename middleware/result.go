package middleware

type Processing struct {
	Result any
}

type Created struct {
	Result any
}

func NewProcessing(result any) *Processing {
	return &Processing{
		Result: result,
	}
}

func NewCreated(result any) *Created {
	return &Created{
		Result: result,
	}
}

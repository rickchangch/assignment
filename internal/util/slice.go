package util

// Convert []A -> []B
func Map[In, Out any](in []In, fn func(In) Out) []Out {
	if in == nil {
		return nil
	}

	out := make([]Out, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = fn(in[i])
	}

	return out
}

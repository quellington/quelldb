package utils

func XorMask(input []byte, key byte) []byte {
	output := make([]byte, len(input))
	for i, b := range input {
		output[i] = b ^ key
	}
	return output
}

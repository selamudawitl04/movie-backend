package utilService

import (
	"math/rand"
	"time"
)

func PublicID()  string {
    // Initialize the random number generator with a seed value based on the current time
	rand.Seed(time.Now().UnixNano())

	// Define the length of the random word to generate
	length := 18

	// Define a string containing all possible characters for the random word
	chars := "abcdefghijklmnopqrstuvwxyz_0123456789-ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Generate a random word of the specified length
	word := ""
	for i := 0; i < length; i++ {
		index := rand.Intn(len(chars))
		word += string(chars[index])
	}
	// return the random word
	return word
}
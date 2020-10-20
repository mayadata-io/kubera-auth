package random

import "math/rand"

//GetRandomNumbers generates random strings, can be used to create ids or random secrets
func GetRandomNumbers(n int) string {
	var letters = []rune("0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

//GetRandomString generates random strings, can be used to create ids or random secrets
func GetRandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz!@#$%^&*ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

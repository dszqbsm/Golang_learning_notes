package main

func reverseWithGenerics[T any](s []T) []T {
	l := len(s)
	r := make([]T, l)
	for i, e := range s {
		r[l-i-1] = e
	}
	return r
}

func main() {

}

package util

import (
    "fmt"
    "strings"
)

func Zip(slices ...[]int) ([][]int, error) {
    if len(slices) == 0 {
        return [][]int{}, nil
    }

    length := len(slices[0])
    for i := 1; i < len(slices); i++ {
        if len(slices[i]) != length {
            return nil, fmt.Errorf("zip: arguments must be of same length")
        }
    }

    r := make([][]int, length)

    for i, e := range slices[0] {
        a := make([]int, len(slices))
        a[0] = e
        for j := 1; j < len(slices); j++ {
            a[j] = slices[j][i]
        }
        r[i] = a
    }

    return r, nil
}

func CaseInsensitiveIntersection(a []string, b []string, which bool) []string {
    set := make([]string, 0)
    hash := make(map[string]bool)

    for i := 0; i < a.Len(); i++ {
        hash[strings.ToLower(a[i])] = true
    }

    for i := 0; i < b.Len(); i++ {
        if _, found := hash[strings.ToLower(b[i])]; found {
            val := a[i]
            if which {
                val = b[i]
            }
            set = append(set, val)
        }
    }

    return set
}

func Contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

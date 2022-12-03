package util

import (
    "math/rand"
    "strings"
)

const alphanet = "abcdefghijqlmnopqrstuvwxyz"

func RandomInt(min, max int64) int64 {
    return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
    var sb strings.Builder

    k := len(alphanet)

    for i := 0; i < n; i++ {
        c := alphanet[rand.Intn(k)]
        sb.WriteByte(c)
    }

    return sb.String()
}

func RandomOwner() string {
    return RandomString(5)
}

func RandomMoney() int64 {
    return RandomInt(1, 1000)
}

func RandomCurrency() string {
    currencies := []string{BRL, USD, EUR}

    a := len(currencies)
    return currencies[rand.Intn(a)]
}

package util

import (
	"math/rand"
	"strings"
	"time"
)



func RandInt(min ,max int)int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min + 1)+min
}

func RandomString(length int)string{
	var sb strings.Builder
	//b := make([]byte,length)
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i :=0;i < length;i++{
		rand.Seed(time.Now().UnixNano())
		c := letterBytes[rand.Intn(len(letterBytes))]
		sb.WriteByte(c)

	}
	return sb.String()
}


func RandomOwner()string{
	return RandomString(6)
}
func RandomMoney()int64{
	return int64(RandInt(0,1000))
}
func RandomCurrency()string{
	currencies := []string{"INR","USD","SGP"}
	rand.Seed(time.Now().UnixNano())
	return currencies[rand.Intn(len(currencies))]
}
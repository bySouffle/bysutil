package num2cn

import (
	"fmt"
	"testing"
)

func TestEncodeFromFloat64(t *testing.T) {
	chNum := "负十七亿零五十三万七千零一十六"
	num, _ := CN2Int64(chNum)
	fmt.Println(num) // -1700537016
	chNumAgain := Int64ToCN(num)
	fmt.Println(chNumAgain) // 负十七亿零五十三万七千零一十六

	chFloatNum := "负零点零七三零六"
	fNum, _ := CN2Float64(chFloatNum)
	fmt.Printf("%f\n", fNum) // -0.073060
	chFloatNumAgain := Float64ToCN(fNum)
	fmt.Println(chFloatNumAgain) // 负零点零七三零六

	fNumber := 1.01
	chfNum := Float64ToCN(fNumber)
	fmt.Println(chfNum)
}

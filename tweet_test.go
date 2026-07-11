package bacot

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	bacot "github.com/yumiodd/bacot/src"
)

type Tweet struct {
	Text string
	Bad  string
}

func TestTweet(t *testing.T) {

	b := bacot.New()

	file, err := os.Open("data_tweet.csv")
	if err != nil {
		t.Fatal("cant open file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// reader.Comma = '\t'

	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		t.Fatal("cant read all")
	}
	// var (
	// 	pass int
	// 	FN   int
	// 	FP   int
	// )
	for _, d := range data[1:] {

		res := b.Text(d[0]).Scan()
		// isBad := d[0] == "HS"

		// if res.IsProfane() == isBad {
		// 	pass += 1
		// } else if res.IsProfane() && isBad == false {

		// 	FP += 1
		// } else if !res.IsProfane() && isBad == true {

		// 	FN += 1
		// }

		if res.IsProfane() {
			fmt.Printf("TEXT: (%s) %s\n", res.First(), d[0])
		}
	}

}

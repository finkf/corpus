package corpus

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

var update bool

func init() {
	flag.BoolVar(&update, "update", false, "update gold file(s)")
}

func dtaToString(dta DTA) string {
	var str string
	err := dta.Tokenize(func(token string) {
		str += token + "\n"
	})
	if err != nil {
		panic(err)
	}
	return str
}

func updateDTAGoldFile() {
	dta, err := NewDTAFile("testdata/dta.xml")
	if err != nil {
		panic(err)
	}
	defer func() {
		if e := dta.Close(); e != nil {
			panic(e)
		}
	}()
	gold := dtaToString(dta)
	if err = ioutil.WriteFile("testdata/dta.gold.txt", []byte(gold), os.ModePerm); err != nil {
		panic(err)
	}
}

func TestDTATokenizeGoldFile(t *testing.T) {
	if update {
		updateDTAGoldFile()
	}
	gold, err := ioutil.ReadFile("testdata/dta.gold.txt")
	if err != nil {
		panic(err)
	}
	dta, err := NewDTAFile("testdata/dta.xml")
	if err != nil {
		panic(err)
	}
	if got := dtaToString(dta); string(gold) != got {
		t.Errorf("expected\n%sgot\n%s", gold, got)
	}
}

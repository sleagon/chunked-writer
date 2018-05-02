package writer

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	// default* files should be cleaned
	day := time.Now().Format(DefaultPattern)
	removeAll(fmt.Sprintf("./%s.%s*", DefaultPrefix, day))
	w, _ := New("", "", 1<<2)
	for k := 0; k < 32; k++ {
		w.Write([]byte("1"))
	}
	f, err := os.Open(fmt.Sprintf("%s.%s.log.7", DefaultPrefix, day))
	if err != nil {
		t.Fail()
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fail()
	}
	if string(data) != "1111" {
		t.Fail()
	}
	_, err = os.Stat(fmt.Sprintf("%s.%s.log.8", DefaultPrefix, day))
	if !os.IsNotExist(err) {
		t.Fail()
	}
	w, _ = New("", "", 1<<2)
	for k := 0; k < 32; k++ {
		w.Write([]byte("1"))
	}
	f, err = os.Open(fmt.Sprintf("%s.%s.log.15", DefaultPrefix, day))
	if err != nil {
		t.Fail()
	}
	data, err = ioutil.ReadAll(f)
	if err != nil {
		t.Fail()
	}
	if string(data) != "1111" {
		t.Fail()
	}
	_, err = os.Stat(fmt.Sprintf("%s.%s.log.16", DefaultPrefix, day))
	if !os.IsNotExist(err) {
		t.Fail()
	}
	removeAll(fmt.Sprintf("./%s.%s*", DefaultPrefix, day))
}

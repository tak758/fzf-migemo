package migemo

import (
	"os"
	"regexp"
	"sync"
	"unicode/utf8"

	"github.com/koron/gomigemo/embedict"
	"github.com/koron/gomigemo/migemo"
)

var cache map[string]*regexp.Regexp
var dict migemo.Dict
var mutex *sync.Mutex

func init() {
	mutex = new(sync.Mutex)
	cache = make(map[string]*regexp.Regexp)
	var err error
	dict, err = embedict.Load()
	if err != nil {
		errorExit(err.Error())
	}
}

func errorExit(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(2)
}

func regex(pattern string) *regexp.Regexp {
	mutex.Lock()
	defer mutex.Unlock()

	v, ok := cache[pattern]
	if ok {
		return v
	}

	re, err := migemo.Compile(dict, pattern)
	if err != nil {
		errorExit(err.Error())
	}

	cache[pattern] = re
	return re
}

func FindStringIndex(s, pattern string) []int {
	v := regex(pattern)

	if i := v.FindStringIndex(s); len(i) != 0 {
		b := []byte(s)
		subString := string(b[i[0]:i[1]])
		return []int{utf8.RuneCountInString(s[0:i[0]]), utf8.RuneCountInString(subString)}
	}

	return nil
}

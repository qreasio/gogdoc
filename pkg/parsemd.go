package gdoc

import (
	"bytes"
	"github.com/yuin/goldmark"
	"io/ioutil"
)

func getMarkdown(filepath string) []byte {
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return dat
}

func ParseMD(filepath string) string {
	var buf bytes.Buffer
	if err := goldmark.Convert(getMarkdown(filepath), &buf); err != nil {
		panic(err)
	}

	return buf.String()
}


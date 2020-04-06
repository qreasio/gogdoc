package gdoc

import (
	"google.golang.org/api/docs/v1"
	"strings"
)

type Document struct {
	GDoc docs.Document
	Taxonomies *Taxonomies
}

type Element struct {
	StructuralElement *docs.StructuralElement
	StringBuilder     strings.Builder
}

type Taxonomies struct {
	Tags *[]string
	Categories *[]string
}
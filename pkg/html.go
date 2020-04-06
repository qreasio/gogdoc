package gdoc

import (
	"fmt"
	"google.golang.org/api/docs/v1"
	"strings"
)

var (
	NamedStyleTag = map[string]string{
		"NORMAL_TEXT": "p",
		"HEADING_1":   "h1",
		"HEADING_2":   "h2",
		"HEADING_3":   "h3",
		"HEADING_4":   "h4",
		"HEADING_5":   "h5",
	}
)

func (d *Document) String() string {
	var str strings.Builder
	for _, el := range d.GDoc.Body.Content {
		elem := Element{StructuralElement: el}
		str.WriteString(elem.String())
	}
	return str.String()
}


func (d *Document) HTML() string {
	var str strings.Builder
	for _, el := range d.GDoc.Body.Content {
		elem := Element{StructuralElement: el}
		str.WriteString(elem.HTML())
	}
	return str.String()
}

func GetOpenTag(p *docs.Paragraph) string {
	tag, ok := NamedStyleTag[p.ParagraphStyle.NamedStyleType]
	if ok {
		return fmt.Sprintf("<%s>", tag)
	}
	return fmt.Sprintf("<span>")
}

func GetClosingTag(p *docs.Paragraph) string {
	tag, ok := NamedStyleTag[p.ParagraphStyle.NamedStyleType]
	if ok {
		return fmt.Sprintf("</%s>", tag)
	}
	return fmt.Sprintf("</span>")
}

func (el *Element) String() string {

	if el.StructuralElement.Paragraph != nil {
		for _, p := range el.StructuralElement.Paragraph.Elements {
			el.StringBuilder.WriteString(p.TextRun.Content)
		}
	}

	return el.StringBuilder.String()
}

func (el *Element) HTML() string {

	if el.StructuralElement.Paragraph != nil {

		openTag := GetOpenTag(el.StructuralElement.Paragraph)
		var content strings.Builder
		for _, el := range el.StructuralElement.Paragraph.Elements {
			content.WriteString(el.TextRun.Content)
		}
		closingTag := GetClosingTag(el.StructuralElement.Paragraph)

		el.StringBuilder.WriteString(fmt.Sprintf("%s%s%s", openTag, content.String(), closingTag))
	}

	return el.StringBuilder.String()
}

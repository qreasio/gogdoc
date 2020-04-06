package gdoc

import (
	"strings"
)

func (d *Document) ParseTaxonomies() *Taxonomies {
	if d.GDoc.DocumentStyle == nil {
		return nil
	}
	var headerContent strings.Builder
	for _, c := range d.GDoc.Headers[d.GDoc.DocumentStyle.DefaultHeaderId].Content {
		for _, el := range c.Paragraph.Elements {
			if el.TextRun != nil {
				headerContent.WriteString(el.TextRun.Content)
			}
		}
	}
	var rows []string
	rows = strings.Split(headerContent.String(), "\n")
	d.Taxonomies = &Taxonomies{}
	for _, row := range rows {
		rowElements := strings.Split(row, ":")
		if strings.ToLower(strings.TrimRight(rowElements[0], " ")) == "tag" || strings.ToLower(strings.TrimRight(rowElements[0], " ")) == "tags" {
			tags := strings.Split(rowElements[1], ",")
			d.Taxonomies.Tags = &tags
		}
		if strings.ToLower(strings.TrimRight(rowElements[0]," ")) == "category" || strings.ToLower(strings.TrimRight(rowElements[0]," ")) == "categories" {
			categories := strings.Split(rowElements[1], ",")
			d.Taxonomies.Categories = &categories
		}
	}
	return d.Taxonomies
}
package gdoc

import (
	"fmt"
	"google.golang.org/api/docs/v1"
	"strings"
)

//reference
//https://learn.netlify.com/en/cont/markdown/

var (
	markdownStyles = map[string]string{
		"NORMAL_TEXT":     "",
		"HEADING_1":       "# ",
		"HEADING_2":       "## ",
		"HEADING_3":       "### ",
		"HEADING_4":       "#### ",
		"HEADING_5":       "##### ",
		"horizontal_rule": "---",
		"blockquote":      "> ",
		"unordered":       "- ",
		"ordered":         "1. ",
	}

	markdownTextStyle = map[string]string{
		"bold":          "**",
		"italic":        "_",
		"strikethrough": "~~",
		"link":          "[%s](%s)",
	}
)

func GetParagraphMarkdownSign(doc docs.Document, p *docs.Paragraph) string {
	orderedList := true
	markdownSign := ""

	if p.Bullet != nil {
		listId := p.Bullet.ListId
		for i, level := range doc.Lists[listId].ListProperties.NestingLevels {
			if i == 0 && level.GlyphSymbol == "" {
				orderedList = false
			}
		}

		if orderedList {
			markdownSign = markdownStyles["ordered"]
		} else {
			markdownSign = markdownStyles["unordered"]
		}

		return markdownSign
	}

	style, ok := markdownStyles[p.ParagraphStyle.NamedStyleType]
	if ok {
		markdownSign = style
		if p.ParagraphStyle.IndentFirstLine != nil && p.ParagraphStyle.IndentFirstLine.Magnitude == p.ParagraphStyle.IndentStart.Magnitude {
			markdownSign = markdownStyles["blockquote"]
		}
	}

	if strings.Contains(markdownSign, "#") {
		markdownSign = "\n" + markdownSign
	}

	return markdownSign
}

func GetMarkdownElement(el *docs.ParagraphElement) string {
	markdownSign := ""
	if el.TextRun != nil {

		if el.TextRun.TextStyle.Bold {
			markdownSign = markdownTextStyle["bold"]
		}
		if el.TextRun.TextStyle.Italic {
			markdownSign = markdownTextStyle["italic"]
		}
		if el.TextRun.TextStyle.Link != nil && el.TextRun.Content != "" {
			markdownSign = markdownTextStyle["link"]
		}

	}

	if el.HorizontalRule != nil {
		return markdownStyles["horizontal_rule"]
	}

	return markdownSign
}

func (d *Document) Markdown() string {
	var str strings.Builder
	for _, el := range d.GDoc.Body.Content {
		elem := Element{StructuralElement: el}
		str.WriteString(elem.Markdown(d.GDoc))
	}
	return str.String()
}

func (el *Element) Markdown(doc docs.Document) string {
	if el.StructuralElement.Paragraph != nil {
		el.FormatParagraph(doc)
	}
	if el.StructuralElement.Table != nil {
		el.FormatTable()
	}
	return el.StringBuilder.String()
}

func (el *Element) FormatTable() {
	//if it is table with 1 row and 1 column (1 cell only) we format using markdown code
	if  el.StructuralElement.Table.Columns == int64(1) &&  el.StructuralElement.Table.Rows == int64(1) {
		el.StringBuilder.WriteString(el.formatAsCode())
	} else {
		el.StringBuilder.WriteString(el.formatAsTable())
	}
}

func (el *Element) FormatParagraph(doc docs.Document) {
	var content strings.Builder
	sign := GetParagraphMarkdownSign(doc, el.StructuralElement.Paragraph)
	content.WriteString(sign)

	for _, e := range el.StructuralElement.Paragraph.Elements {
		content.WriteString(FormatParagraphElement(doc, e, &content))
	}

	el.StringBuilder.WriteString(content.String())
}

func FormatParagraphElement(doc docs.Document, el *docs.ParagraphElement, content *strings.Builder) string {
	markdownSign := GetMarkdownElement(el)

	if el.TextRun != nil {
		if el.TextRun.TextStyle.Link == nil {
			return fmt.Sprintf("%s%s%s", markdownSign, el.TextRun.Content, markdownSign)
		} else {
			return fmt.Sprintf(markdownSign, el.TextRun.Content, el.TextRun.TextStyle.Link.Url)
		}
	}

	if el.HorizontalRule != nil {
		return markdownSign
	}

	if el.InlineObjectElement != nil && doc.InlineObjects[el.InlineObjectElement.InlineObjectId].InlineObjectProperties.EmbeddedObject.ImageProperties != nil {
		imageURI := doc.InlineObjects[el.InlineObjectElement.InlineObjectId].InlineObjectProperties.EmbeddedObject.ImageProperties.ContentUri
		return fmt.Sprintf("![](%s)", imageURI)
	}

	return ""
}

func (el *Element) formatAsCode() string {
	var content strings.Builder
	content.WriteString("```\n")
	for _, ct := range el.StructuralElement.Table.TableRows[0].TableCells[0].Content {
		for _, el := range ct.Paragraph.Elements {
			content.WriteString(el.TextRun.Content)
		}
	}
	content.WriteString("```")
	return content.String()
}

func (el *Element) formatAsTable() string {
	var content strings.Builder

	for i, row := range el.StructuralElement.Table.TableRows {
		rowCells := make([]string, 0)
		for _, cell := range row.TableCells {
			var cellContent strings.Builder
			for _, ct := range cell.Content {
				for _, e := range ct.Paragraph.Elements {
					//fmt.Printf("%s\n", e.TextRun.Content)

					cellContent.WriteString(strings.Replace(e.TextRun.Content, "\n", "", -1))
				}
			}
			rowCells = append(rowCells, cellContent.String())
		}

		rowContent := strings.Join(rowCells, " | ")
		content.WriteString(fmt.Sprintf("| %s |\n", rowContent))

		if i == 0 {
			rowHeader := make([]string, 0)
			for i := 0; i <= len(rowCells); i++ {
				rowHeader = append(rowHeader, "|")
			}
			rowHeaderContent := strings.Join(rowHeader, "--")
			content.WriteString(rowHeaderContent + "\n")
		}
	}
	return content.String()
}

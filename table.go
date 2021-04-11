package main

import (
	"strings"
)

type HtmlTable struct {
	header strings.Builder
	rows   strings.Builder
}

func (t *HtmlTable) setHeaders(headers []string) {
	t.header.WriteString("<thead><tr>")
	t.header.WriteString("<th>" + headers[0] + "</th>")

	for _, title := range headers[1:] {
		t.header.WriteString("<th align=\"right\">" + title + "</th>")
	}
	t.header.WriteString("</tr></thead>")
}

func (t *HtmlTable) addRow(row []string) {
	t.rows.WriteString("<tr>")
	t.rows.WriteString("<td>" + row[0] + "</td>")

	for _, cell := range row[1:] {
		t.rows.WriteString("<td align=\"right\">" + cell + "</td>")
	}
	t.rows.WriteString("</tr>")
}

func (t *HtmlTable) render() string {
	var renderedTable strings.Builder
	renderedTable.WriteString("<table>")
	renderedTable.WriteString(t.header.String())
	renderedTable.WriteString("<tbody>")
	renderedTable.WriteString(t.rows.String())
	renderedTable.WriteString("</tbody></table>")
	return renderedTable.String()
}
package semanticmd_test

import (
	"strings"
	"testing"

	semanticmd "github.com/thorstenpfister/semantic-markdown"
)

func TestSimpleTable(t *testing.T) {
	htmlStr := `
	<table>
		<tr>
			<td>Cell 1</td>
			<td>Cell 2</td>
		</tr>
		<tr>
			<td>Cell 3</td>
			<td>Cell 4</td>
		</tr>
	</table>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	expected := `| Cell 1 | Cell 2 |
| Cell 3 | Cell 4 |`

	if !strings.Contains(result, expected) {
		t.Errorf("Expected table output not found.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestTableWithHeader(t *testing.T) {
	htmlStr := `
	<table>
		<tr>
			<th>Header 1</th>
			<th>Header 2</th>
		</tr>
		<tr>
			<td>Cell 1</td>
			<td>Cell 2</td>
		</tr>
	</table>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have separator row after header
	expected := `| Header 1 | Header 2 |
| --- | --- |
| Cell 1 | Cell 2 |`

	if !strings.Contains(result, expected) {
		t.Errorf("Expected table with header not found.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestTableWithTheadTbody(t *testing.T) {
	htmlStr := `
	<table>
		<thead>
			<tr>
				<th>Name</th>
				<th>Age</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>Alice</td>
				<td>30</td>
			</tr>
			<tr>
				<td>Bob</td>
				<td>25</td>
			</tr>
		</tbody>
	</table>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	expected := `| Name | Age |
| --- | --- |
| Alice | 30 |
| Bob | 25 |`

	if !strings.Contains(result, expected) {
		t.Errorf("Expected table with thead/tbody not found.\nExpected:\n%s\n\nGot:\n%s", expected, result)
	}
}

func TestTableWithColspan(t *testing.T) {
	htmlStr := `
	<table>
		<tr>
			<td colspan="2">Wide Cell</td>
		</tr>
		<tr>
			<td>Cell 1</td>
			<td>Cell 2</td>
		</tr>
	</table>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have colspan comment
	if !strings.Contains(result, "colspan: 2") {
		t.Errorf("Expected colspan comment not found in output:\n%s", result)
	}

	if !strings.Contains(result, "Wide Cell") {
		t.Errorf("Expected cell content not found in output:\n%s", result)
	}
}

func TestTableWithRowspan(t *testing.T) {
	htmlStr := `
	<table>
		<tr>
			<td rowspan="2">Tall Cell</td>
			<td>Cell 1</td>
		</tr>
		<tr>
			<td>Cell 2</td>
		</tr>
	</table>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have rowspan comment
	if !strings.Contains(result, "rowspan: 2") {
		t.Errorf("Expected rowspan comment not found in output:\n%s", result)
	}

	if !strings.Contains(result, "Tall Cell") {
		t.Errorf("Expected cell content not found in output:\n%s", result)
	}
}

func TestTableWithColumnTracking(t *testing.T) {
	htmlStr := `
	<table>
		<tr>
			<th>Column A</th>
			<th>Column B</th>
			<th>Column C</th>
		</tr>
		<tr>
			<td>Data 1</td>
			<td>Data 2</td>
			<td>Data 3</td>
		</tr>
	</table>
	`

	opts := &semanticmd.ConversionOptions{
		EnableTableColumnTracking: true,
	}

	result, err := semanticmd.ConvertString(htmlStr, opts)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Should have column IDs (A, B, C)
	if !strings.Contains(result, "<!-- A -->") {
		t.Errorf("Expected column ID A not found in output:\n%s", result)
	}
	if !strings.Contains(result, "<!-- B -->") {
		t.Errorf("Expected column ID B not found in output:\n%s", result)
	}
	if !strings.Contains(result, "<!-- C -->") {
		t.Errorf("Expected column ID C not found in output:\n%s", result)
	}
}

func TestComplexTableDocument(t *testing.T) {
	htmlStr := `
	<article>
		<h1>Data Report</h1>
		<section>
			<h2>Q1 Results</h2>
			<table>
				<thead>
					<tr>
						<th>Product</th>
						<th>Revenue</th>
						<th>Growth</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>Widget A</td>
						<td>$1,000</td>
						<td>+10%</td>
					</tr>
					<tr>
						<td>Widget B</td>
						<td>$2,000</td>
						<td>+20%</td>
					</tr>
				</tbody>
			</table>
		</section>
		<footer>
			<p>Report generated 2024</p>
		</footer>
	</article>
	`

	result, err := semanticmd.ConvertString(htmlStr, nil)
	if err != nil {
		t.Fatalf("ConvertString failed: %v", err)
	}

	// Check article renders without wrapper
	if strings.Contains(result, "<!-- <article>") {
		t.Error("Article should not have HTML comment wrapper")
	}

	// Check section has horizontal rules
	if !strings.Contains(result, "---") {
		t.Error("Section should have horizontal rules")
	}

	// Check table has header separator
	if !strings.Contains(result, "| --- | --- | --- |") {
		t.Error("Table should have header separator")
	}

	// Check footer has HTML comment wrapper
	if !strings.Contains(result, "<!-- <footer> -->") {
		t.Error("Footer should have HTML comment wrapper")
	}

	// Check content is present
	if !strings.Contains(result, "Data Report") {
		t.Error("Expected heading not found")
	}

	if !strings.Contains(result, "Widget A") {
		t.Error("Expected table content not found")
	}
}

package tidio

import . "github.com/gregoryv/web"

func NewChangelog() *Element {
	return Article(
		H1("Changelog"),

		P(`All notable changes to this project will be documented in
		this file. Project adheres to semantic versioning(v2).`),

		Section(
			H2("[unreleased]"),
			Ul(
				Li("Export / Import system state on startup"),
			),
		),

		Section(
			H2("[0.2.0] 2020-06-04"),
			Ul(
				Li("Protecting apis with apikey in Authorization header"),
				Li("/api/timesheets/{account}/yyyymm.timesheet"),
			),
		),
		Section(
			H2("[0.1.0] 2020-06-02"),
		),
	)
}

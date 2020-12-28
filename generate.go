package main

import (
	"os"
	// html/template mangles the xml header at the start, and should be
	// used for untrusted data anyway
	"text/template"
	"strings"
	"bufio"
	"fmt"
)

type LetterInfo struct {
	Num int
	Dateline string
	Lang string
	Endno int
}

type LangInfo struct {
	Num int
	Language string
	Endno int
}

const se_tmpl = `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" epub:prefix="z3998: http://www.daisy.org/z3998/2012/vocab/structure/, se: https://standardebooks.org/vocab/1.0" xml:lang="en-GB">
	<head>
		<title>Letter {{.Num}}</title>
		<link href="../css/core.css" rel="stylesheet" type="text/css"/>
		<link href="../css/local.css" rel="stylesheet" type="text/css"/>
	</head>
	<body epub:type="bodymatter z3998:non-fiction">
		<section id="letter-{{.Num}}" epub:type="chapter z3998:letter">
			<h2>
				<span epub:type="label">Letter</span>
				<span epub:type="ordinal">{{.Num}}</span>{{if ne .Lang "en" }}<a href="endnotes.xhtml#{{.Endno}}" id="noteref-{{.Endno}}" epub:type="noteref">{{.Endno}}</a>{{end}}
			</h2>
			{{if .Dateline}}
			<p epub:type="se:letter.dateline">{{.Dateline}}</p>{{end}}
			<p epub:type="z3998:salutation"></p>
		</section>
	</body>
</html>
`

const end_tmpl = `<?xml version="1.0" encoding="utf-8"?>
<html xmlns="http://www.w3.org/1999/xhtml" xmlns:epub="http://www.idpf.org/2007/ops" epub:prefix="z3998: http://www.daisy.org/z3998/2012/vocab/structure/, se: https://standardebooks.org/vocab/1.0" xml:lang="en-GB">
	<head>
		<title>Endnotes</title>
		<link href="../css/core.css" rel="stylesheet" type="text/css"/>
		<link href="../css/local.css" rel="stylesheet" type="text/css"/>
	</head>
	<body epub:type="backmatter">
		<section id="endnotes" epub:type="endnotes">
			<h2 epub:type="title">Endnotes</h2>
			<ol>
			{{range .}}
				<li id="note-{{.Endno}}" epub:type="endnote">
					<p>This letter was originally written in {{.Language}}. <a href="letter-{{.Num}}.xhtml#noteref-{{.Endno}}" epub:type="backlink">↩</a></p>
				</li>
			{{end}}
			</ol>
		</section>
	</body>
</html>
`

var tmpl *template.Template

func generate(f LetterInfo) {
	file := fmt.Sprintf("src/epub/text/letter-%d.xhtml", f.Num)
	html, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer html.Close()
	err = tmpl.Execute(html, f)
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		err error
		endnotes []LangInfo
		en_tmpl *template.Template
		letter_count = 0
		lang_count = 0
	)

	tmpl, err = template.New("letter").Parse(se_tmpl)
	if err != nil {
		panic(err)
	}
	en_tmpl, err = template.New("endnotes").Parse(end_tmpl)
	if err != nil {
		panic(err)
	}

	f, err := os.Open("matrix")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		arr := strings.Split(s.Text(), "|")
		letter_count++
		switch (arr[1]) {
		case "fr":
			lang_count++
			endnotes = append(endnotes, LangInfo{lang_count, "French", letter_count})
		case "la":
			lang_count++
			endnotes = append(endnotes, LangInfo{lang_count, "Latin", letter_count})
		}
		generate(LetterInfo{letter_count, arr[0], arr[1], lang_count})
	}

	ef, err := os.Create("src/epub/text/endnotes.xhtml")
	if err != nil {
		panic(err)
	}

	en_tmpl.Execute(ef, endnotes)
}

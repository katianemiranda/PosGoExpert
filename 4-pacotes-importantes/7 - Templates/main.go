package main

import (
	"html/template"
	"os"
	"strings"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func ToUpper(s string) string {
	return strings.ToUpper()
}
func main() {
	// curso := Curso{"Go", 40}
	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}

	//t := template.Must(template.New("content.html").ParseFiles(templates...))

	// tmp := template.New("CursoTemplate")
	// tmp, err := tmp.Parse("Curso: {{.Nome}} - Carga Horária: {{.CargaHoraria}}")
	// if err != nil {
	// 	panic(err)
	// }
	//err := tmp.Execute(os.Stdout, curso)

	t := template.New("content.html")
	t.Funcs(template.FuncMap{"ToUpper": ToUpper})
	t = template.Must(t.ParseFiles(templates...))

	err := t.Execute(os.Stdout, Cursos{
		{"Go", 40},
		{"Java", 20},
		{"Phyton", 10},
	})
	if err != nil {
		panic(err)
	}
}

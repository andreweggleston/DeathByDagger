package controllers

import "html/template"

func InitTemplates() {
	mainTempl = template.Must(template.ParseFiles("views/index.html"))
}

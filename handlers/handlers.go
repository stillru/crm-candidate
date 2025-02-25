package handlers

import (
	"crm-candidate/db"
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"content": "Добро пожаловать в CRM для кандидата!",
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func AddCompanyHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"content": "add_company.html",
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func AddRecruiterHandler(w http.ResponseWriter, r *http.Request) {
	companies, err := db.GetCompanies()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"content":   "add_recruiter.html",
		"Companies": companies,
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	interactions, err := db.GetInteractions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"content":      "status.html",
		"Interactions": interactions,
	}
	tmpl.ExecuteTemplate(w, "base.html", data)
}

func SaveCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// Логика сохранения компании
	http.Redirect(w, r, "/add-company", http.StatusFound)
}

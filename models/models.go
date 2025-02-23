package models

import "time"

// Резюме кандидата
type Resume struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Content     string `json:"content"` // YAML-содержимое
    PDFPath     string `json:"pdf_path"` // Путь к PDF-файлу
    CreatedAt   time.Time `json:"created_at"`
}

// Работодатель
type Employer struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Website     string `json:"website"`
    ContactInfo string `json:"contact_info"`
}

// Вакансия
type Vacancy struct {
    ID          int       `json:"id"`
    EmployerID  int       `json:"employer_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    AppliedAt   time.Time `json:"applied_at"`
    Status      string    `json:"status"` // Например, "Отправлено", "Собеседование", "Отказ"
}

// Компания
type Company struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Website     string `json:"website"`
    ContactInfo string `json:"contact_info"`
}

// Рекрутер
type Recruiter struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Phone     string `json:"phone"`
    CompanyID int    `json:"company_id"`
}

// Взаимодействие
type Interaction struct {
    ID           int       `json:"id"`
    CompanyName  string    `json:"company_name"`
    RecruiterName string   `json:"recruiter_name"`
    Status       string    `json:"status"`
    Date         time.Time `json:"date"`
}

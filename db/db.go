package db

import (
    "database/sql"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "crm-candidate/models"
)

var DB *sql.DB

// Инициализация базы данных
func InitDB() error {
    var err error
    DB, err = sql.Open("sqlite3", "./crm.db")
    if err != nil {
        return fmt.Errorf("ошибка подключения к базе данных: %v", err)
    }

    // Создание таблиц
    if err := createTables(); err != nil {
        return fmt.Errorf("ошибка создания таблиц: %v", err)
    }

    return nil
}

// Создание таблиц
func createTables() error {
    queries := []string{
        `CREATE TABLE IF NOT EXISTS companies (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            website TEXT,
            contact_info TEXT
        )`,
        `CREATE TABLE IF NOT EXISTS recruiters (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT NOT NULL,
            phone TEXT,
            company_id INTEGER NOT NULL,
            FOREIGN KEY (company_id) REFERENCES companies (id)
        )`,
        `CREATE TABLE IF NOT EXISTS interactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            recruiter_id INTEGER NOT NULL,
            status TEXT NOT NULL,
            date DATETIME DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (recruiter_id) REFERENCES recruiters (id)
        )`,
    }

    for _, query := range queries {
        _, err := DB.Exec(query)
        if err != nil {
            return err
        }
    }

    return nil
}

// Получение списка компаний
func GetCompanies() ([]models.Company, error) {
    rows, err := DB.Query("SELECT id, name, website, contact_info FROM companies")
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении компаний: %v", err)
    }
    defer rows.Close()

    var companies []models.Company
    for rows.Next() {
        var company models.Company
        if err := rows.Scan(&company.ID, &company.Name, &company.Website, &company.ContactInfo); err != nil {
            return nil, fmt.Errorf("ошибка при сканировании компании: %v", err)
        }
        companies = append(companies, company)
    }

    return companies, nil
}

// Получение списка взаимодействий
func GetInteractions() ([]models.Interaction, error) {
    rows, err := DB.Query(`
        SELECT i.id, c.name AS company_name, r.name AS recruiter_name, i.status, i.date
        FROM interactions i
        JOIN recruiters r ON i.recruiter_id = r.id
        JOIN companies c ON r.company_id = c.id
    `)
    if err != nil {
        return nil, fmt.Errorf("ошибка при получении взаимодействий: %v", err)
    }
    defer rows.Close()

    var interactions []models.Interaction
    for rows.Next() {
        var interaction models.Interaction
        if err := rows.Scan(&interaction.ID, &interaction.CompanyName, &interaction.RecruiterName, &interaction.Status, &interaction.Date); err != nil {
            return nil, fmt.Errorf("ошибка при сканировании взаимодействия: %v", err)
        }
        interactions = append(interactions, interaction)
    }

    return interactions, nil
}

func SaveCompany(company *models.Company) error {
    query := `INSERT INTO companies (name, website, contact_info) VALUES (?, ?, ?)`
    _, err := DB.Exec(query, company.Name, company.Website, company.ContactInfo)
    if err != nil {
        return fmt.Errorf("ошибка при сохранении компании: %v", err)
    }
    return nil
}

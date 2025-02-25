package services

import (
	"crmcandidate/db"
	"crmcandidate/models"
)

type CandidateService struct {
	db *db.DB
}

func NewCandidateService(db *db.DB) *CandidateService {
	return &CandidateService{db: db}
}

func (s *CandidateService) CreateCandidate(candidate *models.Candidate) error {
	query := `INSERT INTO candidates (first_name, last_name, email, phone) VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, candidate.FirstName, candidate.LastName, candidate.Email, candidate.Phone)
	return err
}

func (s *CandidateService) GetCandidate(id int) (*models.Candidate, error) {
	query := `SELECT id, first_name, last_name, email, phone FROM candidates WHERE id = ?`
	row := s.db.QueryRow(query, id)

	candidate := &models.Candidate{}
	err := row.Scan(&candidate.ID, &candidate.FirstName, &candidate.LastName, &candidate.Email, &candidate.Phone)
	if err != nil {
		return nil, err
	}

	return candidate, nil
}

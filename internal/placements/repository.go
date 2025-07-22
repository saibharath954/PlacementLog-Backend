package placements

import (
	"database/sql"
	"fmt"
)

type PlacementsRepo struct {
	db *sql.DB
}

func NewPlacementsRepo(db *sql.DB) *PlacementsRepo {
	return &PlacementsRepo{db: db}
}

func (r *PlacementsRepo) InsertPlacementCompany(company string, ctc float64, placementDate string) (int, error) {
	var id int
	query := `INSERT INTO placement_companies (company, ctc, placement_date) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, company, ctc, placementDate).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert placement company: %w", err)
	}
	return id, nil
}

func (r *PlacementsRepo) InsertBranchwiseRecords(placementID int, branchCounts []BranchCount) error {
	query := `INSERT INTO placement_branchwise_record (placement_id, branch, count) VALUES ($1, $2, $3)`
	for _, bc := range branchCounts {
		_, err := r.db.Exec(query, placementID, bc.Branch, bc.Count)
		if err != nil {
			return fmt.Errorf("failed to insert branchwise record: %w", err)
		}
	}
	return nil
}

func (r *PlacementsRepo) GetAllPlacements() ([]PlacementCompany, error) {
	placements := []PlacementCompany{}
	rows, err := r.db.Query(`SELECT id, company, ctc, placement_date, created_at FROM placement_companies ORDER BY placement_date DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch placements: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p PlacementCompany
		err := rows.Scan(&p.ID, &p.Company, &p.CTC, &p.PlacementDate, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		branchRows, err := r.db.Query(`SELECT branch, count FROM placement_branchwise_record WHERE placement_id = $1`, p.ID)
		if err != nil {
			return nil, err
		}
		var branchCounts []BranchCount
		for branchRows.Next() {
			var bc BranchCount
			if err := branchRows.Scan(&bc.Branch, &bc.Count); err != nil {
				return nil, err
			}
			branchCounts = append(branchCounts, bc)
		}
		branchRows.Close()
		p.BranchCounts = branchCounts
		placements = append(placements, p)
	}
	return placements, nil
}

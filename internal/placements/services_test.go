package placements

import (
	"errors"
	"reflect"
	"testing"
)

type mockPlacementsRepo struct {
	InsertPlacementCompanyFunc  func(company string, ctc float64, placementDate string) (int, error)
	InsertBranchwiseRecordsFunc func(placementID int, branchCounts []BranchCount) error
	GetAllPlacementsFunc        func() ([]PlacementCompany, error)
	GetCompanyBranchMapFunc     func() ([]CompanyBranch, error)
	GetBranchCompanyMapFunc     func() ([]BranchCompany, error)
}

func (m *mockPlacementsRepo) InsertPlacementCompany(company string, ctc float64, placementDate string) (int, error) {
	return m.InsertPlacementCompanyFunc(company, ctc, placementDate)
}
func (m *mockPlacementsRepo) InsertBranchwiseRecords(placementID int, branchCounts []BranchCount) error {
	return m.InsertBranchwiseRecordsFunc(placementID, branchCounts)
}
func (m *mockPlacementsRepo) GetAllPlacements() ([]PlacementCompany, error) {
	return m.GetAllPlacementsFunc()
}
func (m *mockPlacementsRepo) GetCompanyBranchMap() ([]CompanyBranch, error) {
	return m.GetCompanyBranchMapFunc()
}
func (m *mockPlacementsRepo) GetBranchCompanyMap() ([]BranchCompany, error) {
	return m.GetBranchCompanyMapFunc()
}

func TestPlacementsService_AddPlacement(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockPlacementsRepo{
			InsertPlacementCompanyFunc: func(company string, ctc float64, placementDate string) (int, error) {
				return 1, nil
			},
			InsertBranchwiseRecordsFunc: func(placementID int, branchCounts []BranchCount) error {
				return nil
			},
		}
		s := NewPlacementsService(repo)
		resp, err := s.AddPlacement(PlacementRequest{
			Company:       "TestCo",
			CTC:           10.5,
			PlacementDate: "2024-01-01",
			Students:      []string{"22bcs1234", "22bcs5678"},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Company != "TestCo" || resp.CTC != 10.5 {
			t.Errorf("unexpected response: %+v", resp)
		}
	})
	t.Run("placement company insert error", func(t *testing.T) {
		repo := &mockPlacementsRepo{
			InsertPlacementCompanyFunc: func(company string, ctc float64, placementDate string) (int, error) {
				return 0, errors.New("insert error")
			},
		}
		s := NewPlacementsService(repo)
		_, err := s.AddPlacement(PlacementRequest{Company: "TestCo", CTC: 10.5, Students: []string{"22bcs1234"}})
		if err == nil || err.Error() != "insert error" {
			t.Errorf("expected insert error, got %v", err)
		}
	})
	t.Run("branchwise records insert error", func(t *testing.T) {
		repo := &mockPlacementsRepo{
			InsertPlacementCompanyFunc: func(company string, ctc float64, placementDate string) (int, error) {
				return 1, nil
			},
			InsertBranchwiseRecordsFunc: func(placementID int, branchCounts []BranchCount) error {
				return errors.New("branchwise error")
			},
		}
		s := NewPlacementsService(repo)
		_, err := s.AddPlacement(PlacementRequest{Company: "TestCo", CTC: 10.5, Students: []string{"22bcs1234"}})
		if err == nil || err.Error() != "branchwise error" {
			t.Errorf("expected branchwise error, got %v", err)
		}
	})
}

func TestPlacementsService_GetAllPlacements(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		placements := []PlacementCompany{{ID: 1, Company: "TestCo"}}
		repo := &mockPlacementsRepo{
			GetAllPlacementsFunc: func() ([]PlacementCompany, error) { return placements, nil },
		}
		s := NewPlacementsService(repo)
		got, err := s.GetAllPlacements()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, placements) {
			t.Errorf("expected %v, got %v", placements, got)
		}
	})
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPlacementsRepo{
			GetAllPlacementsFunc: func() ([]PlacementCompany, error) { return nil, errors.New("db error") },
		}
		s := NewPlacementsService(repo)
		_, err := s.GetAllPlacements()
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPlacementsService_GetCompanyBranchMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cb := []CompanyBranch{{Company: "TestCo"}}
		repo := &mockPlacementsRepo{
			GetCompanyBranchMapFunc: func() ([]CompanyBranch, error) { return cb, nil },
		}
		s := NewPlacementsService(repo)
		got, err := s.GetCompanyBranchMap()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, cb) {
			t.Errorf("expected %v, got %v", cb, got)
		}
	})
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPlacementsRepo{
			GetCompanyBranchMapFunc: func() ([]CompanyBranch, error) { return nil, errors.New("db error") },
		}
		s := NewPlacementsService(repo)
		_, err := s.GetCompanyBranchMap()
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestPlacementsService_GetBranchCompanyMap(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		bc := []BranchCompany{{Branch: "bcs"}}
		repo := &mockPlacementsRepo{
			GetBranchCompanyMapFunc: func() ([]BranchCompany, error) { return bc, nil },
		}
		s := NewPlacementsService(repo)
		got, err := s.GetBranchCompanyMap()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, bc) {
			t.Errorf("expected %v, got %v", bc, got)
		}
	})
	t.Run("repo error", func(t *testing.T) {
		repo := &mockPlacementsRepo{
			GetBranchCompanyMapFunc: func() ([]BranchCompany, error) { return nil, errors.New("db error") },
		}
		s := NewPlacementsService(repo)
		_, err := s.GetBranchCompanyMap()
		if err == nil || err.Error() != "db error" {
			t.Errorf("expected db error, got %v", err)
		}
	})
}

func TestGetBranchFromRegNo(t *testing.T) {
	cases := []struct {
		regNo string
		want  string
	}{
		{"22bcs1234", "bcs"},
		{"22mec5678", "mec"},
		{"", ""},
		{"12", ""},
		{"22EEE1234", "eee"}, // upper case, should be lower
		{"22bcs", "bcs"},
		{"22", ""},
	}
	for _, c := range cases {
		got := GetBranchFromRegNo(c.regNo)
		if got != c.want {
			t.Errorf("GetBranchFromRegNo(%q) = %q; want %q", c.regNo, got, c.want)
		}
	}
}

func TestCountBranches(t *testing.T) {
	cases := []struct {
		regNos []string
		want   map[string]int
	}{
		{[]string{"22bcs1234", "22bcs5678", "22mec1234"}, map[string]int{"bcs": 2, "mec": 1}},
		{[]string{}, map[string]int{}},
		{[]string{"", "12"}, map[string]int{}},
		{[]string{"22bcs1234", "22bcs1234"}, map[string]int{"bcs": 2}},
	}
	for _, c := range cases {
		got := CountBranches(c.regNos)
		gotMap := map[string]int{}
		for _, bc := range got {
			gotMap[bc.Branch] = bc.Count
		}
		if !reflect.DeepEqual(gotMap, c.want) {
			t.Errorf("CountBranches(%v) = %v; want %v", c.regNos, gotMap, c.want)
		}
	}
}

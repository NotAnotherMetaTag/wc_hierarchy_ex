package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLastName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Michael Chen", "Chen"},
		{"Julian Early", "Early"},
		{"SingleName", "SingleName"},
		{"Riley O'Conner", "O'Conner"},
	}

	for _, c := range cases {
		got := getLastName(c.input)
		if got != c.want {
			t.Errorf("getLastName(%q) == %q, want %q", c.input, got, c.want)
		}
	}
}

func TestBuildHierarchy_SortsByLastName(t *testing.T) {
	a := Employee{Name: "Manager Andrews", ID: 1}
	b := Employee{Name: "Zelda Bossman", ID: 2, ManagerID: &a.ID}
	c := Employee{Name: "Aaron Clark", ID: 3, ManagerID: &a.ID}

	reports := map[int][]Employee{
		1: {b, c},
	}

	result := buildHierarchy(a, reports, map[int]bool{})

	if len(result.Reports) != 2 {
		t.Fatalf("Expected 2 reports, got %d", len(result.Reports))
	}

	if result.Reports[0].Name != "Zelda Bossman" {
		t.Errorf("Expected first report to be 'Zelda Bossman', got %s", result.Reports[0].Name)
	}
}

func TestBuildHierarchy(t *testing.T) {
	emp1 := Employee{Name: "Manager A", ID: 1}
	emp2 := Employee{Name: "Worker B", ID: 2, ManagerID: &emp1.ID}
	emp3 := Employee{Name: "Worker C", ID: 3, ManagerID: &emp1.ID}

	reports := map[int][]Employee{
		1: {emp3, emp2}, // Unsorted on purpose
	}

	result := buildHierarchy(emp1, reports, map[int]bool{})

	if len(result.Reports) != 2 {
		t.Fatalf("Expected 2 reports, got %d", len(result.Reports))
	}

	if getLastName(result.Reports[0].Name) > getLastName(result.Reports[1].Name) {
		t.Errorf("Reports are not sorted by last name")
	}
}

func TestBuildHierarchy_CircularReference(t *testing.T) {
	aID, bID := 1, 2
	a := Employee{Name: "A", ID: aID, ManagerID: &bID}
	b := Employee{Name: "B", ID: bID, ManagerID: &aID}

	reports := map[int][]Employee{
		aID: {b},
		bID: {a},
	}

	// don't loop, just end at the first one
	result := buildHierarchy(a, reports, map[int]bool{})

	if len(result.Reports) != 1 {
		t.Errorf("Expected 1 report, got %d", len(result.Reports))
	}

	if len(result.Reports[0].Reports) != 0 {
		t.Errorf("Expected circular reference to stop, but got nested reports")
	}
}

func TestValidateEmployees(t *testing.T) {
	tests := []struct {
		name      string
		employees []Employee
		wantErr   bool
	}{
		{
			"valid",
			[]Employee{
				{ID: 1, Name: "A"},
				{ID: 2, Name: "B", ManagerID: intPtr(1)},
			},
			false,
		},
		{
			"duplicate ID",
			[]Employee{
				{ID: 1, Name: "A"},
				{ID: 1, Name: "B"},
			},
			true,
		},
		{
			"invalid manager ID",
			[]Employee{
				{ID: 1, Name: "A", ManagerID: intPtr(99)},
			},
			true,
		},
	}

	for _, tt := range tests {
		err := validateEmployees(tt.employees)
		t.Logf("Running test: %s", tt.name)
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", tt.name)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
		}
	}
}

// for Mock Integration Testing
var mockEmployees = []Employee{
	{ID: 1, Name: "Michael Chen", Title: "CEO", ManagerID: nil},
	{ID: 2, Name: "Barrett Glasauer", Title: "CTO", ManagerID: intPtr(1)},
}

// helper for struct literals
func intPtr(i int) *int { return &i }

func TestEmployeeHandler_ReturnsNestedOrgChart(t *testing.T) {
	// reset data and init mock
	employeeData = make(map[int]Employee)
	for _, e := range mockEmployees {
		employeeData[e.ID] = e
	}

	req := httptest.NewRequest("GET", "/api/employees", nil)
	w := httptest.NewRecorder()
	employeeHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", res.StatusCode)
	}

	var got []Employee
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("Expected 1 top-level employee, got %d", len(got))
	}

	if got[0].Name != "Michael Chen" {
		t.Errorf("Expected CEO to be 'Michael Chen', got '%s'", got[0].Name)
	}

	if len(got[0].Reports) != 1 {
		t.Errorf("Expected CEO to have 1 report, got %d", len(got[0].Reports))
	}
}

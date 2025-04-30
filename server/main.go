// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"sync"
)

type Employee struct {
	Name      string     `json:"name"`
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	ManagerID *int       `json:"manager_id"`
	Reports   []Employee `json:"reports,omitempty"`
}

var (
	employeeData = make(map[int]Employee)
	once         sync.Once
)

const remoteURL = "https://gist.githubusercontent.com/chancock09/6d2a5a4436dcd488b8287f3e3e4fc73d/raw/fa47d64c6d5fc860fabd3033a1a4e3c59336324e/employees.json"

func main() {
	http.HandleFunc("/api/employees", employeeHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		if len(employeeData) == 0 {
			fetchRemoteData()
		}
	})

	orgChart := buildOrgChart()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orgChart)
}

func fetchRemoteData() {
	resp, err := http.Get(remoteURL)
	if err != nil {
		log.Fatal("Failed to fetch remote data:", err)
	}
	defer resp.Body.Close()

	var employees []Employee
	if err := json.NewDecoder(resp.Body).Decode(&employees); err != nil {
		log.Fatal("Failed to decode JSON:", err)
	}

	for _, emp := range employees {
		employeeData[emp.ID] = emp
	}
}

func buildOrgChart() []Employee {
	reports := make(map[int][]Employee)
	for _, emp := range employeeData {
		if emp.ManagerID != nil {
			reports[*emp.ManagerID] = append(reports[*emp.ManagerID], emp)
		}
	}

	var org []Employee
	for _, emp := range employeeData {
		if emp.ManagerID == nil {
			org = append(org, buildHierarchy(emp, reports, map[int]bool{}))
		}
	}
	return org
}

func buildHierarchy(emp Employee, reports map[int][]Employee, visited map[int]bool) Employee {
	if visited[emp.ID] {
		// prevent circular ref from blowing up our app
		return Employee{
			ID:        emp.ID,
			Name:      emp.Name,
			Title:     emp.Title,
			ManagerID: emp.ManagerID,
			Reports:   nil,
		}
	}
	visited[emp.ID] = true

	children := reports[emp.ID]
	sort.Slice(children, func(i, j int) bool {
		return getLastName(children[i].Name) < getLastName(children[j].Name)
	})

	// make a new slice to assign the results to
	var resultChildren []Employee
	for _, child := range children {
		if visited[child.ID] {
			continue // skip appending child if circular relationship
		}

		built := buildHierarchy(child, reports, visited)
		resultChildren = append(resultChildren, built)
	}

	emp.Reports = resultChildren

	return emp
}

func getLastName(name string) string {
	parts := []rune(name)
	lastSpace := -1
	for i, ch := range parts {
		if ch == ' ' {
			lastSpace = i
		}
	}
	if lastSpace != -1 {
		return string(parts[lastSpace+1:])
	}
	return name
}

func validateEmployees(employees []Employee) error {
	ids := make(map[int]bool)
	for _, emp := range employees {
		if ids[emp.ID] {
			return fmt.Errorf("duplicate ID found: %d", emp.ID)
		}
		ids[emp.ID] = true
	}
	for _, emp := range employees {
		if emp.ManagerID != nil && !ids[*emp.ManagerID] {
			return fmt.Errorf("invalid manager ID %d for employee %d", *emp.ManagerID, emp.ID)
		}
	}
	return nil
}

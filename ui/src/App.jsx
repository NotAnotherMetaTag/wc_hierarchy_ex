import React, { useEffect, useState } from 'react';
import './App.css';

function Employee({ employee, level = 0 }) {
  return (
    <div className="OrgHierarchy" style={{ marginLeft: level * 20 }}>
      <strong>{employee.title}</strong>: {employee.name}
      {employee.reports?.length > 0 && (
        <div>
          {employee.reports.map((report) => (
            <Employee key={report.id} employee={report} level={level + 1} />
          ))}
        </div>
      )}
    </div>
  );
}

function App() {
  const [employees, setEmployees] = useState([]);
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchEmployees = () => {
    setLoading(true);
    setError(null);

    fetch('/api/employees')
      .then((res) => {
        if (!res.ok) throw new Error(`Server responded with ${res.status}`);
        return res.json();
      })
      .then((data) => {
        setEmployees(data);
      })
      .catch((err) => {
        console.error(err);
        setError('Could not load employees. Please make sure the server is running.');
        setEmployees([]); // fallback to empty list
      })
      .finally(() => {
        setLoading(false);
      });
  };

  // Initial fetch
  useEffect(() => {
    fetchEmployees();
  }, []);

  return (
    <div className="App">
      <h1>Organization Chart</h1>

      {loading && <p>Loading employees...</p>}

      {!loading && error && (
        <div style={{ color: 'red' }}>
          <p>{error}</p>
          <button onClick={fetchEmployees}>Retry</button>
        </div>
      )}

      {!loading && !error && employees.length === 0 && (
        <p>No employees to display.</p>
      )}

      {!loading && !error && employees.map(emp => (
        <Employee key={emp.id} employee={emp} />
      ))}
    </div>
  );
}

export default App;

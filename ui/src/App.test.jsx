/* eslint-disable no-undef */
import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import App from './App';

beforeEach(() => {
  // Clear mock before each test
  global.fetch = vi.fn();
});

test('renders employee list on successful fetch', async () => {
  fetch.mockResolvedValueOnce({
    ok: true,
    json: async () => [
      { id: 1, name: 'Michael Chen', title: 'CEO', reports: [] },
    ],
  });

  render(<App />);
  expect(screen.getByText(/loading employees/i)).toBeInTheDocument();

  await waitFor(() => {
    expect(
      screen.getByText((_, element) =>
        element?.textContent === 'CEO: Michael Chen'
      )
    ).toBeInTheDocument();
  });
});

test('shows error and retry button on failed fetch', async () => {
  fetch.mockRejectedValueOnce(new Error('Network error'));

  render(<App />);
  await waitFor(() =>
    expect(
      screen.getByText(/could not load employees/i)
    ).toBeInTheDocument()
  );

  expect(
    screen.getByRole('button', { name: /retry/i })
  ).toBeInTheDocument();
});

test('clicking Retry triggers a new fetch and renders data', async () => {
  // initially fails
  fetch
    .mockRejectedValueOnce(new Error('Network down'))
    // retry succeeds
    .mockResolvedValueOnce({
      ok: true,
      json: async () => [
        { id: 1, name: 'Michael Chen', title: 'CEO', reports: [] },
      ],
    });

  render(<App />);

  await waitFor(() =>
    expect(
      screen.getByText(/could not load employees/i)
    ).toBeInTheDocument()
  );

  // click retry
  fireEvent.click(screen.getByRole('button', { name: /retry/i }));

  await waitFor(() =>
    expect(
      screen.getByText((_, element) =>
        element?.textContent === 'CEO: Michael Chen'
      )
    ).toBeInTheDocument()
  );
});

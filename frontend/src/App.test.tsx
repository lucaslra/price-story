import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import App from './App'
vi.mock('@/api/client', () => ({
  get: async () => ({ products: [] })
}))

describe('App', () => {
  it('renders Products heading after data loads', async () => {
    // Ensure route is /products
    // (API client is mocked above to return empty list)

    // Navigate to the products route so App renders the Products page
    window.history.pushState({}, '', '/products')

    render(<App />)

    const heading = await screen.findByRole('heading', { level: 1, name: 'Products' })
    expect(heading).toBeInTheDocument()
  })
})

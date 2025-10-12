import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import App from './App'

describe('App', () => {
  it('renders Products heading after data loads', async () => {
    vi.spyOn(global, 'fetch').mockResolvedValueOnce(
      new Response(JSON.stringify({ products: [] }), {
        status: 200,
        headers: { 'Content-Type': 'application/json' }
      } as ResponseInit)
    )

    render(<App />)

    const heading = await screen.findByRole('heading', { level: 1, name: 'Products' })
    expect(heading).toBeInTheDocument()
  })
})

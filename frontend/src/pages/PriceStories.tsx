import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { get } from '@/api/client'
import type { Product } from '@/types'

type PriceStory = {
  id: string
  product: Product
  created_datetime?: string
}

type PriceStoriesResponse = { priceStories: PriceStory[] }

export default function PriceStoriesPage() {
  const [stories, setStories] = useState<PriceStory[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    get<PriceStoriesResponse>('/price-stories')
      .then((res) => {
        const list = Array.isArray(res.priceStories) ? res.priceStories : []
        setStories(list)
        setLoading(false)
      })
      .catch((e) => {
        setError(e?.message ?? 'Failed to load price stories')
        setLoading(false)
      })
  }, [])

  if (loading) return <p style={{ padding: 24 }}>Loading price stories…</p>
  if (error) return <p style={{ padding: 24, color: 'crimson' }}>{error}</p>

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          marginBottom: 16
        }}
      >
        <h1 style={{ margin: 0 }}>Price Stories</h1>
        <Link to="/price-stories/new">New Price Story</Link>
      </div>
      {stories.length === 0 ? (
        <p>
          No price stories found. <Link to="/price-stories/new">Create one</Link>
        </p>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 12 }}>
          {stories.map((ps) => (
            <li
              key={ps.id}
              style={{
                border: '1px solid #ddd',
                borderRadius: 8,
                padding: 12,
                display: 'grid',
                gridTemplateColumns: '96px 1fr',
                gap: 12
              }}
            >
              {ps.product?.product_image_url ? (
                <img
                  src={ps.product.product_image_url}
                  alt={ps.product.product_name}
                  loading="lazy"
                  style={{ width: 96, height: 96, objectFit: 'cover', borderRadius: 6 }}
                />
              ) : (
                <div
                  style={{
                    width: 96,
                    height: 96,
                    borderRadius: 6,
                    background: '#f3f3f3',
                    display: 'grid',
                    placeItems: 'center',
                    color: '#888'
                  }}
                >
                  No image
                </div>
              )}
              <div>
                <div style={{ display: 'flex', alignItems: 'baseline', gap: 8 }}>
                  <Link to={`/price-stories/${ps.id}`} style={{ textDecoration: 'none' }}>
                    <strong>{ps.product?.product_name ?? 'Unknown product'}</strong>
                  </Link>
                  {ps.product?.product_url ? (
                    <a href={ps.product.product_url} target="_blank" rel="noreferrer">
                      Visit
                    </a>
                  ) : null}
                </div>
                <div style={{ marginTop: 6, color: '#333' }}>
                  <span>ID: {ps.id}</span>
                  {ps.created_datetime ? (
                    <span style={{ marginLeft: 8, color: '#666' }}>
                      Created: {new Date(ps.created_datetime).toLocaleString()}
                    </span>
                  ) : null}
                </div>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

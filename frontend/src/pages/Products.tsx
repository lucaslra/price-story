import { useEffect, useState } from 'react'
import { get } from '@/api/client'
import type { Product } from '@/types'

type ProductsResponse = { products: Product[] }

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    get<ProductsResponse>('/products')
      .then((res) => {
        const list = Array.isArray(res.products) ? res.products : []
        setProducts(list)
        setLoading(false)
      })
      .catch((e) => {
        setError(e?.message ?? 'Failed to load products')
        setLoading(false)
      })
  }, [])

  if (loading) return <p style={{ padding: 24 }}>Loading products…</p>
  if (error) return <p style={{ padding: 24, color: 'crimson' }}>{error}</p>

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <h1 style={{ marginBottom: 16 }}>Products</h1>
      {products.length === 0 ? (
        <p>No products found.</p>
      ) : (
        <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 12 }}>
          {products.map((p) => (
            <li
              key={p.id}
              style={{
                border: '1px solid #ddd',
                borderRadius: 8,
                padding: 12,
                display: 'flex',
                alignItems: 'center',
                gap: 12
              }}
            >
              {p.product_image_url ? (
                <img
                  src={p.product_image_url}
                  alt={p.product_name}
                  loading="lazy"
                  style={{ width: 96, height: 96, objectFit: 'cover', borderRadius: 6 }}
                />
              ) : null}
              <div style={{ flex: 1 }}>
                <div style={{ display: 'flex', alignItems: 'baseline', gap: 8 }}>
                  <strong>{p.product_name}</strong>
                  <a href={p.product_url} target="_blank" rel="noreferrer">
                    Visit
                  </a>
                </div>
                {p.product_description ? (
                  <p style={{ margin: '4px 0 0', color: '#555' }}>{p.product_description}</p>
                ) : null}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

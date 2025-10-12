import { useEffect, useMemo, useState } from 'react'
import { useParams } from 'react-router-dom'
import { get } from '@/api/client'
import type { Product } from '@/types'

type PricePoint = {
  id: string
  price: number
  timestamp: string
  product: Product
}

type PricePointsResponse = { pricePoints: PricePoint[] }

export default function PricePointsPage() {
  const { productId } = useParams<{ productId: string }>()
  const [points, setPoints] = useState<PricePoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    get<PricePointsResponse>('/price-points')
      .then((res) => {
        const list = Array.isArray(res.pricePoints) ? res.pricePoints : []
        setPoints(list)
        setLoading(false)
      })
      .catch((e) => {
        setError(e?.message ?? 'Failed to load price points')
        setLoading(false)
      })
  }, [])

  const visible = useMemo(() => {
    if (!productId) return points
    return points.filter((pp) => pp.product?.id === productId)
  }, [points, productId])

  const product = useMemo(() => {
    return visible.length > 0 ? visible[0].product : undefined
  }, [visible])

  const sorted = useMemo(() => {
    return [...visible].sort(
      (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    )
  }, [visible])

  if (loading) return <p style={{ padding: 24 }}>Loading price points…</p>
  if (error) return <p style={{ padding: 24, color: 'crimson' }}>{error}</p>

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      {productId ? (
        <>
          <h1 style={{ marginBottom: 12 }}>Price Points</h1>
          {product ? (
            <div
              style={{
                border: '1px solid #ddd',
                borderRadius: 8,
                padding: 12,
                display: 'grid',
                gridTemplateColumns: '96px 1fr',
                gap: 12,
                marginBottom: 16
              }}
            >
              {product.product_image_url ? (
                <img
                  src={product.product_image_url}
                  alt={product.product_name}
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
                  <strong>{product.product_name ?? 'Unknown product'}</strong>
                  {product.product_url ? (
                    <a href={product.product_url} target="_blank" rel="noreferrer">
                      Visit
                    </a>
                  ) : null}
                </div>
                {productId ? (
                  <div style={{ marginTop: 4, color: '#666' }}>ID: {productId}</div>
                ) : null}
              </div>
            </div>
          ) : null}

          {sorted.length === 0 ? (
            <p>No price points found for this product.</p>
          ) : (
            <table style={{ width: '100%', borderCollapse: 'collapse' }}>
              <thead>
                <tr>
                  <th
                    style={{
                      textAlign: 'left',
                      padding: '8px 6px',
                      borderBottom: '1px solid #ddd'
                    }}
                  >
                    Date
                  </th>
                  <th
                    style={{
                      textAlign: 'right',
                      padding: '8px 6px',
                      borderBottom: '1px solid #ddd'
                    }}
                  >
                    Price
                  </th>
                </tr>
              </thead>
              <tbody>
                {sorted.map((pp) => (
                  <tr key={pp.id}>
                    <td style={{ padding: '8px 6px', borderBottom: '1px solid #eee' }}>
                      {new Date(pp.timestamp).toLocaleString()}
                    </td>
                    <td
                      style={{
                        padding: '8px 6px',
                        borderBottom: '1px solid #eee',
                        textAlign: 'right',
                        fontVariantNumeric: 'tabular-nums'
                      }}
                    >
                      ${pp.price.toFixed(2)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      ) : (
        <>
          <h1 style={{ marginBottom: 16 }}>Price Points</h1>
          {visible.length === 0 ? (
            <p>No price points found.</p>
          ) : (
            <ul style={{ listStyle: 'none', padding: 0, margin: 0, display: 'grid', gap: 12 }}>
              {visible.map((pp) => (
                <li
                  key={pp.id}
                  style={{
                    border: '1px solid #ddd',
                    borderRadius: 8,
                    padding: 12,
                    display: 'grid',
                    gridTemplateColumns: '96px 1fr',
                    gap: 12
                  }}
                >
                  {pp.product?.product_image_url ? (
                    <img
                      src={pp.product.product_image_url}
                      alt={pp.product.product_name}
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
                      <strong>{pp.product?.product_name ?? 'Unknown product'}</strong>
                      {pp.product?.product_url ? (
                        <a href={pp.product.product_url} target="_blank" rel="noreferrer">
                          Visit
                        </a>
                      ) : null}
                    </div>
                    <div style={{ marginTop: 6, color: '#333' }}>
                      <span style={{ fontVariantNumeric: 'tabular-nums' }}>
                        ${pp.price.toFixed(2)}
                      </span>
                      <span style={{ marginLeft: 8, color: '#666' }}>
                        {new Date(pp.timestamp).toLocaleString()}
                      </span>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </>
      )}
    </div>
  )
}

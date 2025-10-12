import { useEffect, useMemo, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { get } from '@/api/client'
import type { Product } from '@/types'

type PriceStory = {
  id: string
  product: Product
  created_datetime?: string
}

type PricePoint = {
  id: string
  price: number
  timestamp: string
  product: Product
}

type PricePointsResponse = { pricePoints: PricePoint[] }

export default function PriceStoryDetailsPage() {
  const { id } = useParams<{ id: string }>()
  const [story, setStory] = useState<PriceStory | null>(null)
  const [points, setPoints] = useState<PricePoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let mounted = true
    async function load() {
      try {
        if (!id) throw new Error('Missing price story id')
        const ps = await get<PriceStory>(`/price-stories/${id}`)
        if (!mounted) return
        setStory(ps)
        const res = await get<PricePointsResponse>('/price-points')
        if (!mounted) return
        const list = Array.isArray(res.pricePoints) ? res.pricePoints : []
        const filtered = list.filter((pp) => pp.product?.id === ps.product?.id)
        setPoints(filtered)
        setLoading(false)
      } catch (e: unknown) {
        if (!mounted) return
        setError(e instanceof Error ? e.message : 'Failed to load price story')
        setLoading(false)
      }
    }
    load()
    return () => {
      mounted = false
    }
  }, [id])

  const sorted = useMemo(() => {
    return [...points].sort(
      (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    )
  }, [points])

  const chartData = useMemo(() => {
    if (sorted.length === 0) return null
    const w = 720
    const h = 240
    const padding = 24
    const times = sorted.map((p) => new Date(p.timestamp).getTime())
    const prices = sorted.map((p) => p.price)
    const minT = Math.min(...times)
    const maxT = Math.max(...times)
    const minP = Math.min(...prices)
    const maxP = Math.max(...prices)
    const tSpan = Math.max(1, maxT - minT)
    const pSpan = Math.max(1e-9, maxP - minP)
    const x = (t: number) => padding + ((t - minT) / tSpan) * (w - padding * 2)
    const y = (p: number) => padding + (1 - (p - minP) / pSpan) * (h - padding * 2)
    const pointsStr = sorted
      .map((p) => `${x(new Date(p.timestamp).getTime())},${y(p.price)}`)
      .join(' ')
    return { w, h, padding, minP, maxP, minT, maxT, pointsStr }
  }, [sorted])

  if (loading) return <p style={{ padding: 24 }}>Loading price story…</p>
  if (error) return <p style={{ padding: 24, color: 'crimson' }}>{error}</p>
  if (!story) return <p style={{ padding: 24 }}>Price story not found.</p>

  const product = story.product

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
        <h1 style={{ margin: 0 }}>Price Story</h1>
        <span style={{ color: '#666' }}>ID: {story.id}</span>
      </div>

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
          <div style={{ marginTop: 4, color: '#666' }}>Product ID: {product.id}</div>
        </div>
      </div>

      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          marginBottom: 8
        }}
      >
        <h2 style={{ margin: 0, fontSize: 18 }}>Price evolution</h2>
        <Link to={`/price-points/${product.id}`}>View price points table</Link>
      </div>

      {sorted.length === 0 ? (
        <p>No price points found for this product.</p>
      ) : (
        <div style={{ border: '1px solid #ddd', borderRadius: 8, padding: 12 }}>
          {chartData && (
            <svg
              width={chartData.w}
              height={chartData.h}
              role="img"
              aria-label="Price evolution chart"
            >
              <rect x={0} y={0} width={chartData.w} height={chartData.h} fill="#fff" />
              {/* axes */}
              <line
                x1={chartData.padding}
                y1={chartData.h - chartData.padding}
                x2={chartData.w - chartData.padding}
                y2={chartData.h - chartData.padding}
                stroke="#ccc"
              />
              <line
                x1={chartData.padding}
                y1={chartData.padding}
                x2={chartData.padding}
                y2={chartData.h - chartData.padding}
                stroke="#ccc"
              />
              {/* line */}
              <polyline fill="none" stroke="#2a7bff" strokeWidth={2} points={chartData.pointsStr} />
              {/* min/max labels */}
              <text x={chartData.padding} y={chartData.padding - 6} fontSize={12} fill="#555">
                ${chartData.maxP.toFixed(2)}
              </text>
              <text
                x={chartData.padding}
                y={chartData.h - chartData.padding + 14}
                fontSize={12}
                fill="#555"
              >
                {new Date(chartData.minT).toLocaleDateString()}
              </text>
              <text
                x={chartData.w - chartData.padding}
                y={chartData.h - chartData.padding + 14}
                fontSize={12}
                fill="#555"
                textAnchor="end"
              >
                {new Date(chartData.maxT).toLocaleDateString()}
              </text>
              <text
                x={chartData.w - chartData.padding}
                y={chartData.padding - 6}
                fontSize={12}
                fill="#555"
                textAnchor="end"
              >
                ${chartData.minP.toFixed(2)}
              </text>
            </svg>
          )}
        </div>
      )}
    </div>
  )
}

import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { post } from '@/api/client'
import type { Product } from '@/types'

export default function NewPriceStoryPage() {
  const navigate = useNavigate()

  const [productName, setProductName] = useState('')
  const [productUrl, setProductUrl] = useState('')
  const [productImageUrl, setProductImageUrl] = useState('')
  const [price, setPrice] = useState('')
  const [timestamp, setTimestamp] = useState('') // date value (day precision)

  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)

  function validate(): string[] {
    const errors: string[] = []
    if (!productName.trim()) errors.push('Product name is required.')
    if (!productUrl.trim()) errors.push('Product URL is required.')
    else {
      try {
        // Basic URL validation
        const u = new URL(productUrl)
        if (!(u.protocol === 'http:' || u.protocol === 'https:'))
          throw new Error('Invalid protocol')
      } catch {
        errors.push('Product URL must be a valid http(s) URL.')
      }
    }
    const priceNum = Number(price)
    if (!Number.isFinite(priceNum) || priceNum <= 0) errors.push('Price must be a positive number.')
    if (!timestamp) errors.push('Date is required.')
    else if (isNaN(new Date(`${timestamp}T00:00:00Z`).getTime())) errors.push('Invalid date.')
    if (productImageUrl.trim()) {
      try {
        const u = new URL(productImageUrl)
        if (!(u.protocol === 'http:' || u.protocol === 'https:'))
          throw new Error('Invalid protocol')
      } catch {
        errors.push('Product image URL must be a valid http(s) URL.')
      }
    }
    return errors
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setSubmitError(null)
    const errs = validate()
    if (errs.length) {
      setSubmitError(errs[0])
      return
    }
    setSubmitting(true)
    try {
      const product = await post<
        { product_name: string; product_url: string; product_image_url?: string },
        Product
      >('/products', {
        product_name: productName.trim(),
        product_url: productUrl.trim(),
        product_image_url: productImageUrl.trim() || undefined
      })
      const priceStory = await post<{ product_id: string }, { id: string }>('/price-stories', {
        product_id: product.id
      })
      const iso = `${timestamp}T00:00:00Z`
      await post<{ price: number; timestamp: string; product_id: string }, unknown>(
        '/price-points',
        { price: Number(price), timestamp: iso, product_id: product.id }
      )
      navigate(`/price-stories/${priceStory.id}`)
    } catch (e) {
      setSubmitError(e instanceof Error ? e.message : String(e))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <div style={{ marginBottom: 16 }}>
        <Link to="/price-stories">← Back to Price Stories</Link>
      </div>
      <h1>Create Price Story</h1>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gap: 16, maxWidth: 720 }}>
        <section>
          <h2 style={{ fontSize: 18, margin: 0 }}>Product</h2>
          <div style={{ display: 'grid', gap: 8 }}>
            <label style={{ display: 'grid', gap: 4 }}>
              <div>Product Name</div>
              <input
                type="text"
                value={productName}
                onChange={(e) => setProductName(e.target.value)}
                placeholder="e.g. Example Widget"
                required
              />
            </label>
            <label style={{ display: 'grid', gap: 4 }}>
              <div>Product URL</div>
              <input
                type="url"
                value={productUrl}
                onChange={(e) => setProductUrl(e.target.value)}
                placeholder="https://store.example.com/product/123"
                required
              />
            </label>
            <label style={{ display: 'grid', gap: 4 }}>
              <div>Product Image URL (optional)</div>
              <input
                type="url"
                value={productImageUrl}
                onChange={(e) => setProductImageUrl(e.target.value)}
                placeholder="https://images.example.com/product.jpg"
              />
            </label>
          </div>
        </section>

        <section>
          <h2 style={{ fontSize: 18, margin: 0 }}>Initial Price Point</h2>
          <div style={{ display: 'grid', gap: 8 }}>
            <label style={{ display: 'grid', gap: 4 }}>
              <div>Price</div>
              <input
                type="number"
                min="0"
                step="0.01"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                placeholder="0.00"
                required
              />
            </label>
            <label style={{ display: 'grid', gap: 4 }}>
              <div>Date</div>
              <input
                type="date"
                value={timestamp}
                onChange={(e) => setTimestamp(e.target.value)}
                required
              />
            </label>
          </div>
        </section>

        {submitError ? <p style={{ color: 'crimson' }}>{submitError}</p> : null}

        <div>
          <button type="submit" disabled={submitting}>
            {submitting ? 'Creating…' : 'Create Price Story'}
          </button>
        </div>
      </form>
    </div>
  )
}

import { useEffect, useMemo, useState } from 'react'
import { useAuth } from '@/auth/AuthContext'
import type { User } from '@/types'
import { put } from '@/api/client'

export default function ProfilePage() {
  const { user, login } = useAuth()
  const [email, setEmail] = useState(user?.email ?? '')
  const [password, setPassword] = useState('')
  const [preferredCurrency, setPreferredCurrency] = useState(user?.preferred_currency ?? 'USD')
  const [decimalPlaces, setDecimalPlaces] = useState<number>(user?.decimal_places ?? 2)
  const [thousandSeparator, setThousandSeparator] = useState<string>(
    user?.thousand_separator ?? ','
  )
  const [symbolPlacement, setSymbolPlacement] = useState<'before' | 'after'>(
    (user?.currency_symbol_placement as 'before' | 'after') ?? 'before'
  )
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const currencyOptions = useMemo(
    () => ['USD', 'EUR', 'JPY', 'GBP', 'AUD', 'CAD', 'CHF', 'CNY'],
    []
  )

  useEffect(() => {
    if (preferredCurrency === 'JPY') {
      setDecimalPlaces(0)
    } else if (decimalPlaces === 0) {
      setDecimalPlaces(2)
    }
  }, [preferredCurrency, decimalPlaces])

  if (!user) {
    return (
      <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
        <p>You must be logged in to view your profile.</p>
      </div>
    )
  }

  function validate(): string[] {
    const errs: string[] = []
    if (!email.trim()) errs.push('Email is required.')
    if (!password.trim()) errs.push('Password is required to update profile.')
    if (!/^[A-Z]{3}$/.test(preferredCurrency)) {
      errs.push('Preferred currency must be a 3-letter uppercase code.')
    }
    if (preferredCurrency === 'JPY' && decimalPlaces !== 0) {
      errs.push('JPY requires 0 decimal places.')
    }
    if (decimalPlaces < 0 || decimalPlaces > 4) {
      errs.push('Decimal places must be between 0 and 4.')
    }
    if (![',', '.', ' '].includes(thousandSeparator)) {
      errs.push('Thousand separator must be one of ",", ".", or space.')
    }
    if (!['before', 'after'].includes(symbolPlacement)) {
      errs.push('Currency symbol placement must be "before" or "after".')
    }
    return errs
  }

  async function handleSave(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    setError(null)
    setSuccess(null)
    const userId = user?.id
    if (!userId) {
      setError('Missing user id')
      setSaving(false)
      return
    }
    const errs = validate()
    if (errs.length > 0) {
      setError(errs.join(' '))
      setSaving(false)
      return
    }
    try {
      const updated = await put<
        {
          email: string
          password_hash: string
          preferred_currency: string
          decimal_places: number
          thousand_separator: string
          currency_symbol_placement: 'before' | 'after'
        },
        User
      >(`/users/${userId}`, {
        email,
        password_hash: password,
        preferred_currency: preferredCurrency,
        decimal_places: decimalPlaces,
        thousand_separator: thousandSeparator,
        currency_symbol_placement: symbolPlacement
      })
      login(updated)
      setSuccess('Profile updated successfully.')
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setSaving(false)
    }
  }

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <h1>Profile</h1>
      <p style={{ color: '#555', marginTop: -8 }}>Update your account and display preferences.</p>

      <form onSubmit={handleSave} style={{ display: 'grid', gap: 12, maxWidth: 560 }}>
        <fieldset style={{ border: '1px solid #eee', borderRadius: 8, padding: 12 }}>
          <legend>Account</legend>
          <label>
            <div>Email</div>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              style={{ width: '100%' }}
            />
          </label>
          <label>
            <div>Password</div>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="Enter current or new password"
              style={{ width: '100%' }}
            />
          </label>
        </fieldset>

        <fieldset style={{ border: '1px solid #eee', borderRadius: 8, padding: 12 }}>
          <legend>Display Preferences</legend>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
            <label>
              <div>Preferred Currency</div>
              <select
                value={preferredCurrency}
                onChange={(e) => setPreferredCurrency(e.target.value)}
                style={{ width: '100%' }}
              >
                {currencyOptions.map((c) => (
                  <option key={c} value={c}>
                    {c}
                  </option>
                ))}
              </select>
            </label>
            <label>
              <div>Decimal Places</div>
              <input
                type="number"
                value={decimalPlaces}
                min={0}
                max={4}
                onChange={(e) => setDecimalPlaces(Number(e.target.value))}
                style={{ width: '100%' }}
              />
            </label>
            <label>
              <div>Thousand Separator</div>
              <select
                value={thousandSeparator}
                onChange={(e) => setThousandSeparator(e.target.value)}
                style={{ width: '100%' }}
              >
                <option value=",">Comma (,)</option>
                <option value=".">Dot (.)</option>
                <option value=" ">Space ( )</option>
              </select>
            </label>
            <label>
              <div>Currency Symbol Placement</div>
              <select
                value={symbolPlacement}
                onChange={(e) => setSymbolPlacement(e.target.value as 'before' | 'after')}
                style={{ width: '100%' }}
              >
                <option value="before">Before (e.g., $1,234.56)</option>
                <option value="after">After (e.g., 1,234.56$)</option>
              </select>
            </label>
          </div>
        </fieldset>

        <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
          <button type="submit" disabled={saving}>
            {saving ? 'Saving…' : 'Save Changes'}
          </button>
          {error ? <span style={{ color: 'crimson' }}>{error}</span> : null}
          {success ? <span style={{ color: '#0a0' }}>{success}</span> : null}
        </div>
      </form>
    </div>
  )
}

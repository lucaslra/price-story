import { useEffect, useMemo, useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { get, post } from '@/api/client'
import { useAuth } from '@/auth/AuthContext'
import type { User } from '@/types'

type UsersResponse = { users: User[] }

export default function LoginPage() {
  const { user, login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const from = (location.state as { from?: Location })?.from?.pathname || '/'

  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [creating, setCreating] = useState(false)
  const [createError, setCreateError] = useState<string | null>(null)
  const [preferredCurrency, setPreferredCurrency] = useState('USD')
  const [decimalPlaces, setDecimalPlaces] = useState<number>(2)
  const [thousandSeparator, setThousandSeparator] = useState<string>(',')
  const [symbolPlacement, setSymbolPlacement] = useState<'before' | 'after'>('before')

  const currencyOptions = useMemo(
    () => ['USD', 'EUR', 'JPY', 'GBP', 'AUD', 'CAD', 'CHF', 'CNY'],
    []
  )

  useEffect(() => {
    // Adjust decimal places for JPY
    if (preferredCurrency === 'JPY') {
      setDecimalPlaces(0)
    } else if (decimalPlaces === 0) {
      setDecimalPlaces(2)
    }
  }, [preferredCurrency, decimalPlaces])

  useEffect(() => {
    let isMounted = true
    setLoading(true)
    get<UsersResponse>('/users')
      .then((resp) => {
        if (!isMounted) return
        setUsers(resp.users)
        setError(null)
      })
      .catch((e) => {
        if (!isMounted) return
        setError(e instanceof Error ? e.message : String(e))
      })
      .finally(() => {
        if (isMounted) setLoading(false)
      })
    return () => {
      isMounted = false
    }
  }, [])

  function handleSelect(u: User) {
    login(u)
    navigate(from, { replace: true })
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setCreating(true)
    setCreateError(null)
    try {
      if (!/^[A-Z]{3}$/.test(preferredCurrency)) {
        throw new Error('Preferred currency must be a 3-letter uppercase code')
      }
      if (preferredCurrency === 'JPY' && decimalPlaces !== 0) {
        throw new Error('JPY requires 0 decimal places')
      }
      if (decimalPlaces < 0 || decimalPlaces > 4) {
        throw new Error('Decimal places must be between 0 and 4')
      }
      if (![',', '.', ' '].includes(thousandSeparator)) {
        throw new Error('Thousand separator must be one of ",", ".", or space')
      }
      if (!['before', 'after'].includes(symbolPlacement)) {
        throw new Error('Currency symbol placement must be "before" or "after"')
      }

      const created = await post<
        {
          email: string
          password_hash: string
          preferred_currency: string
          decimal_places: number
          thousand_separator: string
          currency_symbol_placement: 'before' | 'after'
        },
        User
      >('/users', {
        email,
        password_hash: password,
        preferred_currency: preferredCurrency,
        decimal_places: decimalPlaces,
        thousand_separator: thousandSeparator,
        currency_symbol_placement: symbolPlacement
      })
      login(created)
      navigate(from, { replace: true })
    } catch (e) {
      setCreateError(e instanceof Error ? e.message : String(e))
    } finally {
      setCreating(false)
    }
  }

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <div style={{ marginBottom: 16 }}>
        <Link to="/">← Back to Home</Link>
      </div>
      <h1>Login</h1>
      {user ? <p style={{ color: '#0a0' }}>Currently signed in as {user.email}.</p> : null}

      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24 }}>
        <section>
          <h2 style={{ fontSize: 18 }}>Select existing user</h2>
          {loading ? (
            <p>Loading users…</p>
          ) : error ? (
            <p style={{ color: 'crimson' }}>Failed to load users: {error}</p>
          ) : users.length === 0 ? (
            <p>No users found.</p>
          ) : (
            <ul style={{ listStyle: 'none', padding: 0 }}>
              {users.map((u) => (
                <li
                  key={u.id}
                  style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}
                >
                  <span style={{ flex: 1 }}>{u.email}</span>
                  <button onClick={() => handleSelect(u)}>Sign in</button>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section>
          <h2 style={{ fontSize: 18 }}>Create new user</h2>
          <form onSubmit={handleCreate} style={{ display: 'grid', gap: 8, maxWidth: 480 }}>
            <label>
              <div>Email</div>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                placeholder="you@example.com"
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
                placeholder="••••••••"
                style={{ width: '100%' }}
              />
            </label>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
              <label>
                <div>Preferred currency</div>
                <select
                  value={preferredCurrency}
                  onChange={(e) => setPreferredCurrency(e.target.value)}
                  required
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
                <div>Decimal places</div>
                <input
                  type="number"
                  value={decimalPlaces}
                  onChange={(e) => setDecimalPlaces(Number(e.target.value))}
                  min={preferredCurrency === 'JPY' ? 0 : 0}
                  max={preferredCurrency === 'JPY' ? 0 : 4}
                  required
                  style={{ width: '100%' }}
                />
              </label>
            </div>
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 8 }}>
              <label>
                <div>Thousand separator</div>
                <select
                  value={thousandSeparator}
                  onChange={(e) => setThousandSeparator(e.target.value)}
                  required
                  style={{ width: '100%' }}
                >
                  <option value="," label="," />
                  <option value="." label="." />
                  <option value=" " label="Space" />
                </select>
              </label>
              <label>
                <div>Currency symbol placement</div>
                <select
                  value={symbolPlacement}
                  onChange={(e) => setSymbolPlacement(e.target.value as 'before' | 'after')}
                  required
                  style={{ width: '100%' }}
                >
                  <option value="before">Before</option>
                  <option value="after">After</option>
                </select>
              </label>
            </div>
            <div>
              <button type="submit" disabled={creating}>
                {creating ? 'Creating…' : 'Create & Sign in'}
              </button>
            </div>
            {createError ? <p style={{ color: 'crimson' }}>{createError}</p> : null}
          </form>
        </section>
      </div>
    </div>
  )
}

import { useEffect, useState } from 'react'
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
      const created = await post<{ email: string; password_hash: string }, User>('/users', {
        email,
        password_hash: password
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
          <form onSubmit={handleCreate} style={{ display: 'grid', gap: 8, maxWidth: 360 }}>
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

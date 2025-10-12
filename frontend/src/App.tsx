import { StrictMode, useState } from 'react'
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import ProductsPage from '@/pages/Products'
import PricePointsPage from '@/pages/PricePoints'
import PriceStoriesPage from '@/pages/PriceStories'
import PriceStoryDetailsPage from '@/pages/PriceStoryDetails'
import LoginPage from '@/pages/Login'
import NewPriceStoryPage from '@/pages/NewPriceStory'
import ProfilePage from '@/pages/Profile'
import { AuthProvider, useAuth } from '@/auth/AuthContext'
import RequireAuth from '@/auth/RequireAuth'
import { post } from '@/api/client'
import type { User } from '@/types'

function AppHeader() {
  const { user, logout } = useAuth()
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        gap: 12,
        padding: '8px 16px',
        borderBottom: '1px solid #eee',
        fontFamily: 'system-ui, sans-serif'
      }}
    >
      <Link to="/" style={{ fontWeight: 600 }}>
        Price Story
      </Link>
      {user ? (
        <nav style={{ display: 'flex', gap: 12 }}>
          <Link to="/products">Products</Link>
          <Link to="/price-points">Price Points</Link>
          <Link to="/price-stories">Price Stories</Link>
          <Link to="/profile">Profile</Link>
        </nav>
      ) : null}
      <div style={{ marginLeft: 'auto', display: 'flex', alignItems: 'center', gap: 8 }}>
        {user ? (
          <>
            <span style={{ color: '#555' }}>Signed in as {user.email}</span>
            <button onClick={logout}>Logout</button>
          </>
        ) : (
          <Link to="/login">Login</Link>
        )}
      </div>
    </div>
  )
}

function HomeContent() {
  const { user, login } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleRegister(e: React.FormEvent) {
    e.preventDefault()
    setCreating(true)
    setError(null)
    try {
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
        preferred_currency: 'USD',
        decimal_places: 2,
        thousand_separator: ',',
        currency_symbol_placement: 'before'
      })
      login(created)
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setCreating(false)
    }
  }

  if (user) {
    return (
      <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
        <h1>Your Price Stories</h1>
        <p style={{ color: '#555' }}>Browse and open your stories.</p>
        <PriceStoriesPage />
      </div>
    )
  }

  return (
    <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
      <h1>Welcome to Price Story</h1>
      <p style={{ color: '#555', maxWidth: 720 }}>
        Track products over time with price points and stories. Create a user to start tracking, or
        log in if you already have an account.
      </p>
      <div
        style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 24, alignItems: 'start' }}
      >
        <section>
          <h2 style={{ fontSize: 18 }}>Register</h2>
          <form onSubmit={handleRegister} style={{ display: 'grid', gap: 8, maxWidth: 360 }}>
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
                {creating ? 'Registering…' : 'Register & Continue'}
              </button>
            </div>
            {error ? <p style={{ color: 'crimson' }}>{error}</p> : null}
          </form>
        </section>
        <section>
          <h2 style={{ fontSize: 18 }}>Login</h2>
          <p>Already have an account? Sign in to access your stories.</p>
          <Link to="/login">Go to Login</Link>
        </section>
      </div>
    </div>
  )
}

function App() {
  return (
    <StrictMode>
      <AuthProvider>
        <BrowserRouter>
          <AppHeader />
          <Routes>
            <Route path="/" element={<HomeContent />} />
            <Route path="/products" element={<ProductsPage />} />
            <Route path="/price-points" element={<PricePointsPage />} />
            <Route path="/price-points/:productId" element={<PricePointsPage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/profile"
              element={
                <RequireAuth>
                  <ProfilePage />
                </RequireAuth>
              }
            />
            <Route
              path="/price-stories/new"
              element={
                <RequireAuth>
                  <NewPriceStoryPage />
                </RequireAuth>
              }
            />
            <Route
              path="/price-stories"
              element={
                <RequireAuth>
                  <PriceStoriesPage />
                </RequireAuth>
              }
            />
            <Route
              path="/price-stories/:id"
              element={
                <RequireAuth>
                  <PriceStoryDetailsPage />
                </RequireAuth>
              }
            />
          </Routes>
        </BrowserRouter>
      </AuthProvider>
    </StrictMode>
  )
}

export default App

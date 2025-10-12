import { StrictMode } from 'react'
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import ProductsPage from '@/pages/Products'
import PricePointsPage from '@/pages/PricePoints'
import PriceStoriesPage from '@/pages/PriceStories'
import PriceStoryDetailsPage from '@/pages/PriceStoryDetails'

function App() {
  return (
    <StrictMode>
      <BrowserRouter>
        <Routes>
          <Route
            path="/"
            element={
              <div style={{ padding: 24, fontFamily: 'system-ui, sans-serif' }}>
                <h1>Home</h1>
                <p style={{ color: '#555' }}>Scratch home page</p>
                <ul style={{ listStyle: 'none', padding: 0 }}>
                  <li>
                    <Link to="/products">Products</Link>
                  </li>
                  <li>
                    <Link to="/price-points">Price Points</Link>
                  </li>
                  <li>
                    <Link to="/price-stories">Price Stories</Link>
                  </li>
                </ul>
              </div>
            }
          />
          <Route path="/products" element={<ProductsPage />} />
          <Route path="/price-points" element={<PricePointsPage />} />
          <Route path="/price-points/:productId" element={<PricePointsPage />} />
          <Route path="/price-stories" element={<PriceStoriesPage />} />
          <Route path="/price-stories/:id" element={<PriceStoryDetailsPage />} />
        </Routes>
      </BrowserRouter>
    </StrictMode>
  )
}

export default App

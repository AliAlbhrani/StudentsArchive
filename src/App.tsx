import { Routes, Route, Navigate } from 'react-router-dom'
import Navbar from './components/Navbar'
import ProtectedRoute from './components/ProtectedRoute'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import ProfilePage from './pages/ProfilePage'
import PostsPage from './pages/PostsPage'
import MyPostsPage from './pages/MyPostsPage'
import CreatePostPage from './pages/CreatePostPage'
import { useAuth } from './context/AuthContext'

function HomePage() {
  const { user } = useAuth()
  return (
    <div className="page">
      <div className="hero-section">
        <h1>Students Archive</h1>
        <p>A platform for students to share and archive their academic work.</p>
        {user ? (
          <div className="hero-actions">
            <a href="/posts" className="btn btn-primary">Browse Posts</a>
            <a href="/create-post" className="btn btn-secondary">Create Post</a>
          </div>
        ) : (
          <div className="hero-actions">
            <a href="/login" className="btn btn-primary">Sign In</a>
            <a href="/register" className="btn btn-secondary">Get Started</a>
          </div>
        )}
      </div>
    </div>
  )
}

function App() {
  return (
    <div className="app">
      <Navbar />
      <main className="main-content">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/posts" element={<PostsPage />} />
          <Route element={<ProtectedRoute />}>
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/my-posts" element={<MyPostsPage />} />
            <Route path="/create-post" element={<CreatePostPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </main>
    </div>
  )
}

export default App

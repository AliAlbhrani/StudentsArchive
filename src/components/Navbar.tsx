import { Link, useLocation } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export default function Navbar() {
  const { user, logout } = useAuth()
  const location = useLocation()

  const isActive = (path: string) =>
    location.pathname === path ? 'nav-link active' : 'nav-link'

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <Link to="/">Students Archive</Link>
      </div>
      <div className="navbar-links">
        <Link to="/posts" className={isActive('/posts')}>
          Posts
        </Link>
        {user && (
          <>
            <Link to="/my-posts" className={isActive('/my-posts')}>
              My Posts
            </Link>
            <Link to="/create-post" className={isActive('/create-post')}>
              New Post
            </Link>
            <Link to="/profile" className={isActive('/profile')}>
              Profile
            </Link>
            <span className="navbar-user">@{user.username}</span>
            <button onClick={logout} className="btn btn-logout">
              Logout
            </button>
          </>
        )}
        {!user && (
          <>
            <Link to="/login" className={isActive('/login')}>
              Sign In
            </Link>
            <Link to="/register" className="btn btn-primary btn-sm">
              Register
            </Link>
          </>
        )}
      </div>
    </nav>
  )
}

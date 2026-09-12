import { useState, useEffect } from 'react'
import { getPosts } from '../api'
import type { Post } from '../types'

export default function PostsPage() {
  const [posts, setPosts] = useState<Post[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    getPosts()
      .then((data) => {
        setPosts(data.posts)
        setTotal(data.total)
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="loading">Loading posts...</div>

  return (
    <div className="page">
      <div className="page-header">
        <h1>All Posts</h1>
        <span className="badge">{total} total</span>
      </div>
      {error && <div className="error-msg">{error}</div>}
      {posts.length === 0 ? (
        <div className="empty-state">
          <p>No posts yet. Be the first to create one!</p>
        </div>
      ) : (
        <div className="posts-grid">
          {posts.map((post) => (
            <div key={post.id} className="post-card">
              <div className="post-card-header">
                <h3>{post.title}</h3>
                <span className="post-id">#{post.id}</span>
              </div>
              {post.content && <p className="post-content">{post.content}</p>}
              {post.images && post.images.length > 0 && (
                <div className="post-images">
                  {post.images.map((img, i) => (
                    <img key={i} src={img} alt="" className="post-image" />
                  ))}
                </div>
              )}
              <div className="post-footer">
                <span className="post-date">
                  {post.created_at
                    ? new Date(post.created_at).toLocaleDateString()
                    : 'Unknown date'}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

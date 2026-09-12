import { useState, useEffect, type FormEvent } from 'react'
import { useAuth } from '../context/AuthContext'
import { updateProfile } from '../api'

export default function ProfilePage() {
  const { user, loading } = useAuth()
  const [bio, setBio] = useState('')
  const [photoUrl, setPhotoUrl] = useState('')
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (user) {
      setBio(user.bio || '')
      setPhotoUrl(user.photo_url || '')
    }
  }, [user])

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setSaving(true)
    setError('')
    setMessage('')
    try {
      await updateProfile(bio || null, photoUrl || null)
      setMessage('Profile updated successfully')
    } catch {
      setError('Failed to update profile')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="loading">Loading...</div>
  if (!user) return <div className="error-msg">Not authenticated</div>

  return (
    <div className="page">
      <div className="page-header">
        <h1>My Profile</h1>
      </div>
      <div className="profile-view">
        <div className="profile-avatar">
          {user.photo_url ? (
            <img src={user.photo_url} alt={user.full_name} />
          ) : (
            <div className="avatar-placeholder">
              {user.full_name.charAt(0).toUpperCase()}
            </div>
          )}
        </div>
        <div className="profile-info">
          <div className="info-row">
            <span className="label">Name</span>
            <span className="value">{user.full_name}</span>
          </div>
          <div className="info-row">
            <span className="label">Username</span>
            <span className="value">@{user.username}</span>
          </div>
          <div className="info-row">
            <span className="label">Stage</span>
            <span className="value">{user.stage}</span>
          </div>
          <div className="info-row">
            <span className="label">Member since</span>
            <span className="value">
              {new Date(user.created_at).toLocaleDateString()}
            </span>
          </div>
        </div>
      </div>

      <div className="profile-edit">
        <h2>Edit Profile</h2>
        {message && <div className="success-msg">{message}</div>}
        {error && <div className="error-msg">{error}</div>}
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="bio">Bio</label>
            <textarea
              id="bio"
              value={bio}
              onChange={(e) => setBio(e.target.value)}
              placeholder="Tell us about yourself"
              rows={3}
            />
          </div>
          <div className="form-group">
            <label htmlFor="photoUrl">Photo URL</label>
            <input
              id="photoUrl"
              type="url"
              value={photoUrl}
              onChange={(e) => setPhotoUrl(e.target.value)}
              placeholder="https://example.com/photo.jpg"
            />
          </div>
          <button type="submit" className="btn btn-primary" disabled={saving}>
            {saving ? 'Saving...' : 'Save Changes'}
          </button>
        </form>
      </div>
    </div>
  )
}

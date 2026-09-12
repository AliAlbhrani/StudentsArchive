import type { TokensResponse, UserProfile, Post } from './types'

const API_BASE = ''

function getAuthHeaders(): Record<string, string> {
  const token = localStorage.getItem('access_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export async function register(
  full_name: string,
  username: string,
  password: string,
  stage: number
): Promise<TokensResponse> {
  const res = await fetch(`${API_BASE}/users/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ full_name, username, password, stage }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ detail: 'Request failed' }))
    throw new Error(err.detail || 'Registration failed')
  }
  return res.json()
}

export async function login(
  username: string,
  password: string
): Promise<TokensResponse> {
  const res = await fetch(`${API_BASE}/users/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ detail: 'Request failed' }))
    throw new Error(err.detail || 'Login failed')
  }
  return res.json()
}

export async function getProfile(): Promise<UserProfile> {
  const res = await fetch(`${API_BASE}/users/profile`, {
    headers: getAuthHeaders(),
  })
  if (!res.ok) throw new Error('Failed to fetch profile')
  return res.json()
}

export async function updateProfile(
  bio: string | null,
  photo_url: string | null
): Promise<string> {
  const res = await fetch(`${API_BASE}/users/profile`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ bio, photo_url }),
  })
  if (!res.ok) throw new Error('Failed to update profile')
  return res.json()
}

export async function getPosts(): Promise<{ posts: Post[]; total: number }> {
  const res = await fetch(`${API_BASE}/posts`, {
    headers: getAuthHeaders(),
  })
  if (res.status === 204 || res.ok) {
    const dataHeader = res.headers.get('Data')
    const totalHeader = res.headers.get('Total')
    const total = totalHeader ? parseInt(totalHeader, 10) : 0
    if (dataHeader) {
      const post = JSON.parse(dataHeader) as Post
      return { posts: [post], total }
    }
    try {
      const text = await res.text()
      if (text) {
        const posts = JSON.parse(text)
        return { posts: Array.isArray(posts) ? posts : [posts], total }
      }
    } catch {}
    return { posts: [], total }
  }
  throw new Error('Failed to fetch posts')
}

export async function createPost(
  title: string,
  content: string,
  images: string[] | null
): Promise<string> {
  const res = await fetch(`${API_BASE}/posts`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', ...getAuthHeaders() },
    body: JSON.stringify({ title, content, images }),
  })
  if (!res.ok) throw new Error('Failed to create post')
  return res.json()
}

export async function getMyPosts(): Promise<{
  posts: Post[]
  total: number
}> {
  const res = await fetch(`${API_BASE}/users/posts`, {
    headers: getAuthHeaders(),
  })
  if (res.status === 204 || res.ok) {
    const dataHeader = res.headers.get('Data')
    const totalHeader = res.headers.get('Total')
    const total = totalHeader ? parseInt(totalHeader, 10) : 0
    if (dataHeader) {
      const post = JSON.parse(dataHeader) as Post
      return { posts: [post], total }
    }
    try {
      const text = await res.text()
      if (text) {
        const posts = JSON.parse(text)
        return { posts: Array.isArray(posts) ? posts : [posts], total }
      }
    } catch {}
    return { posts: [], total }
  }
  throw new Error('Failed to fetch posts')
}

export async function uploadFile(file: File): Promise<string> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await fetch(`${API_BASE}/storage/upload`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: formData,
  })
  if (!res.ok) throw new Error('Failed to upload file')
  return res.json()
}

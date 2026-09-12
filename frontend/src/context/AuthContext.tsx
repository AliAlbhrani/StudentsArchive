import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import type { UserProfile } from '../types'
import { getProfile } from '../api'

interface AuthContextType {
  token: string | null
  user: UserProfile | null
  setToken: (token: string | null) => void
  logout: () => void
  loading: boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(
    localStorage.getItem('access_token')
  )
  const [user, setUser] = useState<UserProfile | null>(null)
  const [loading, setLoading] = useState(true)

  const setToken = (newToken: string | null) => {
    setTokenState(newToken)
    if (newToken) {
      localStorage.setItem('access_token', newToken)
    } else {
      localStorage.removeItem('access_token')
    }
  }

  const logout = () => {
    setToken(null)
    setUser(null)
  }

  useEffect(() => {
    if (token) {
      getProfile()
        .then(setUser)
        .catch(() => {
          setToken(null)
          setUser(null)
        })
        .finally(() => setLoading(false))
    } else {
      setUser(null)
      setLoading(false)
    }
  }, [token])

  return (
    <AuthContext.Provider value={{ token, user, setToken, logout, loading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

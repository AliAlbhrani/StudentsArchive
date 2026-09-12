export interface TokensResponse {
  access_token: string
  refresh_token: string
}

export interface UserProfile {
  UserID: number
  full_name: string
  username: string
  bio: string | null
  stage: number
  photo_url: string | null
  created_at: string
}

export interface Post {
  id: number
  user_id: number | null
  title: string
  content: string | null
  images: string[] | null
  created_at: string | null
  deleted_at: string | null
  updated_at: string | null
}

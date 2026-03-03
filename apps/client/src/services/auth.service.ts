import apiClient from './api'

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  password: string
  confirmPassword: string
}

export interface LoginResponse {
  data: {
    access_token: string
    refresh_token: string
  }
  message?: string
}

export interface RegisterResponse {
  data: Record<string, never>
  message?: string
}

export interface MeResponse {
  data: {
    user_id: string
    username: string
  }
  message?: string
}

export const authService = {
  async login(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/login', credentials)
    return response.data
  },

  async register(data: RegisterRequest): Promise<RegisterResponse> {
    const response = await apiClient.post<RegisterResponse>('/register', {confirm_password: data.confirmPassword, ...data})
    return response.data
  },

  async me(): Promise<MeResponse> {
    const response = await apiClient.get<MeResponse>('/me')
    return response.data
  },
}

export default authService

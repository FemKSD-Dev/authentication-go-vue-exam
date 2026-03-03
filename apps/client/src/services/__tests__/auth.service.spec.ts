import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { authService } from '../auth.service'
import apiClient from '../api'

vi.mock('../api')

describe('authService', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('login', () => {
    it('should call POST /login with credentials', async () => {
      const mockResponse = {
        data: {
          data: {
            access_token: 'mock-access-token',
            refresh_token: 'mock-refresh-token',
          },
          message: 'Login successful',
        },
      }

      vi.mocked(apiClient.post).mockResolvedValue(mockResponse)

      const result = await authService.login({
        username: 'testuser',
        password: 'password123',
      })

      expect(apiClient.post).toHaveBeenCalledWith('/login', {
        username: 'testuser',
        password: 'password123',
      })

      expect(result).toEqual(mockResponse.data)
      expect(result.data.access_token).toBe('mock-access-token')
      expect(result.data.refresh_token).toBe('mock-refresh-token')
    })

    it('should throw error on failed login', async () => {
      const mockError = new Error('Invalid credentials')
      vi.mocked(apiClient.post).mockRejectedValue(mockError)

      await expect(
        authService.login({
          username: 'testuser',
          password: 'wrongpassword',
        }),
      ).rejects.toThrow('Invalid credentials')
    })
  })

  describe('register', () => {
    it('should call POST /register with user data', async () => {
      const mockResponse = {
        data: {
          data: {},
          message: 'User registered successfully',
        },
      }

      vi.mocked(apiClient.post).mockResolvedValue(mockResponse)

      const result = await authService.register({
        username: 'newuser',
        password: 'password123',
        confirmPassword: 'password123',
      })

      expect(apiClient.post).toHaveBeenCalledWith('/register', {
        confirm_password: 'password123',
        username: 'newuser',
        password: 'password123',
        confirmPassword: 'password123',
      })

      expect(result).toEqual(mockResponse.data)
    })

    it('should throw error when registration fails', async () => {
      const mockError = new Error('Username already exists')
      vi.mocked(apiClient.post).mockRejectedValue(mockError)

      await expect(
        authService.register({
          username: 'existinguser',
          password: 'password123',
          confirmPassword: 'password123',
        }),
      ).rejects.toThrow('Username already exists')
    })
  })

  describe('me', () => {
    it('should call GET /me', async () => {
      const mockResponse = {
        data: {
          data: {
            user_id: '123',
            username: 'testuser',
          },
          message: 'User information retrieved successfully',
        },
      }

      vi.mocked(apiClient.get).mockResolvedValue(mockResponse)

      const result = await authService.me()

      expect(apiClient.get).toHaveBeenCalledWith('/me')
      expect(result).toEqual(mockResponse.data)
      expect(result.data.user_id).toBe('123')
      expect(result.data.username).toBe('testuser')
    })

    it('should throw error when not authenticated', async () => {
      const mockError = new Error('Unauthorized')
      vi.mocked(apiClient.get).mockRejectedValue(mockError)

      await expect(authService.me()).rejects.toThrow('Unauthorized')
    })
  })
})

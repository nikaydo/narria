import { useState } from 'react'
import { NarriaApi } from '../bindings/narria/backend/app'
import { UserData } from '../bindings/narria/backend/models'

function LoginPage() {
  const [isLogin, setIsLogin] = useState(true)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const resetForm = () => {
    setUsername('')
    setPassword('')
    setConfirmPassword('')
    setError('')
    setSuccess(false)
  }

  const handleToggleMode = () => {
    setIsLogin(!isLogin)
    resetForm()
  }

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)
    setLoading(true)

    try {
      const userData = new UserData({
        uuid: '00000000-0000-0000-0000-000000000000',
        username: username.trim(),
        password: password,
        recovery: ''
      })

      const result = await NarriaApi.AuthUser(userData)
      
      if (result && result.uuid) {
        setSuccess(true)
        setError('')
        console.log('Login successful:', result)
      } else {
        setError('Неверный логин или пароль')
      }
    } catch (err: any) {
      console.error('Login error:', err)
      const errorMessage = err.message || 'Ошибка при входе. Проверьте логин и пароль.'
      setError(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)

    // Валидация
    if (password.length < 6) {
      setError('Пароль должен содержать минимум 6 символов')
      return
    }

    if (password !== confirmPassword) {
      setError('Пароли не совпадают')
      return
    }

    if (username.trim().length < 3) {
      setError('Логин должен содержать минимум 3 символа')
      return
    }

    setLoading(true)

    try {
      const userData = new UserData({
        uuid: '00000000-0000-0000-0000-000000000000',
        username: username.trim(),
        password: password,
        recovery: ''
      })

      const userUUID = await NarriaApi.CreateUser(userData)
      
      if (userUUID) {
        setSuccess(true)
        setError('')
        console.log('Registration successful:', userUUID)
        // Автоматически переключаемся на форму входа после успешной регистрации
        setTimeout(() => {
          setIsLogin(true)
          resetForm()
        }, 2000)
      } else {
        setError('Ошибка при регистрации')
      }
    } catch (err: any) {
      console.error('Registration error:', err)
      let errorMessage = 'Ошибка при регистрации'
      
      // Проверяем, не занят ли логин
      if (err.message && err.message.includes('UNIQUE constraint')) {
        errorMessage = 'Пользователь с таким логином уже существует'
      } else if (err.message) {
        errorMessage = err.message
      }
      
      setError(errorMessage)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-container">
      <div className="login-box">
        <div className="auth-header">
          <h1>{isLogin ? 'Вход в систему' : 'Регистрация'}</h1>
          <button 
            type="button"
            onClick={handleToggleMode}
            className="toggle-button"
            disabled={loading}
          >
            {isLogin ? 'Нет аккаунта? Зарегистрироваться' : 'Уже есть аккаунт? Войти'}
          </button>
        </div>
        
        <form onSubmit={isLogin ? handleLogin : handleRegister}>
          <div className="form-group">
            <label htmlFor="username">Логин:</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Введите логин"
              required
              disabled={loading}
              className="input"
              minLength={3}
            />
          </div>
          
          <div className="form-group">
            <label htmlFor="password">Пароль:</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder={isLogin ? "Введите пароль" : "Введите пароль (минимум 6 символов)"}
              required
              disabled={loading}
              className="input"
              minLength={isLogin ? undefined : 6}
            />
          </div>

          {!isLogin && (
            <div className="form-group">
              <label htmlFor="confirmPassword">Подтверждение пароля:</label>
              <input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                placeholder="Повторите пароль"
                required
                disabled={loading}
                className="input"
                minLength={6}
              />
            </div>
          )}

          {error && <div className="error-message">{error}</div>}
          {success && (
            <div className="success-message">
              {isLogin ? 'Вход выполнен успешно!' : 'Регистрация выполнена успешно! Переход на страницу входа...'}
            </div>
          )}
          
          <button 
            type="submit" 
            disabled={loading}
            className="login-button"
          >
            {loading 
              ? (isLogin ? 'Вход...' : 'Регистрация...') 
              : (isLogin ? 'Войти' : 'Зарегистрироваться')
            }
          </button>
        </form>
      </div>
    </div>
  )
}

export default LoginPage


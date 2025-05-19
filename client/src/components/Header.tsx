import { Link } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/lib/auth'
import { Moon, Sun } from 'lucide-react'
import { useTheme } from '@/lib/theme'

export default function Header() {
  const { isAuthenticated, logout } = useAuthStore()
  const { theme, toggleTheme } = useTheme()

  return (
    <header className="header">
      <div className="header-container">
        <div className="header-nav">
          <Link to="/" className="header-logo">
            <span className="header-logo-text">Todo App</span>
          </Link>
        </div>
        <div className="flex-1" />
        <Button
          variant="ghost"
          size="icon"
          onClick={toggleTheme}
          className="mr-2"
        >
          {theme === 'dark' ? (
            <Sun className="h-5 w-5" />
          ) : (
            <Moon className="h-5 w-5" />
          )}
        </Button>
        {isAuthenticated && (
          <Button variant="ghost" onClick={logout}>
            Logout
          </Button>
        )}
      </div>
    </header>
  )
}
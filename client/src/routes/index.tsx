import { createFileRoute } from '@tanstack/react-router'
import { useAuthStore } from '@/lib/auth'
import { TodoList } from '@/components/todo/todo-list'
import { Dialog, DialogContent, DialogTrigger } from "@/components/ui/dialog"
import { Button } from '@/components/ui/button'
import { LoginForm } from '@/components/auth/login-form'
import { RegisterForm } from '@/components/auth/register-form'


export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  const { isAuthenticated } = useAuthStore()

  return (
    <div className="app-container">
      <main className="app-main">
        {isAuthenticated ? (
          <TodoList />
        ) : (
          <div className="auth-buttons-container">
            <Dialog>
              <DialogTrigger asChild>
                <Button variant="outline">Login</Button>
              </DialogTrigger>
              <DialogContent>
                <LoginForm />
              </DialogContent>
            </Dialog>

            <Dialog>
              <DialogTrigger asChild>
                <Button>Register</Button>
              </DialogTrigger>
              <DialogContent>
                <RegisterForm />
              </DialogContent>
            </Dialog>
          </div>
        )}
      </main>
    </div>
  )
}

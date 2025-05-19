import { Outlet, createRootRoute } from '@tanstack/react-router'
import { Toaster } from 'sonner'
import Header from '../components/Header'
import { ThemeProvider } from '@/lib/theme'

export const Route = createRootRoute({
  component: () => (
    <ThemeProvider>
      <Toaster richColors />
      <Header />
      <Outlet />
    </ThemeProvider>
  ),
})

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/lib/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { TodoItem } from './todo-item'
import { toast } from "sonner"
import { Loader2 } from "lucide-react"

type FilterType = 'all' | 'active' | 'completed'

interface Todo {
  id: string
  title: string
  description: string
  completed: boolean
}

export function TodoList() {
  const [todos, setTodos] = useState<Todo[]>([]) // Initialize with empty array
  const [isLoadingTodos, setIsLoadingTodos] = useState(false)
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [filter, setFilter] = useState<FilterType>('all')
  const { token } = useAuthStore()

  const fetchTodos = async () => {
    setIsLoadingTodos(true)
    try {
      const res = await fetch('http://localhost:8080/api/v1/todos', { 
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!res.ok) throw new Error('Failed to fetch todos')
      const data = await res.json()
      setTodos(data || []) // Ensure we always set an array
    } catch (error) {
      const json = {
        variant: "destructive",
        title: "Error",
        description: error instanceof Error? error.message : "Something went wrong",
      }
      toast(JSON.stringify(json, null, 2))
      setTodos([]) // Set empty array on error
    } finally {
      setIsLoadingTodos(false)
    }
  }

  const addTodo = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    try {
      const res = await fetch('http://localhost:8080/api/v1/todos', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ title, description }),
      })
      if (!res.ok) throw new Error('Failed to add todo')
      const newTodo = await res.json()
      setTodos([...todos, newTodo])
      setTitle('')
      setDescription('')
    } catch (error) {
      const json = {
        variant: "destructive",
        title: "Error",
        description: "Failed to add todo",
      }
      toast(JSON.stringify(json, null, 2))
    } finally {
      setIsLoading(false)
    }
  }

  const onDelete = async (id: string) => {
    try {
      const res = await fetch(`http://localhost:8080/api/v1/todos/${id}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
      if (!res.ok) throw new Error('Failed to delete todo')
      setTodos(todos.filter(todo => todo.id !== id))
    } catch (error) {
      const json = {
        variant: "destructive",
        title: "Error",
        description: error instanceof Error ? error.message : "Failed to delete todo",
      }
      toast(JSON.stringify(json, null, 2))
    }
  }

  const onToggle = async (id: string, completed: boolean) => {
    try {
      const res = await fetch(`http://localhost:8080/api/v1/todos/${id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ completed }),
      })
      if (!res.ok) throw new Error('Failed to update todo')
      setTodos(todos.map(todo => 
        todo.id === id ? { ...todo, completed } : todo
      ))
    } catch (error) {
      const json = {
        variant: "destructive",
        title: "Error",
        description: error instanceof Error ? error.message : "Failed to update todo",
      }
      toast(JSON.stringify(json, null, 2))
    }
  }

  useEffect(() => {
    if (token) {
      fetchTodos()
    }
  }, [token])

  const filteredTodos = (todos || []).filter(todo => {
    switch (filter) {
      case 'active':
        return !todo.completed
      case 'completed':
        return todo.completed
      default:
        return true
    }
  })

  return (
    <div className="todo-container">
      <div className="flex justify-center space-x-4 mb-6">
        <Button
          variant={filter === 'all' ? 'default' : 'outline'}
          onClick={() => setFilter('all')}
        >
          All
        </Button>
        <Button
          variant={filter === 'active' ? 'default' : 'outline'}
          onClick={() => setFilter('active')}
        >
          Active
        </Button>
        <Button
          variant={filter === 'completed' ? 'default' : 'outline'}
          onClick={() => setFilter('completed')}
        >
          Completed
        </Button>
      </div>

      <form onSubmit={addTodo} className="todo-form">
        <h2 className="auth-title">Add New Todo</h2>
        <Input
          placeholder="Todo title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="form-input"
          disabled={isLoading}
        />
        <Input
          placeholder="Description"
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          disabled={isLoading}
        />
        <Button type="submit" className="primary-button" disabled={isLoading}>
          {isLoading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Adding...
            </>
          ) : (
            'Add Todo'
          )}
        </Button>
      </form>

      <div className="todo-list">
        {isLoadingTodos ? (
          <div className="flex justify-center">
            <Loader2 className="h-6 w-6 animate-spin" />
          </div>
        ) : (
          filteredTodos.map((todo) => (
            <TodoItem key={todo.id} todo={todo} onDelete={onDelete} onToggle={onToggle} />
          ))
        )}
      </div>
    </div>
  )
}
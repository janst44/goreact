import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Trash2 } from "lucide-react"

interface Todo {
  id: string
  title: string
  description: string
  completed: boolean
}

interface TodoItemProps {
  todo: Todo
  onDelete: (id: string) => void
  onToggle: (id: string, completed: boolean) => void
}

export function TodoItem({ todo, onDelete, onToggle }: TodoItemProps) {
  return (
    <Card>
      <CardContent className="todo-item-content">
        <div className="todo-item-left">
          <input
            type="checkbox"
            checked={todo.completed}
            onChange={(e) => onToggle(todo.id, e.target.checked)}
            className="todo-checkbox"
          />
          <div>
            <p className={`todo-item-text ${todo.completed ? 'todo-item-text-completed' : ''}`}>
              {todo.title}
            </p>
            {todo.description && (
              <p className="todo-item-description">{todo.description}</p>
            )}
          </div>
        </div>
        <Button
          variant="ghost"
          size="icon"
          onClick={() => onDelete(todo.id)}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </CardContent>
    </Card>
  )
}

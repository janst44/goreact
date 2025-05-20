import { Checkbox } from "@/components/ui/checkbox"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Trash2, Pencil, Save } from "lucide-react"
import { useState } from "react"

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
  onEdit: (id: string, title: string, description: string) => void
}

export function TodoItem({ todo, onDelete, onToggle, onEdit }: TodoItemProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [editTitle, setEditTitle] = useState(todo.title)
  const [editDescription, setEditDescription] = useState(todo.description)

  const handleSave = () => {
    onEdit(todo.id, editTitle, editDescription)
    setIsEditing(false)
  }

  return (
    <Card>
      <CardContent className="todo-item-content">
        <div className="todo-item-left">
          <Checkbox
            checked={todo.completed}
            onCheckedChange={(checked) => onToggle(todo.id, checked as boolean)}
            className="todo-checkbox"
          />
          {isEditing ? (
            <div className="flex flex-col gap-2 flex-1">
              <Input
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                placeholder="Todo title"
              />
              <Input
                value={editDescription}
                onChange={(e) => setEditDescription(e.target.value)}
                placeholder="Description"
              />
            </div>
          ) : (
            <div>
              <p className={`todo-item-text ${todo.completed ? 'todo-item-text-completed' : ''}`}>
                {todo.title}
              </p>
              {todo.description && (
                <p className="todo-item-description">{todo.description}</p>
              )}
            </div>
          )}
        </div>
        <div className="flex gap-2">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => isEditing ? handleSave() : setIsEditing(true)}
          >
            {isEditing ? (
              <Save className="h-4 w-4" />
            ) : (
              <Pencil className="h-4 w-4" />
            )}
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(todo.id)}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </CardContent>
    </Card>
  )
}

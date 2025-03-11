'use client';

import { useState, useEffect } from 'react';
import { Todo } from '../types/todo';

export default function TodoList() {
    const [todos, setTodos] = useState<Todo[]>([]);
    const [newTodo, setNewTodo] = useState('');
    const [ws, setWs] = useState<WebSocket | null>(null);

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        
        const response = await fetch('http://localhost:8080/api/todos', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title: newTodo }),
        });

        if (response.ok) {
            setNewTodo('');
        }
    };

    async function toggleTodo(id: string) {
        const response = await fetch(`http://localhost:8080/api/todos/${id}/toggle`, {
            method: 'PUT',
        });

        if (response.ok) {
            const updatedTodo = await response.json();
            setTodos(prev => prev.map(todo => 
                todo.id === id ? updatedTodo : todo
            ));
        }
    };

    useEffect(() => {
        fetch('http://localhost:8080/api/todos')
            .then(res => res.json())
            .then(data => setTodos(data));

        const websocket = new WebSocket('ws://localhost:8080/ws');
        websocket.onmessage = (event) => {
            const newTodo = JSON.parse(event.data);
            setTodos(prev => {
                const exists = prev.some(todo => todo.id === newTodo.id);
                if (exists) {
                    return prev.map(todo => 
                        todo.id === newTodo.id ? newTodo : todo
                    );
                }
                return [...prev, newTodo];
            });
        };

        setWs(websocket);

        return () => websocket.close();
    }, []);

    return (
        <div className="max-w-3xl mx-auto mt-10 p-8 bg-gray-900 rounded-xl shadow-2xl">
            <h1 className="text-3xl font-bold mb-8 text-white bg-gradient-to-r from-purple-500 to-pink-500 bg-clip-text text-transparent">
                Minhas Tarefas
            </h1>
            
            <form onSubmit={handleSubmit} className="mb-8">
                <div className="flex gap-3">
                    <input
                        type="text"
                        value={newTodo}
                        onChange={(e) => setNewTodo(e.target.value)}
                        className="flex-1 px-4 py-3 bg-gray-800 border border-gray-700 rounded-lg 
                                 text-gray-100 placeholder-gray-500 focus:outline-none 
                                 focus:ring-2 focus:ring-purple-500 focus:border-transparent
                                 transition-all duration-200"
                        placeholder="Adicionar nova tarefa..."
                    />
                    <button 
                        type="submit"
                        className="px-6 py-3 bg-gradient-to-r from-purple-500 to-pink-500 
                                 text-white rounded-lg font-medium hover:opacity-90 
                                 transform hover:scale-105 transition-all duration-200
                                 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    >
                        Adicionar
                    </button>
                </div>
            </form>

            <ul className="space-y-3">
                {todos.map(todo => (
                    <li 
                        key={todo.id}
                        className={`flex items-center gap-3 p-4 bg-gray-800 rounded-lg
                                  border border-gray-700 transform transition-all duration-200
                                  hover:border-purple-500 ${
                                    todo.completed ? 'opacity-60' : ''
                                  }`}
                    >
                        <input 
                            type="checkbox" 
                            checked={todo.completed}
                            onChange={() => toggleTodo(todo.id)}
                            className="w-5 h-5 rounded-md border-gray-600 text-purple-500 
                                     focus:ring-purple-500 focus:ring-offset-gray-800
                                     cursor-pointer transition-all duration-200"
                        />
                        <span className={`text-gray-100 flex-1 text-lg ${
                            todo.completed ? 'line-through text-gray-500' : ''
                        }`}>
                            {todo.title}
                        </span>
                        <span className="text-xs text-gray-500">
                            {new Date(todo.created_at).toLocaleDateString()}
                        </span>
                    </li>
                ))}
            </ul>

            {todos.length === 0 && (
                <div className="text-center py-10 text-gray-500">
                    <p>Nenhuma tarefa ainda. Adicione uma!</p>
                </div>
            )}
        </div>
    );
}
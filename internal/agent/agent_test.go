package agent

import (
	"github.com/Yorshik/Go-Final/internal/ast"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected float64
	}{
		{
			name:     "Addition",
			task:     Task{Arg1: 2, Arg2: 3, Operation: "+", OperationTime: 0},
			expected: 5,
		},
		{
			name:     "Subtraction",
			task:     Task{Arg1: 5, Arg2: 3, Operation: "-", OperationTime: 0},
			expected: 2,
		},
		{
			name:     "Multiplication",
			task:     Task{Arg1: 4, Arg2: 3, Operation: "*", OperationTime: 0},
			expected: 12,
		},
		{
			name:     "Division",
			task:     Task{Arg1: 6, Arg2: 2, Operation: "/", OperationTime: 0},
			expected: 3,
		},
		{
			name:     "Unknown operation",
			task:     Task{Arg1: 2, Arg2: 3, Operation: "%", OperationTime: 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compute(tt.task)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetOperationTime(t *testing.T) {
	tests := []struct {
		name       string
		op         string
		envKey     string
		envValue   string
		defaultVal int
		expected   int
	}{
		{
			name:       "Addition with env",
			op:         "+",
			envKey:     "TIME_ADDITION_MS",
			envValue:   "500",
			defaultVal: 1000,
			expected:   500,
		},
		{
			name:       "Subtraction without env",
			op:         "-",
			envKey:     "TIME_SUBTRACTION_MS",
			defaultVal: 1000,
			expected:   1000,
		},
		{
			name:       "Unknown operation",
			op:         "%",
			defaultVal: 1000,
			expected:   1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" && tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}
			s := &Server{}
			result := s.getOperationTime(tt.op)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEvaluateNode(t *testing.T) {
	s := &Server{
		Tasks:   make(chan Task, 100),
		Results: make(map[int]float64),
		mu:      sync.Mutex{},
	}
	go s.Worker() // Запускаем воркер для обработки задач

	tests := []struct {
		name     string
		node     *ast.Node
		expected float64
	}{
		{
			name:     "Simple number",
			node:     &ast.Node{Value: 42},
			expected: 42,
		},
		{
			name: "Addition",
			node: &ast.Node{
				Operator: "+",
				Left:     &ast.Node{Value: 2},
				Right:    &ast.Node{Value: 3},
			},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr := &expression{Node: tt.node}
			result := s.EvaluateNode(tt.node, expr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

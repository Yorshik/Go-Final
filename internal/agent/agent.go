package agent

import (
	"context"
	"fmt"
	"github.com/Yorshik/Go-Final/internal/ast"
	agentpb "github.com/Yorshik/Go-Final/internal/proto/gen"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}
type expression struct {
	ID     string    `json:"id"`
	Status string    `json:"status"`
	Result *float64  `json:"result"`
	Node   *ast.Node `json:"-"`
	Tasks  []Task    `json:"-"`
}

type Server struct {
	agentpb.UnimplementedAgentServer
	Tasks   chan Task
	Results map[int]float64
	taskID  int
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func (s *Server) getOperationTime(op string) int {
	switch op {
	case "+":
		return getEnvInt("TIME_ADDITION_MS", 1000)
	case "-":
		return getEnvInt("TIME_SUBTRACTION_MS", 1000)
	case "*":
		return getEnvInt("TIME_MULTIPLICATIONS_MS", 1000)
	case "/":
		return getEnvInt("TIME_DIVISIONS_MS", 1000)
	default:
		return 1000
	}
}

func compute(t Task) float64 {
	time.Sleep(time.Duration(t.OperationTime) * time.Millisecond)
	switch t.Operation {
	case "+":
		return t.Arg1 + t.Arg2
	case "-":
		return t.Arg1 - t.Arg2
	case "*":
		return t.Arg1 * t.Arg2
	case "/":
		return t.Arg1 / t.Arg2
	default:
		return 0
	}
}

func (s *Server) Worker() {
	for t := range s.Tasks {
		result := compute(t)
		s.mu.Lock()
		s.Results[t.ID] = result
		s.mu.Unlock()
	}
}

func (s *Server) EvaluateNode(node *ast.Node, expr *expression) float64 {
	if node.Operator == "" {
		return node.Value
	}
	var leftValue, rightValue float64
	if node.Left.Operator != "" {
		leftValue = s.EvaluateNode(node.Left, expr)
	} else {
		leftValue = node.Left.Value
	}
	if node.Right.Operator != "" {
		rightValue = s.EvaluateNode(node.Right, expr)
	} else {
		rightValue = node.Right.Value
	}
	s.mu.Lock()
	s.taskID++
	task := Task{
		ID:            s.taskID,
		Arg1:          leftValue,
		Arg2:          rightValue,
		Operation:     node.Operator,
		OperationTime: s.getOperationTime(node.Operator),
	}
	expr.Tasks = append(expr.Tasks, task)
	s.mu.Unlock()

	s.Tasks <- task

	for {
		s.mu.Lock()
		if result, ok := s.Results[task.ID]; ok {
			delete(s.Results, s.taskID)
			s.mu.Unlock()
			return result
		}
		s.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *Server) SendExpression(ctx context.Context, req *agentpb.ExpressionRequest) (*agentpb.ExpressionResult, error) {
	reqExpression := req.GetExpression()
	node, err := ast.Parse(reqExpression)
	if err != nil {
		return nil, err
	}
	expr := &expression{ID: req.GetId(), Status: "pending", Node: node}

	result := s.EvaluateNode(expr.Node, expr)
	return &agentpb.ExpressionResult{
		Id:     req.GetId(),
		Result: fmt.Sprintf("%f", result),
	}, nil
}

func StartServer() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при запуске gRPC-сервера: %v", err)
	}
	s := &Server{
		Tasks:   make(chan Task, 100),
		Results: make(map[int]float64),
	}
	power := getEnvInt("COMPUTING_POWER", 1)
	if power <= 0 {
		power = 1
	}
	for i := 0; i < power; i++ {
		go s.Worker()
	}
	grpcServer := grpc.NewServer()
	agentpb.RegisterAgentServer(grpcServer, s)
	fmt.Println("Запуск gRPC-сервера на порту 50051")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Ошибка при запуске gRPC-сервера: %v", err)
	}
}

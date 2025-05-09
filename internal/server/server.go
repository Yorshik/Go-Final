package server

import (
	"fmt"
	auth "github.com/Yorshik/Go-Final/internal/auth"
	"github.com/Yorshik/Go-Final/internal/database"
	"github.com/Yorshik/Go-Final/internal/models"
	agentpb "github.com/Yorshik/Go-Final/internal/proto/gen"
	"github.com/Yorshik/Go-Final/internal/utils"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ExpressionResponse struct {
	ID         uint   `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
	UserID     uint   `json:"user_id"`
}

func RegisterHandler(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		log.Println("Error binding input:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}
	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "User already exists"})
	}
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to hash password"})
	}
	user.Password = hashedPassword
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to register user"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Registration successful"})
}

func LoginHandler(c echo.Context) error {
	var credentials models.User
	if err := c.Bind(&credentials); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}
	var user models.User
	if err := database.DB.Where("username = ?", credentials.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid credentials"})
	}
	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid credentials"})
	}
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to generate token"})
	}
	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func CalculateHandler(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		fmt.Println("Missing token in Authorization header")
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Missing token"})
	}
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		fmt.Printf("Invalid token: %v\n", err)
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token"})
	}
	fmt.Printf("User ID from token: %d\n", claims.ID)
	var expression models.Expression
	if err := c.Bind(&expression); err != nil {
		fmt.Printf("Error binding expression data: %v\n", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid input"})
	}
	fmt.Printf("Received expression: %s\n", expression.Expression)
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect to gRPC server: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to connect to gRPC server"})
	}
	defer conn.Close()
	client := agentpb.NewAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	fmt.Printf("Sending expression to gRPC server: %s\n", expression.Expression)
	resp, err := client.SendExpression(ctx, &agentpb.ExpressionRequest{Expression: expression.Expression})
	if err != nil {
		fmt.Printf("Failed to send expression: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to send expression"})
	}
	fmt.Printf("Received result from gRPC server: %s\n", resp.Result)
	expression.Result = resp.Result
	expression.UserID = claims.ID
	if err := database.DB.Create(&expression).Error; err != nil {
		fmt.Printf("Failed to save expression to database: %v\n", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to save expression"})
	}
	fmt.Println("Expression calculated and saved successfully")
	return c.JSON(http.StatusOK, echo.Map{"message": "Expression calculated", "result": expression.Result})
}

func GetAllExpressionsHandler(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Missing token"})
	}
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token"})
	}
	userID := claims.ID
	var expressions []models.Expression
	if err := database.DB.Where("user_id = ?", userID).Find(&expressions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to retrieve expressions"})
	}
	var response []ExpressionResponse
	for _, expr := range expressions {
		response = append(response, ExpressionResponse{
			ID:         expr.ID,
			Expression: expr.Expression,
			Result:     expr.Result,
			UserID:     expr.UserID,
		})
	}

	return c.JSON(http.StatusOK, response)
}

func GetExpressionByIDHandler(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Missing token"})
	}
	claims, err := auth.ValidateToken(tokenString)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "Invalid token"})
	}
	userID := claims.ID
	idParam := c.Param("id")
	exprID, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid expression ID"})
	}
	var expression models.Expression
	if err := database.DB.Where("id = ? AND user_id = ?", exprID, userID).First(&expression).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "Expression not found"})
	}
	response := ExpressionResponse{
		ID:         expression.ID,
		Expression: expression.Expression,
		Result:     expression.Result,
		UserID:     expression.UserID,
	}
	return c.JSON(http.StatusOK, response)
}

func StartServer() {
	auth.InitJWT()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	database.ConnectDatabase()
	e := echo.New()
	v1 := e.Group("/api/v1")
	v1.POST("/register", RegisterHandler)
	v1.POST("/login", LoginHandler)
	v1.POST("/calculate", CalculateHandler)
	v1.GET("/expressions", GetAllExpressionsHandler)
	v1.GET("/expression/:id", GetExpressionByIDHandler)
	e.Logger.Fatal(e.Start(":8080"))

}

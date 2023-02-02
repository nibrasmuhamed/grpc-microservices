package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"login-service/protobuf"
	"net"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var db *sql.DB

type server struct {
	protobuf.UnimplementedRegServer
}

func (*server) Login(ctx context.Context, r *protobuf.LoginRequest) (*protobuf.LoginResponse, error) {
	email := r.GetEmail()
	password := r.GetPassword()

	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	query := `select email, password from users where email = $1`
	row := db.QueryRow(query, email)
	err := row.Scan(&payload.Email, &payload.Password)
	if err != nil {
		log.Println("error while reading rows: ", err)
		return nil, errors.New("error while reading rows")
	}
	valid, _ := PasswordMatches(password, payload.Password)
	if !valid {
		log.Println("password incorrect: ")
		return nil, errors.New("incorrect creadentials")
	}
	return &protobuf.LoginResponse{
		Email:  payload.Email,
		Status: "you are logged in",
	}, nil
}

func main() {
	fmt.Println("Login Microservice is starting. gRPC Server Powering up.")

	db = ConnectToDB()

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error while setting up a listener %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	protobuf.RegisterRegServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error while starting the Server %v", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	fmt.Println(host, password)
	dsn := fmt.Sprintf("host=%s port=5432 user=postgres password=%s dbname=postgres sslmode=disable timezone=UTC connect_timeout=5", host, password)
	count := 0
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres Not Yet Started")
			count++
		} else {
			log.Println("Connected to postgres")
			return connection
		}
		if count > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two seconds")
		time.Sleep(2 * time.Second)
		continue
	}
}

func PasswordMatches(plainText, Password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

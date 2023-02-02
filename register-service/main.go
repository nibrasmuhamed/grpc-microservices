package main

import (
	"api-service/protobuf"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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

func (*server) SignUp(ctx context.Context, r *protobuf.Request) (*protobuf.Response, error) {
	name := r.GetName()
	email := r.GetEmail()
	password := r.GetPassword()
	query := `insert into users (email, name, password) values ($1,$2,$3)`
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error while hashing")
		return nil, errors.New("error while hashing")
	}
	_, err = db.Exec(query, email, name, hashedPassword)
	if err != nil {
		log.Println("error while inserting", err)
		return nil, errors.New("error while inserting")
	}
	return &protobuf.Response{
		Name:  name,
		Email: email,
		Data:  "success",
	}, nil
}

func main() {
	fmt.Println("Starting gRPC Server")
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

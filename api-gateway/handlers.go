package main

import (
	"api-service/protobuf"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"
)

func (a *Config) login(w http.ResponseWriter, r *http.Request) {
	var loginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginPayload)
	if err != nil {
		a.errorJson(w, errors.New("cannot decode json"), 400)
		return
	}
	cc, err := grpc.Dial("login-service:50051", grpc.WithInsecure())
	if err != nil {
		a.errorJson(w, errors.New("internal server error"), 500)
		return
	}
	defer cc.Close()
	c := protobuf.NewRegClient(cc)
	req := &protobuf.LoginRequest{
		Email:    loginPayload.Email,
		Password: loginPayload.Password,
	}
	res, err := c.Login(context.Background(), req)
	if err != nil {
		fmt.Println("error is :", err)
		a.errorJson(w, errors.New("internal server error"), 500)
		return
	}
	log.Println(res)
	var resp struct {
		Email string `json:"email"`
		Data  string `json:"data"`
	}
	resp.Email = res.Email
	resp.Data = res.Status
	a.writeJson(w, 200, resp)
}

func (a *Config) register(w http.ResponseWriter, r *http.Request) {
	var registerpayload struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&registerpayload)
	if err != nil {
		a.errorJson(w, errors.New("cannot decode json"), 400)
		return
	}
	cc, err := grpc.Dial("register-service:50051", grpc.WithInsecure())
	if err != nil {
		a.errorJson(w, errors.New("internal server error"), 500)
		return
	}
	defer cc.Close()
	c := protobuf.NewRegClient(cc)
	req := &protobuf.Request{
		Name:     registerpayload.Name,
		Email:    registerpayload.Email,
		Password: registerpayload.Password,
	}
	res, err := c.SignUp(context.Background(), req)
	if err != nil {
		fmt.Println("error is :", err)
		a.errorJson(w, errors.New("internal server error"), 500)
	}
	var resp struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Data  string `json:"data"`
	}
	fmt.Println(res)
	resp.Data = res.Data
	resp.Email = res.Email
	resp.Name = res.Name
	a.writeJson(w, 200, resp)
}

/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"

	pb "grpc_test/helloworld/helloworld"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"hoge.net/cancelablelistener"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	name := in.GetName()

	log.Printf("Received: %v", name)
	switch name {
	case "error":
		return nil, status.Errorf(codes.Unavailable, "Forced Error")
	case "error with details":
		st := status.New(codes.Unavailable, "Forced Error (with details)")
		details := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{Field: "Name", Description: "Name is Forced Error"},
			},
		}
		st, _ = st.WithDetails(details)
		return nil, st.Err()
	default:
		return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
	}
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("%v", err)
	}
}

func run(ctx context.Context) error {
	flag.Parse()
	lis, err := cancelablelistener.Listen(ctx, "tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to listen: %v", err))
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			lis.Close()
			lis.Close()
		}
	}(ctx)
	if err := s.Serve(lis); err != nil {
		return errors.New(fmt.Sprintf("failed to serve: %v", err))
	}
	return nil
}

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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "grpc_test/helloworld/helloworld"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	SendUnary(c)
}

func SendUnary(c pb.GreeterClient) {
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		// Conversion for transform Error to Status
		st, ok := status.FromError(err)
		if !ok {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("GRPC Error :  Coce [%d], Message [%s]", st.Code(), st.Message())
		if len(st.Details()) > 0 {
			for _, detail := range st.Details() {
				switch d := detail.(type) {
				case *errdetails.BadRequest:
					log.Fatalf("Details: BadRequest: %v", d)
				case *errdetails.DebugInfo:
					log.Fatalf("Details: DebugInfo: %v", d)
				default:
					log.Printf("Details: Unknown: %v", d)
				}
			}
		}
	}
	log.Printf("Greeting: %s", r.GetMessage())
}

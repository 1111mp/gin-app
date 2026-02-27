package integration_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	protov1 "github.com/1111mp/gin-app/docs/proto/v1"
	"github.com/1111mp/gin-app/ent"
	"github.com/1111mp/gin-app/internal/dto"
	rmqClient "github.com/1111mp/gin-app/pkg/rabbitmq/rmq_rpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// Base settings
	// docker-compose service name
	host = "app"
	// test local server
	// host     = "127.0.0.1"
	attempts = 20

	// Attempts connection
	httpURL = "http://" + host + ":8080"

	requestTimeout = 5 * time.Second

	// HTTP REST
	basePathV1 = httpURL + "/api/v1"
	healthPath = basePathV1 + "/healthz"

	// gRPC
	grpcURL = host + ":8081"

	// RPC configs
	rpcServerExchange = "rpc_server"
	rpcClientExchange = "rpc_client"
	requests          = 10

	// RabbitMQ RPC
	// docker-compose service name
	rmqURL = "amqp://guest:guest@rabbitmq:5672/"
	// test local server
	// rmqURL = "amqp://guest:guest@127.0.0.1:5672/"

	// RabbitMQ RPC
	natsURL = "nats://guest:guest@nats:4222/"
)

var errHealthCheck = fmt.Errorf("url %s is not available", healthPath)

func doWebRequestWithTimeout(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

func getHealthCheck(url string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func healthCheck(attempts int) error {
	for attempts > 0 {
		statusCode, err := getHealthCheck(healthPath)
		if err != nil {
			return err
		}

		if statusCode == http.StatusOK {
			return nil
		}

		log.Printf("Integration tests: url %s is not available, attempts left: %d", healthPath, attempts)

		time.Sleep(time.Second)

		attempts--
	}

	return errHealthCheck
}

func TestMain(m *testing.M) {
	err := healthCheck(attempts)
	if err != nil {
		log.Fatalf("Integration tests: httpURL %s is not available: %s", httpURL, err)
	}

	log.Printf("Integration tests: httpURL %s is available", httpURL)

	code := m.Run()
	os.Exit(code)
}

// HTTP GET: /api/v1/user/:id
func TestHTTPGetUserV1(t *testing.T) {
	url := basePathV1 + "/users/1"
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	resp, err := doWebRequestWithTimeout(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// gRPC Client V1: GetHistory.
func TestClientGRPCV1(t *testing.T) {
	grpcConn, err := grpc.NewClient(grpcURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal("gRPC Client - init error - grpc.NewClient", err)
	}

	defer func() {
		err = grpcConn.Close()
		if err != nil {
			t.Fatal("gRPC Client - shutdown error - grpcClientV1.GetPostById", err)
		}
	}()

	grpcClientV1 := protov1.NewPostClient(grpcConn)

	for range requests {
		req := protov1.GetPostByIdRequest{Id: 1}
		post, err := grpcClientV1.GetPostById(t.Context(), &req)
		if err != nil {
			t.Fatal("gRPC Client - remote call error - grpcClientV1.GetPostById", err)
		}

		if post == nil {
			t.Fatalf("Post ID not found: expected %d", req.Id)
		}

		if post.Id != req.Id {
			t.Fatalf("Post ID mismatch: expected %d, got %d", req.Id, post.Id)
		}
	}
}

// RabbitMQ RPC Client V1: GetPostById.
func TestClientRMQRPCV1(t *testing.T) { //nolint: dupl,gocritic,nolintlint
	client, err := rmqClient.New(rmqURL, rpcServerExchange, rpcClientExchange)
	if err != nil {
		t.Fatal("RabbitMQ RPC Client - init error - rmqClient.New", err)
	}

	defer func() {
		err = client.Shutdown()
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - shutdown error - client.RemoteCall", err)
		}
	}()

	for range requests {
		req := dto.GetPostByIdRequest{ID: 1}
		var post ent.PostEntity

		err = client.RemoteCall("v1.get_post_by_id", req, &post)
		if err != nil {
			t.Fatal("RabbitMQ RPC Client - remote call error - client.RemoteCall", err)
		}

		if post.ID != req.ID {
			t.Fatalf("Post ID mismatch: expected %d, got %d", req.ID, post.ID)
		}
	}
}

package comunication

import (
	"CamMonitoring/src/utility"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"CamMonitoring/src/_interface"
	"CamMonitoring/src/message"
	"google.golang.org/grpc"
)

type CameraServiceServer struct {
	message.UnimplementedSubscriptionServiceServer
	CamIdCallback _interface.CamIdCallback
	server        *grpc.Server
	listener      net.Listener
	wg            sync.WaitGroup
}

// NewSubscriptionServiceServer initializes a new CameraServiceServer with the provided CamIdCallback.
func NewSubscriptionServiceServer(CamIdCallback _interface.CamIdCallback) *CameraServiceServer {
	return &CameraServiceServer{
		CamIdCallback: CamIdCallback,
	}
}

// GetCamIds implements the gRPC service method for retrieving CamIds.
func (s *CameraServiceServer) GetCamIds(_ context.Context, request *message.UserIdRequest) (
	*message.CamIdsResponse, error,
) {
	if s.CamIdCallback == nil {
		utility.ErrorLog().Println("GetCamIdsCallback not set")
		return nil, fmt.Errorf("GetCamIdsCallback not set")
	}

	// Your logic to retrieve cam_ids associated with user_id using the callback function
	userID := request.UserId
	camIds, err := s.CamIdCallback.GetAllCamIdsByUser(userID)
	if err != nil {
		utility.ErrorLog().Printf("Error retrieving CamIds: %v", err)
		return nil, err
	}

	return &message.CamIdsResponse{CamIds: camIds}, nil
}

// Start starts the gRPC server.
func (s *CameraServiceServer) Start(address string) {
	s.wg.Add(1)
	defer s.wg.Done()

	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		utility.ErrorLog().Fatalf("Failed to listen: %v", err)
	}

	s.server = grpc.NewServer()
	message.RegisterSubscriptionServiceServer(s.server, s)

	utility.InfoLog().Printf("Server is running on %s\n", address)

	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			utility.ErrorLog().Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT or SIGTERM
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	utility.InfoLog().Println("Shutting down gRPC server...")
	s.server.GracefulStop()
	utility.InfoLog().Println("gRPC server stopped")
}

// Stop stops the gRPC server.
func (s *CameraServiceServer) Stop() {
	if s.server != nil {
		s.server.Stop()
	}
	if s.listener != nil {
		err := s.listener.Close()
		utility.ErrorLog().Printf("Error closing listener: %v", err)
	}
	s.wg.Wait()
}

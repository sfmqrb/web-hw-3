package cache_client

import (
	"context"
	"google.golang.org/grpc"
	pb "hw3/BackEnd/cacheproto"
	"log"
	"time"
)

var conn grpc.ClientConnInterface
var err error

var C pb.CacheManagementClient
var cancel context.CancelFunc

type CacheNoteResponse struct {
	Note      string
	NoteId    string
	Exist     bool
	Access    bool
	MissCache bool
}
type CacheNoteRequest struct {
	RequestType int
	NoteId      string
	AuthorId    string
	Note        string
}
type CacheLoginRequest struct {
	RequestType int
	User        string
	Pass        string
}
type CacheLoginResponse struct {
	UserId    string
	WrongPass bool
	Exist     bool
	MissCache bool
}

const (
	address = "localhost:50051"
)

func Connect() {
	conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	C = pb.NewCacheManagementClient(conn)
	//todo
}
func RequestNoteCache(requestType int, note string, noteTitle string, noteType string, noteId string, authorId string) (*pb.CacheNoteResponse, error) {
	loginReq := &pb.CacheNoteRequest{
		RequestType: int32(requestType),
		NoteId:      noteId,
		AuthorId:    authorId,
		Note:        note,
		NoteTitle:   noteTitle,
		Type:        noteType,
	}
	Ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	cacheNoteResponse, err := C.CacheNoteRPC(Ctx, loginReq)
	if err != nil {
		print(err)
		return nil, err
	}
	if err != nil {
		print(err)
		return nil, err
	}
	recv, err := cacheNoteResponse.Recv()
	if err != nil {
		return nil, err
	}
	return recv, nil
}
func RequestLoginCache(requestType int, userName string, name string, pass string) (*pb.CacheLoginResponse, error) {
	Ctx, _ := context.WithTimeout(context.Background(), 10*time.Minute)
	loginReq := &pb.CacheLoginRequest{
		RequestType: int32(requestType),
		User:        userName,
		Pass:        pass,
		Name:        name,
	}
	cacheLoginResponse, err := C.CacheLoginRPC(Ctx, loginReq)
	if err != nil {
		return nil, err
	}
	recv, err := cacheLoginResponse.Recv()
	if err != nil {
		return nil, err
	}
	return recv, nil
}

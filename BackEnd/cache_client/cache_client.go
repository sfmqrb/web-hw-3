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
var Ctx context.Context
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
	Ctx, _ = context.WithTimeout(context.Background(), time.Minute)
}
func RequestNoteCache(requestType int, note string, noteTitle string, noteId string, authorId string) *pb.CacheNoteResponse {
	cacheNoteResponse, err := C.CacheNoteRPC(Ctx, &pb.CacheNoteRequest{
		RequestType: int32(requestType),
		NoteId:      noteId,
		AuthorId:    authorId,
		Note:        note,
		NoteTitle:   noteTitle,
	})
	if err != nil {
		return nil
	}
	return cacheNoteResponse
}
func RequestLoginCache(requestType int, userName string, name string, pass string) *pb.CacheLoginResponse {
	Ctx, _ = context.WithTimeout(context.Background(), time.Minute)
	loginReq := &pb.CacheLoginRequest{
		RequestType: int32(requestType),
		User:        userName,
		Pass:        pass,
		Name:        name,
	}
	cacheLoginResponse, err := C.CacheLoginRPC(Ctx, loginReq)
	if err != nil {
		return nil
	}
	recv, err := cacheLoginResponse.Recv()
	if err != nil {
		return nil
	}
	return recv
}

package cache_client

import (
	"context"
	"fmt"
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

func main() {
	fmt.Println("in client main")
	Connect()
	SendNoteRequestCache(CacheNoteRequest{
		RequestType: 0,
		NoteId:      "n id",
		AuthorId:    "a id",
		Note:        "noooooooooooote",
	})
}
func Connect() {
	conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	C = pb.NewCacheManagementClient(conn)
	Ctx, _ = context.WithTimeout(context.Background(), time.Second)
}
func SendNoteRequestCache(cReq CacheNoteRequest) CacheNoteResponse {
	res, err := C.CacheNoteRPC(Ctx, &pb.CacheNoteRequest{RequestType: int32(cReq.RequestType),
		NoteId:   cReq.NoteId,
		AuthorId: cReq.AuthorId,
		Note:     cReq.Note})
	if err != nil {
		log.Fatalf(err.Error())
	}
	cRes := CacheNoteResponse{
		Note:      res.Note,
		NoteId:    res.NoteId,
		Exist:     res.Exist,
		Access:    res.Access,
		MissCache: res.MissCache,
	}
	log.Printf("%t %t %v %t ", cRes.Exist, cRes.Access, cRes.Note, cRes.MissCache)
	return cRes
}
func SendLoginRequestCache(cReq CacheLoginRequest) CacheLoginResponse {
	res, err := C.CacheLoginRPC(Ctx, &pb.CacheLoginRequest{RequestType: int32(cReq.RequestType),
		User: cReq.User,
		Pass: cReq.Pass,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}
	cRes := CacheLoginResponse{
		UserId:    res.UserId,
		WrongPass: res.WrongPass,
		Exist:     res.Exist,
		MissCache: res.MissCache,
	}
	return cRes
}

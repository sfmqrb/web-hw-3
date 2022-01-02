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

var c pb.CacheManagementClient
var ctx context.Context
var cancel context.CancelFunc

type CacheResponse struct {
	Note      string
	NoteId    string
	Exist     bool
	Access    bool
	MissCache bool
}
type CacheRequest struct {
	RequestType int
	NoteId      string
	AuthorId    string
	Note        string
}

const (
	address = "localhost:50051"
)

func main() {
	fmt.Println("in client main")
	Connect()
	SendRequestCache(CacheRequest{
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
	c = pb.NewCacheManagementClient(conn)
	ctx, _ = context.WithTimeout(context.Background(), time.Second)
}
func SendRequestCache(cReq CacheRequest) CacheResponse {

	res, err := c.CacheNoteRPC(ctx, &pb.CacheRequest{RequestType: int32(cReq.RequestType),
		NoteId:   cReq.NoteId,
		AuthorId: cReq.AuthorId,
		Note:     cReq.Note})
	if err != nil {
		log.Fatalf(err.Error())
	}
	cRes := CacheResponse{
		Note:      res.Note,
		NoteId:    res.NoteId,
		Exist:     res.Exist,
		Access:    res.Access,
		MissCache: res.MissCache,
	}
	log.Printf("%t %t %v %t ", cRes.Exist, cRes.Access, cRes.Note, cRes.MissCache)
	return cRes
}

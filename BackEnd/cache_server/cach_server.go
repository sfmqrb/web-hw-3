package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"google.golang.org/grpc"
	pb "hw3/BackEnd/cacheproto"
	"log"
	"net"
)

type user struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	userId        int    `bun:"user_id,pk,autoincrement"`
	userName      string `bun:"user_name,notnull"`
	password      string `bun:"password,notnull"`
}

type note struct {
	bun.BaseModel `bun:"table:notes,alias:u"`
	noteId        int    `bun:"note_id,pk,autoincrement"`
	note          string `bun:"note,notnull"`
	authorId      int    `bun:"author_id,fk"`
}

func (u user) String() string {
	return fmt.Sprintf("User<%d %s %s>", u.userId, u.userName, u.password)
}

const (
	port = ":50051"
)

type CacheManagementServer struct {
	pb.UnimplementedCacheManagementServer
}

func (s *CacheManagementServer) CacheNoteRPC(ctx context.Context, in *pb.CacheRequest) (*pb.CacheResponse, error) {
	//log.Printf("Received: %v", in.GetName())
	log.Printf("Recived Cache Request: %v , %v , %d , %v ", in.Note, in.NoteId, in.RequestType, in.AuthorId)
	//var user_id int32 = int32(rand.Intn(100))
	//todo handle request
	return &pb.CacheResponse{
		Note:      "amir",
		NoteId:    "6556AC5",
		Exist:     false,
		Access:    true,
		MissCache: false,
	}, nil
}
func connectToDB() {
	ctx := context.Background()
	// Open a PostgreSQL database.
	dsn := "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	// Create a Bun db on top of it.
	db := bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	count, _ := db.NewSelect().Model((*user)(nil)).Count(ctx)
	fmt.Println(count)
	// Select a random number.
	exists, _ := db.NewSelect().Model((*note)(nil)).Where("author_id = 2").Exists(ctx)
	fmt.Println(exists)

}
func main() {
	connectToDB()
	//startGrpcServer()
}

func startGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCacheManagementServer(s, &CacheManagementServer{})
	//pb.RegisterRequestCacheServer(s, &RequestCacheServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

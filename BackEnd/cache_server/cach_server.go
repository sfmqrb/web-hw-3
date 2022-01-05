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
	"strconv"
)

const (
	// action types of cache note request
	save = 1
	del  = 2
	get  = 3
	// action types of cache login request
	signIn = 1
	signUp = 2
)

type user struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	userId        int    `bun:"user_id,pk,autoincrement"`
	userName      string `bun:"user_name,notnull"`
	password      string `bun:"password,notnull"`
}

type note struct {
	bun.BaseModel `bun:"table:notes"`
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

var db *bun.DB

func (s *CacheManagementServer) CacheLoginRPC(ctx context.Context, in *pb.CacheLoginRequest) (*pb.CacheLoginResponse, error) {
	//todo handle request
	return &pb.CacheLoginResponse{
		UserId:    "3",
		WrongPass: false,
		Exist:     false,
		MissCache: false,
	}, nil
}
func (s *CacheManagementServer) CacheNoteRPC(ctx context.Context, in *pb.CacheNoteRequest) (*pb.CacheNoteResponse, error) {
	//log.Printf("Received: %v", in.GetName())
	log.Printf("Recived Cache Request: %v , %v , %d , %v ", in.Note, in.NoteId, in.RequestType, in.AuthorId)
	var res *pb.CacheNoteResponse
	switch in.RequestType {
	case save:
		aId, _ := strconv.Atoi(in.AuthorId)
		noteObj := &note{
			BaseModel: bun.BaseModel{},
			note:      in.Note,
			authorId:  aId,
		}
		exec, err := db.NewInsert().Model(noteObj).Exec(ctx)
		if err != nil {
			id, err := exec.LastInsertId()
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			res.NoteId = strconv.FormatInt(id, 10)
		}
	case del:
		nId, _ := strconv.Atoi(in.NoteId)
		aId, _ := strconv.Atoi(in.AuthorId)
		noteObj := &note{noteId: nId, authorId: aId}
		_, err := db.NewDelete().Model(noteObj).WherePK().Exec(ctx)
		if err != nil {
			res.Exist = true
			res.Access = true
		} else {
			res.Access = false
		}
	case get:
		noteObj := new(note)
		err := db.NewSelect().Model(noteObj).Where("note_id = ?", in.NoteId).Scan(ctx)
		if err != nil {
			fmt.Println(err)
			res.Exist = false
		} else {
			res.Note = noteObj.note
			res.NoteId = string(rune(noteObj.noteId))
			res.Exist = true
			res.Access = in.AuthorId == string(rune(noteObj.authorId))
			//todo missCache
		}
	}
	//var user_id int32 = int32(rand.Intn(100))
	//todo handle request
	return &pb.CacheNoteResponse{
		Note:      "amir",
		NoteId:    "6556AC5",
		Exist:     false,
		Access:    true,
		MissCache: false,
	}, nil
}
func connectToDB() {
	// Open a PostgreSQL database.
	dsn := "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	// Create a Bun db on top of it.
	db = bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	userObj := new(user)
	ctx := context.Background()
	err := db.NewSelect().Model(userObj).Where("user_id = ?", 1).Scan(ctx)
	userObj = &user{
		userId:   10,
		userName: "soo",
		password: "soo1234",
	}
	res, err := db.NewInsert().Model(userObj).Exec(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

}
func main() {
	connectToDB()
	startGrpcServer()

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

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
	Login  = 1
	signUp = 2
)

type user struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	UserId        int    `bun:"user_id,pk,autoincrement"`
	UserName      string `bun:"user_name,notnull"`
	Password      string `bun:"password,notnull"`
}

type note struct {
	bun.BaseModel `bun:"table:notes,alias:u"`
	NoteId        int    `bun:"note_id,pk,autoincrement"`
	Note          string `bun:"note,notnull"`
	AuthorId      int    `bun:"author_id"`
}

func (u user) String() string {
	return fmt.Sprintf("User<%d %s %s>", u.UserId, u.UserName, u.Password)
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
	var res *pb.CacheLoginResponse
	res = &pb.CacheLoginResponse{}
	switch in.RequestType {
	case Login:
		userObj := &user{}
		err := db.NewSelect().Model(userObj).Where("user_name = ? AND password = ?", in.User, in.Pass).Scan(ctx)
		if err != nil {
			fmt.Println(err)
			res.WrongPass = true
		} else {
			res.UserId = strconv.Itoa(userObj.UserId)
			res.Exist = true
		}
	case signUp:
		userObj := &user{}
		err := db.NewSelect().Model(userObj).Where("user_name = ?", in.User).Scan(ctx)
		if err == nil {
			res.Exist = true
		} else if err == sql.ErrNoRows {
			res.Exist = false
			userObj = &user{
				BaseModel: bun.BaseModel{},
				UserName:  in.User,
				Password:  in.Pass,
			}
			exec, err := db.NewInsert().Model(userObj).Exec(ctx)
			if err != nil {
				id, err := exec.LastInsertId()
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				res.UserId = strconv.FormatInt(id, 10)
			}
		}
	}
	return res, nil
}
func (s *CacheManagementServer) CacheNoteRPC(ctx context.Context, in *pb.CacheNoteRequest) (*pb.CacheNoteResponse, error) {
	//log.Printf("Received: %v", in.GetName())
	log.Printf("Recived Cache Request: %v , %v , %d , %v ", in.Note, in.NoteId, in.RequestType, in.AuthorId)
	var res *pb.CacheNoteResponse
	res = &pb.CacheNoteResponse{
		Note:      "",
		NoteId:    "",
		Exist:     true,
		Access:    false,
		MissCache: false,
	}
	switch in.RequestType {
	case save:
		aId, _ := strconv.Atoi(in.AuthorId)
		noteObj := &note{
			BaseModel: bun.BaseModel{},
			Note:      in.Note,
			AuthorId:  aId,
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
		noteObj := &note{NoteId: nId, AuthorId: aId}
		_, err := db.NewDelete().Model(noteObj).Where("note_id = ? AND author_id = ?", nId, aId).Exec(ctx)
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
			res = &pb.CacheNoteResponse{
				Note:      noteObj.Note,
				NoteId:    strconv.Itoa(noteObj.NoteId),
				Exist:     true,
				Access:    in.AuthorId == strconv.Itoa(noteObj.AuthorId),
				MissCache: false,
			}
			//todo missCache
		}
	}
	//var user_id int32 = int32(rand.Intn(100))
	//todo handle request
	return res, nil
}
func connectToDB() {
	// Open a PostgreSQL database.
	dsn := "postgres://postgres:admin@localhost:5432/postgres?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	// Create a Bun db on top of it.
	db = bun.NewDB(pgdb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	//userObj := &user{}
	//err := db.NewSelect().Model(userObj).Where("user_name = ? AND password = ?", "amir", "Xamm2666").Scan(context.Context())
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

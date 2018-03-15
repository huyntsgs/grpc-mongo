package main

import (
	"flag"
	"log"
	"net"

	"github.com/golang/glog"
	pb "github.com/huyntsgs/grpc-mongo/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type productServer struct {
	session *mgo.Session
}

func (s productServer) GetProducts(c context.Context, req *pb.ProductReq) (*pb.ProductRes, error) {
	ss := s.session.Copy()
	defer ss.Close()

	items := ss.DB("waterfall").C("newitems")

	var allItems []*pb.Item
	log.Println("limit getProducts", req.Offset)

	err := items.Find(bson.M{"IdCat": req.IdCat, "Status": bson.M{"$lt": 4}}).Sort("Status").Limit(int(req.Limit)).Skip(int(req.Offset)).All(&allItems)

	if err != nil {
		log.Println("err getProducts:", err)
		return &pb.ProductRes{}, err
	}
	//var ar []*pb.Item
	//item := pb.Item{ID: []byte("123456789abc"), CusName: "user1", IdCat: 1}
	//ar = append(ar, &item)
	return &pb.ProductRes{Items: allItems}, nil
}
func (s productServer) Login(c context.Context, req *pb.User) (*pb.LoginRes, error) {
	ss := s.session.Copy()
	defer ss.Close()

	users := ss.DB("waterfall").C("user")
	user := pb.User{}

	loginRes := &pb.LoginRes{ErrCode: 0}

	log.Println("user info ", req)

	err := users.Find(bson.M{"name": req.Name, "pass": req.Pass}).One(&user)

	if err != nil {
		log.Println("err:", err)
		loginRes.ErrCode = 1
	} else {
		loginRes.UserInfo = &user
	}

	return loginRes, nil
}

func Run() error {
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}
	server := grpc.NewServer()

	//create mongo session
	session := CreateSession()
	defer session.Close()

	pb.RegisterProductServiceServer(server, productServer{session})
	pb.RegisterUserServiceServer(server, productServer{session})

	server.Serve(listen)
	return nil
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := Run(); err != nil {
		glog.Fatal(err)
	}
}

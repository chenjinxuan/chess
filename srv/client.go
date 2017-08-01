package srv
import (
    "log"
   // "os"

    "chess/srv/srv-auth/vendor/golang.org/x/net/context"
    "chess/srv/srv-auth/vendor/google.golang.org/grpc"
    pb "chess/srv/srv-auth/proto"
)

const (
    address     = "localhost:50001"
)

func main() {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
	log.Fatal("did not connect: %v", err)
    }
    defer conn.Close()
    c := pb.NewAuthServiceClient(conn)

//user:=2
//    if len(os.Args) >1 {
//	user = os.Args[1]
//    }
    r, err := c.Auth(context.Background(), &pb.AuthArgs{UserId: 10000001,Token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MDE3MzgyOTUsInVzZXJpZCI6IjEwMDAwMDAxIn0.rJMc5nUJTsbTEolnbt0y7dFK2ovFuYCU_MaGu13V6AY"})
    if err != nil {
	log.Fatal("could not greet: %v", err)
    }
    log.Printf("Greeting: %s", r.Ret)
}
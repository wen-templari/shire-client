package core_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"strings"
	"testing"

	"github.com/templari/shire-client/core"
	"github.com/templari/shire-client/model"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

type ServiceA struct {
	id int
}

func (s *ServiceA) Add(args *Args, reply *int) error {
	*reply = s.id
	return nil
}

func TestRPC(t *testing.T) {
	makeServer := func(id int) {
		arith := new(Arith)
		server := rpc.NewServer()
		server.HandleHTTP(fmt.Sprint("/", id), fmt.Sprint("/debug/", id))
		server.Register(arith)
		serviceA := &ServiceA{
			id: id,
		}
		server.Register(serviceA)
	}
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	makeServer(1)
	makeServer(2)

	callServer := func(id int) {
		client, err := rpc.DialHTTPPath("tcp", "127.0.0.1:1234", fmt.Sprint("/", id))
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		args := &Args{7, 8}
		var reply int
		err = client.Call("Arith.Multiply", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
		err = client.Call("ServiceA.Add", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)
	}

	callServer(1)
	callServer(2)
}

// func TestRPCConnect(t *testing.T) {

// 	rpc.HandleHTTP()
// 	l, e := net.Listen("tcp", ":1234")
// 	if e != nil {
// 		log.Fatal("listen error:", e)
// 	}
// 	go http.Serve(l, nil)

// 	_, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
// 	if err != nil {
// 		log.Fatal("dialing:", err)
// 	}
// }

func TestMultiHandlers(t *testing.T) {
	l, _ := net.Listen("tcp", ":1234")
	c := core.MakeCore("http://localhost:3011")
	go core.StartHttpServer(c, l)
	t.Log("Starting RPC server")
	rpcServer := rpc.NewServer()
	rpcServer.HandleHTTP("/rpc", "/debug")
	rpcServer.Register(new(Arith))
	rpcServer.Register(new(ServiceA))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		t.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	http.ListenAndServe(":1234", nil)
	// for {

	// }

	// client, err := rpc.DialHTTP("tcp", "
}

type StringifyArgs struct {
	Id   int
	Name string
}

func TestStringify(t *testing.T) {
	s := "{{4b34119e-3e0b-40ae-ba2c-e707dafd2d0b 2} APPEND 61 {'from':206,'to':0,'content':'hello Again From Alice','groupId':61,'time':'123'}}"
	s = s[1 : len(s)-1]
	arr := strings.Split(s, " ")
	log.Println(arr)
	log.Println(arr[len(arr)-1])
	msg := model.Message{}
	l := arr[len(arr)-1]
	json.Unmarshal([]byte(l), &msg)
	log.Print(msg)
}

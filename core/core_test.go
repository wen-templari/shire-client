package core_test

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"testing"

	"github.com/templari/shire-client/core"
	"github.com/templari/shire-client/core/model"
)

func TestCoreLogin(t *testing.T) {
	core := core.MakeCore("http://localhost:3011", nil)
	_, err := core.Login(1, "12346")
	if err != nil {
		t.Error(err)
	}
}

func TestSendMessage(t *testing.T) {
	bob := core.MakeCore("http://localhost:3011", nil)

	if _, err := bob.Register("bob", "12345"); err != nil {
		t.Error(err)
	}

	messageChan := make(chan model.Message)
	bob.Subscribe(messageChan)

	go func() {
		for {
			message := <-messageChan
			log.Printf("received message: %v", message)
		}
	}()

	sender := core.MakeCore("http://localhost:3011", nil)
	if _, err := bob.Register("tom", "12345"); err != nil {
		t.Error(err)
	}

	receiver, err := sender.GetUserById(bob.GetUser().Id)
	if err != nil {
		t.Error(err)
	}

	err = sender.SendMessage(model.Message{
		To:      receiver.Id,
		From:    sender.GetUser().Id,
		Content: "hello",
		Time:    "123",
	})
	if err != nil {
		t.Errorf("failed to send message: %v", err)
	}

}

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

type ServiceA struct{}

func (s *ServiceA) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func TestRPC(t *testing.T) {
	arith := new(Arith)
	rpc.Register(arith)
	serviceA := new(ServiceA)
	rpc.Register(serviceA)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
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

package main

import (
	"bufio"
	"fmt"
	"net"
	//"time"
)

type Model string

var (
	LampDimmable    = Model("lm12")
	LampNotDimmable = Model("lm15")
	Appliance       = Model("am12")
	Tranceiver      = Model("tm12")
)

/*
lm12: lamp module, dimmable
lm15: lamp module, not dimmable
am12: appliance module, not dimmable
tm12: tranceiver module, not dimmable. Its unitcode is 1.
*/

type X10 struct {
	Name      string
	Model     Model
	HouseCode string
	UnitCode  int
}

func (this *X10) define() string {
	return fmt.Sprintf("define %s X10 %s %s %v", this.Name, this.Model, this.HouseCode, this.UnitCode)
}

func (this *X10) delete() string {
	return fmt.Sprintf("delete %s", this.Name)
}

func (this *X10) Command(cm command) {
	c := fmt.Sprintf("%s\nset %s %s\n%s", this.define(), this.Name, cm, this.delete())
	cmd(c)
}

func main() {
	x10 := &X10{
		Name:      "golamp",
		Model:     LampDimmable,
		HouseCode: "A",
		UnitCode:  2,
	}

	//fmt.Println(x10.define())
	x10.Command(On)
	x10.Command(Off)

	//cmd("define golamp15 X10 lm12 A 15\nset golamp15 on")
	//time.Sleep(time.Second * 5)
	//cmd("define golamp16 X10 lm12 A 16\nset golamp16 on")

}

type command string

var (
	DimDown    = command("dimdown")
	DimUp      = command("dimup")
	Off        = command("off")
	On         = command("on")
	OnTill     = command("on-till")
	OnForTimer = command("on-for-timer")
)

/*dimdown           # requires argument, see the note
  dimup             # requires argument, see the note
  off
  on
  on-till           # Special, see the note
  on-for-timer      # Special, see the note
*/

func cmd(s string) {
	fmt.Println("connect")
	fmt.Println(s)
	conn, err := net.Dial("tcp", "localhost:7072")
	if err != nil {
		fmt.Println(err)
	}

	send("", conn)
	send("\n"+s+"\nquit", conn)
	err = conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("disconnect")

}

func send(cmd string, conn net.Conn) {
	fmt.Fprintf(conn, cmd+"\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("> " + status)
}

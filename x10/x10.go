package x10

import (
	"bufio"
	"fmt"
	"net"
	"runtime"
	//"time"
)

type model string

var (
	LampDimmable    = model("lm12")
	LampNotDimmable = model("lm15")
	Appliance       = model("am12")
	Tranceiver      = model("tm12")
)

var (
	Models = map[string]model{
		"lm12": LampDimmable,
		"lm15": LampNotDimmable,
		"am12": Appliance,
		"tm12": Tranceiver,
	}
)

/*
lm12: lamp module, dimmable
lm15: lamp module, not dimmable
am12: appliance module, not dimmable
tm12: tranceiver module, not dimmable. Its unitcode is 1.
*/

func InitCM11() {
	device := "/dev/ttyUSB0"
	if runtime.GOOS == "windows" {
		device = "com6"
	}
	cmd(fmt.Sprintf("define x10if CM11 %s", device))
}

type X10 struct {
	Name      string
	Model     model
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
	//c := fmt.Sprintf("%s\nset %s %s\n%s", this.define(), this.Name, cm, "")

	//c := fmt.Sprintf("%s\nset %s %s\n%s", "", this.Name, cm, "")
	//c := fmt.Sprintf("%s\nset %s %s\n%s", this.define(), this.Name, cm, this.delete())
	//c := fmt.Sprintf("set %s %s", this.Name, cm)
	cmd(c)
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

var FhemServer = "localhost:7072"
var Debug = true

func cmd(s string) {
	if Debug {
		fmt.Println("connect")
		fmt.Println(s)
	}
	conn, err := net.Dial("tcp", FhemServer)
	if err != nil {
		panic(err.Error())
	}

	send("", conn)
	send("\n"+s+"\nquit", conn)
	err = conn.Close()
	if err != nil {
		panic(err.Error())
	}
	if Debug {
		fmt.Println("disconnect")
	}
}

func send(cmd string, conn net.Conn) {
	fmt.Fprintf(conn, cmd+"\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err.Error())
	}
	if Debug {
		fmt.Println("> " + status)
	}
}

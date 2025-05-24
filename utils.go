package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ui funcs

func takeOrder() Order {
	var orderName string
	var items int
	var current_location string
	var destination string
	var ship string 

	// take input
	fmt.Print("Enter order name > ")
	fmt.Scanf("%s", &orderName)

	fmt.Print("Enter items it'll store > ")
	fmt.Scanf("%v", &items)

	fmt.Print("Enter current location > ")
	fmt.Scanf("%v", &current_location)

	fmt.Print("Enter destination > ")
	fmt.Scanf("%v", &destination)

	fmt.Print("Enter the ship i'll go in > ")
	fmt.Scanf("%v", &ship)

	order := CreateOrder(orderName, items, current_location, destination, ship)

	return order
}

func DisplayPortsStatus(ports map[string]*Port) {
    fmt.Println("\n📦 Port Occupancy Status:")
    fmt.Println("----------------------------")
    for name, port := range ports {
        percentUsed := float64(port.shipsResidingInPort) / float64(port.spaceInPort) * 100
        bar := strings.Repeat("█", int(percentUsed/5)) + strings.Repeat("░", 20-int(percentUsed/5))
		title := cases.Title(language.English)
		fmt.Printf("Port %-8s [%s] %2.0f%% (%d/%d)\n", title.String(name), bar, percentUsed, port.shipsResidingInPort, port.spaceInPort)
    }
    fmt.Println()
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cls")
	default:
		runCmd("clear")
	}
}

func loadingBar(iteration, total int, prefix, suffix string, length int, fill string) {
	percent := float64(iteration) / float64(total)
	filledLength := int(length * iteration / total)

	end := ">"

	if iteration == total {
		end = "="
	}

	bar := strings.Repeat(fill, filledLength) + end + strings.Repeat("-", (length - filledLength))
	fmt.Printf("\r%s [%s] %f%% %s", prefix, bar, percent, suffix)

	if iteration == total {
		fmt.Println()
	}
}

func MainCycle() {
	steps := make([]int, 30)	
	for i := range steps {
		loadingBar(i+1, 30, "Order", "Delivered", 25, "=")
		time.Sleep(100 * time.Millisecond)
	}
}
// logic funcs

func CreateOrder(name string, items int, current_location string, destination string, ship string) Order {
	return Order {
		storage: items,
		current_location: current_location,
		destination: destination,
		holder: ship,
		name: name,
	}
}

//initialise ports
func initialisePorts() map[string]*Port {
    ports := map[string]*Port{
        "alpha":   {name: "alpha", spaceInPort: 5},
        "beta":    {name: "beta", spaceInPort: 10},
        "gamma":   {name: "gamma", spaceInPort: 15},
        "delta":   {name: "delta", spaceInPort: 20},
        "epsilon": {name: "epsilon", spaceInPort: 25},
        "eta":     {name: "eta", spaceInPort: 30},
        "theta":   {name: "theta", spaceInPort: 35},
    }

    return ports
}

// docking

func (p *Port) HasSpace() bool {
	return p.shipsResidingInPort < p.spaceInPort
}

func (p *Port) DockShip(ship Ship) bool {
	if p.HasSpace() {
		p.shipsResidingInPort++
		fmt.Printf("Ship %s docked at port %s.\n", ship.name, p.name)
		return true
	} else {
		fmt.Printf("Port %s is full. Ship %s cannot dock.\n", p.name, ship.name)
	}

	return false
}

func (p *Port) ReleaseShip(ship Ship) {
	if p.shipsResidingInPort > 0 {
		p.shipsResidingInPort--
		fmt.Printf("Ship %s has left port %s.\n", ship.name, p.name)
	}
}

// misc

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

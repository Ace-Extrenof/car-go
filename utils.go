package main

import (
	"fmt"
	"math"
	"math/rand"
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
	percentUsed := 0.0
    for name, port := range ports {
        percentUsed = float64(port.shipsResidingInPort) / float64(port.spaceInPort) * 100

		blocks := int(math.Round(percentUsed / 5))
		if percentUsed > 0 && blocks == 0 {
			blocks = 1
		}
		if blocks > 20 {
			blocks = 20
		}

        bar := strings.Repeat("█", blocks) + strings.Repeat(" ", 20-blocks)

		title := cases.Title(language.English)
		fmt.Printf("Port %-8s [%s] %2.0f%% (%d/%d)\n", title.String(name), bar, percentUsed, port.shipsResidingInPort, port.spaceInPort)
    }
    fmt.Println()
}

func DisplayCargoSummary(cargo map[string]int) {
	fmt.Println("\n📊 Cargo Delivery Summary:")
	fmt.Println("-------------------------------")

	max := 0
	for _, amount := range cargo {
		if amount > max {
			max = amount
		}
	}

	for port, amount := range cargo {
		barLen := 0
		if max > 0 {
			barLen = int(float64(amount) / float64(max) * 30)
		}
		bar := strings.Repeat("█", barLen)
		title := cases.Title(language.English)
		fmt.Printf("%-10s | %-30s (%d units)\n", title.String(port), bar, amount)
	}
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
		time.Sleep(300 * time.Millisecond)
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

var samplePortNames = []string{"alpha", "beta", "gamma", "delta", "epsilon", "eta", "theta"}
var sampleShipNames = []string{"Titan", "Orion", "Voyager", "Enterprise", "Falcon", "Nebula", "Pegasus"}

func generateRandomOrder() Order {
	name := fmt.Sprintf("order-%d", rand.Intn(1000))
	storage := rand.Intn(1000)

	from := samplePortNames[rand.Intn(len(samplePortNames))]
	to := samplePortNames[rand.Intn(len(samplePortNames))]

	for from == to {
		to = samplePortNames[rand.Intn(len(samplePortNames))]
	}

	ship := sampleShipNames[rand.Intn(len(sampleShipNames))]

	return CreateOrder(name, storage, from, to, ship)
}

//initialise ports
// func initialisePorts() map[string]*Port {
//     ports := map[string]*Port{
//         "alpha":   {name: "alpha", spaceInPort: 5},
//         "beta":    {name: "beta", spaceInPort: 10},
//         "gamma":   {name: "gamma", spaceInPort: 15},
//         "delta":   {name: "delta", spaceInPort: 20},
//         "epsilon": {name: "epsilon", spaceInPort: 25},
//         "eta":     {name: "eta", spaceInPort: 30},
//         "theta":   {name: "theta", spaceInPort: 35},
//     }
//
//     return ports
// }

func initialisePorts() map[string]*Port {
    return map[string]*Port{
        "alpha":   {name: "alpha", spaceInPort: 5, shipsResidingInPort: 0},
        "beta":    {name: "beta", spaceInPort: 10, shipsResidingInPort: 0},
        "gamma":   {name: "gamma", spaceInPort: 15, shipsResidingInPort: 0},
        "delta":   {name: "delta", spaceInPort: 20, shipsResidingInPort: 0},
        "epsilon": {name: "epsilon", spaceInPort: 25, shipsResidingInPort: 0},
        "eta":     {name: "eta", spaceInPort: 30, shipsResidingInPort: 0},
        "theta":   {name: "theta", spaceInPort: 35, shipsResidingInPort: 0},
    }
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

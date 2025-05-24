package main

import (
	"fmt"
	"strings"
)

// import "fmt"

func main() {
	var numOrders int
	var auto string
	fmt.Print("Randomised orders? (y/n) > ")
	fmt.Scanf("%s", &auto)

	fmt.Print("How many orders would you like to process? > ")
	fmt.Scanf("%d", &numOrders)

	ports := initialisePorts()

	orders := make([]Order, 0)
	for i := 0; i < numOrders; i++ {
		fmt.Printf("\n📦 Enter details for Order #%d\n", i+1)
		if strings.ToLower(auto) == "y" {
			order := generateRandomOrder()
			fmt.Printf("Auto-created order: %+v\n", order)

			orders = append(orders, order)
		} else {
			order := takeOrder()
			orders = append(orders, order)
		}
	}
	
	for i, order := range orders {
		fmt.Printf("\n🚚 Processing Order #%d: %s\n", i+1, order.name)

		ship := Ship{name: order.holder, cargo: order.storage}
		currentPort, ok1 := ports[strings.ToLower(order.current_location)]
		destinationPort, ok2 := ports[strings.ToLower(order.destination)]

		if !ok1 || !ok2 {
			fmt.Printf("❌ Error: Invalid port in order '%s'. Skipping.\n", order.name)
			continue
		}

		DisplayPortsStatus(ports)

		if !currentPort.DockShip(ship) {
			fmt.Printf("⚠️ Ship %s could not start journey from %s.\n", ship.name, currentPort.name)
			continue
		}

		DisplayPortsStatus(ports)
		fmt.Println("🚀 En route...")
		MainCycle()

		currentPort.ReleaseShip(ship)
		DisplayPortsStatus(ports)

		if !destinationPort.DockShip(ship) {
			fmt.Printf("🚫 Ship %s reached %s but cannot dock (port full).\n", ship.name, destinationPort.name)
			continue
		}

		fmt.Printf("✅ Order '%s' delivered successfully to %s!\n", order.name, destinationPort.name)
		DisplayPortsStatus(ports)
	}

}

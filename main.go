package main

import (
    "fmt"
    "strings"
)

func main() {
    dashboardResult, err := ShowDashboard()
    if err != nil {
        fmt.Println("Err running dashboard:", err)
        return
    }
    if dashboardResult == nil {
        fmt.Println("Exited Dashboard.")
        return
    }

    auto := dashboardResult.Randomised
    numOrders := dashboardResult.NumOrders

    ClearTerminal()

    ports := initialisePorts() // Make sure this initializes your ports map

    orders := make([]Order, 0, numOrders)
    for i := 0; i < numOrders; i++ {
        fmt.Printf("\n📦 Enter details for Order #%d\n", i+1)
        if auto {
            order := generateRandomOrder()
            fmt.Printf("Auto-created order: %+v\n", order)
            orders = append(orders, order)
            ClearTerminal()
        } else {
            order := takeOrder()
            orders = append(orders, order)
        }
    }

    cargoDelivered := make(map[string]int)

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
        MainCycle() // Simulate travel / passage of time

        currentPort.ReleaseShip(ship)
        DisplayPortsStatus(ports)

        if !destinationPort.DockShip(ship) {
            fmt.Printf("🚫 Ship %s reached %s but cannot dock (port full).\n", ship.name, destinationPort.name)
            continue
        }

        fmt.Printf("✅ Order '%s' delivered successfully to %s!\n", order.name, destinationPort.name)
        cargoDelivered[destinationPort.name] += order.storage
        DisplayPortsStatus(ports)
    }

    DisplayCargoSummary(cargoDelivered)

    fmt.Println("Press enter to exit...")
    fmt.Scanln()
}


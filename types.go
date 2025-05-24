package main

type Order struct {
	storage int
	current_location string
	destination string
	holder string
	name string
}

type Port struct {
	name string
	spaceInPort int
	shipsResidingInPort int
}

type Ship struct {
	name string
	cargo int
}

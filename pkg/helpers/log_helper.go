package helpers

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

func LogStartup() {
	fmt.Print("\n\n")
	figure.NewFigure("Xanny", "doom", true).Print()
	fmt.Print("\n\n")
	fmt.Println("Xanny Go Template: A scalable and efficient Go boilerplate for modern web applications.")
	fmt.Println("Support me! ETH/POL:0x4418f0009606B8a5666B353ac5Fe49E874Fb8b61")
	fmt.Println("-------------------------------------------------------------------------------------------")
	fmt.Print("\n\n")
}

package types

import "fmt"

func PrintXLSX(sheetName string, header []string, body interface{}) {
	printLine()
	fmt.Println(sheetName)
	printLine()
	for _, v := range header {
		fmt.Print(v, " | ")
	}
	println()
	printLine()
	switch b := body.(type) {
	case []string:
		{
			for _, col := range b {
				fmt.Print(col, " | ")
			}
			println()
			printLine()
		}
	case [][]string:
		{
			for _, row := range b {
				for _, col := range row {
					fmt.Print(col, " | ")
				}
				println()
				printLine()
			}
		}
	}
	println()
}

func printLine() {
	fmt.Println("--------------------------------------------------------------------------------------------------------------------------------")
}

func println() {
	fmt.Println("")
}

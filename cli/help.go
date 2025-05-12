package main

import "fmt"

func showHelp() {
	fmt.Println(headerStyle.Render("Orca CLI"))

	fmt.Println(subHeaderStyle.Render("\nUsage:"))
	fmt.Println(infoStyle.Render("  command [options]"))

	fmt.Println(subHeaderStyle.Render("\nCommands:"))
	fmt.Println(infoStyle.Render("  start   ") + "- Start the Orca stack")
	fmt.Println(infoStyle.Render("  stop    ") + "- Stop the Orca stack")
	fmt.Println(infoStyle.Render("  status  ") + "- Show status of the Orca components")
	fmt.Println(infoStyle.Render("  destroy ") + "- Delete all Orca resources")
	fmt.Println(
		infoStyle.Render("  help    ") + "- Show this help message or help for a specific command",
	)

	fmt.Println(subHeaderStyle.Render("\nExamples:"))
	fmt.Println(infoStyle.Render("  orca start"))
	fmt.Println(infoStyle.Render("  orca stop"))
	fmt.Println(infoStyle.Render("  orca status"))
	fmt.Println(infoStyle.Render("  orca destroy"))
	fmt.Println(infoStyle.Render("  orca help start"))
}

func showCommandHelp(command string) {
	switch command {
	case "start":
		fmt.Println(subHeaderStyle.Render("'start' command - Start the Orca stack"))
		fmt.Println(infoStyle.Render("\nUsage:"))
		fmt.Println(infoStyle.Render("  orca start"))
		fmt.Println(infoStyle.Render("\nExample:"))
		fmt.Println(infoStyle.Render("  orca start"))

	case "stop":
		fmt.Println(subHeaderStyle.Render("'stop' command - Stop the Orca stack"))
		fmt.Println(infoStyle.Render("\nUsage:"))
		fmt.Println(infoStyle.Render("  orca stop"))
		fmt.Println(infoStyle.Render("\nExample:"))
		fmt.Println(infoStyle.Render("  orca stop"))

	case "status":
		fmt.Println(subHeaderStyle.Render("'status' command - Show status of the Orca services"))
		fmt.Println(infoStyle.Render("\nUsage:"))
		fmt.Println(infoStyle.Render("  orca status"))
		fmt.Println(infoStyle.Render("\nExample:"))
		fmt.Println(infoStyle.Render("  orca status"))

	case "destroy":
		fmt.Println(subHeaderStyle.Render("'destroy' command - Tear down the Orca environment"))
		fmt.Println(infoStyle.Render("\nUsage:"))
		fmt.Println(infoStyle.Render("  orca destroy"))
		fmt.Println(infoStyle.Render("\nExample:"))
		fmt.Println(infoStyle.Render("  orca destroy"))

	default:
		fmt.Println(errorStyle.Render(fmt.Sprintf("Unknown command: %s", command)))
		showHelp()
	}
}

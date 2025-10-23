package cmd

import (
	"fmt"
	"os"
	"todo/internal/items"
	"todo/internal/items/console"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/mbndr/figlet4go"
)

type VaultMenuState int

const (
	VaultMenuNone = iota
	VaultMenuInsert
	VaultMenuDelete
	VaultMenuEdit
	VaultMenuHistory
)

func RenderMenu(index int) {
	vault := items.LoadVaultFromFile()

	console.Clear()
	title()
	for itemIndex, item := range vault {
		if itemIndex == index {
			fmt.Printf(console.Blue+"➡ "+console.Reset+"%v\n", item.Name)
		} else {
			fmt.Println(item.Name)
		}
	}

	fmt.Println()

	// fmt.Println("E - edit || R - Remove || I - Insert || H - History || Q = Quit")
	fmt.Println(vaultMenuOptions([]string{"Insert", "Remove", "Edit", "History", "Quit"}))

	vaultState := VaultMenuNone

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.Down {
			index++
			vaultIndexCorrection(&index, vault)

			return true, nil
		} else if key.Code == keys.Up {
			index--
			vaultIndexCorrection(&index, vault)

			return true, nil
		} else if key.String() == "q" || key.Code == keys.CtrlC {
			console.Clear()
			os.Exit(0)
		} else if key.String() == "i" {
			vaultState = VaultMenuInsert
			return true, nil
		} else if key.String() == "r" {
			vaultState = VaultMenuDelete
			vault = append(vault[:index], vault[index+1:]...)
			items.SaveVaultToFile(vault)
			return true, nil
		} else if key.String() == "e" {
			vaultState = VaultMenuEdit
			return true, nil
		} else if key.String() == "h" {
			itemHistory(vault[index])
		}

		return false, nil
	})

	switch vaultState {
	case VaultMenuInsert:
		insertItemMenu()
	case VaultMenuEdit:
		editItemMenu(index)
	}

	RenderMenu(index)
}

// Corrects the current index of the vault menu, so it wouldn't overflow
func vaultIndexCorrection(index *int, vault []items.Item) {
	if *index >= len(vault) {
		*index = 0
	} else if *index < 0 {
		*index = len(vault) - 1
	}
}

func insertItemMenu() {
	console.Clear()
	var data string

	fmt.Print("Item Name: ")
	fmt.Scan(&data)

	vault := items.LoadVaultFromFile()
	vault = append(vault, items.NewItem(data))
	items.SaveVaultToFile(vault)

	RenderMenu(0)
}

func editItemMenu(index int) {
	vault := items.LoadVaultFromFile()

	console.Clear()
	fmt.Print("Editing item: ")

	originalItemName := vault[index].Name
	itemName := vault[index].Name
	fmt.Scan(&itemName)

	if originalItemName != itemName {
		vault[index].History = append(vault[index].History, originalItemName)
	}

	vault[index].Name = itemName

	items.SaveVaultToFile(vault)
}

func itemHistory(item items.Item) {
	fmt.Println("\nHistory:")

	if len(item.History) == 0 {
		fmt.Println("No history found for this item")
	}

	for i, val := range item.History {
		fmt.Printf("%d. %v\n", i, val)
	}
}

// fmt.Println("E - edit || R - Remove || I - Insert || H - History || Q = Quit")
// Creates a string of the vault menu options
func vaultMenuOptions(options []string) string {
	str := ""

	for i, option := range options {
		firstLetter := console.Blue + string(option[0]) + console.Reset
		str += firstLetter + option[1:]
		if i != len(options)-1 {
			str += " || "
		}
	}

	return str
}

func title() {
	ascii := figlet4go.NewAsciiRender()

	renderStr, _ := ascii.Render("gTodo")
	renderStr = console.Blue + renderStr + console.Reset
	renderStr += "\n" + console.Green + "ver 1.0" + console.Reset + "\n"

	fmt.Println(renderStr)
}

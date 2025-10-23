package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"todo/internal/items"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func clearCli() {
	switch runtime.GOOS {
	case "windows":
		log.Fatal("This program doesn't currently work on Windows ATM")
	default:
		fmt.Println("\033[2J")
		fmt.Print("\033[H")
	}
}

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

	clearCli()
	for itemIndex, item := range vault {
		if itemIndex == index {
			fmt.Printf(">> %v\n", item.Name)
		} else {
			fmt.Println(item.Name)
		}
	}

	fmt.Println("E - edit || R - Remove || I - Insert || H - History || Q = Quit")

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
			clearCli()
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
	clearCli()
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

	clearCli()
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

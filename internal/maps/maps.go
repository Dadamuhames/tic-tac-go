package maps

import (
	"fmt"
)

var gameMaps = make(map[string]*[3][3]rune)

func GetKey(fromId int, toId int) string {
	var key = ""

	if fromId < toId {
		key = fmt.Sprintf("%d|%d", fromId, toId)
	} else {
		key = fmt.Sprintf("%d|%d", toId, fromId)
	}

	fmt.Println(key)

	return key
}

func InitMap(fromId int, toId int) {
	key := GetKey(fromId, toId)

	_, exists := gameMaps[key]

	if exists {
		ClearBoard(fromId, toId)
	}

	gameMaps[key] = createMatrix()
}

func AddToMap(fromId, toId int, row, col int, xo string) bool {
	key := GetKey(fromId, toId)

	fmt.Printf("key: %s, row: %d, col: %d\n", key, row, col)

	boardPtr, exists := gameMaps[key]
	if !exists {
		return false
	}
	if (*boardPtr)[row][col] != ' ' {
		return false // already occupied
	}
	if xo == "X" {
		(*boardPtr)[row][col] = 'X'
	} else {
		(*boardPtr)[row][col] = 'O'
	}
	return true
}

func CheckWinner(fromId int, toId int) string {
	key := GetKey(fromId, toId)

	boardPtr, exists := gameMaps[key]

	if !exists {
		return " "
	}

	b := *boardPtr

	for i := 0; i < 3; i++ {
		if b[i][0] != ' ' && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			return string(b[i][0])
		}
	}

	// Check columns
	for j := 0; j < 3; j++ {
		if b[0][j] != ' ' && b[0][j] == b[1][j] && b[1][j] == b[2][j] {
			return string(b[0][j])
		}
	}

	// Check diagonals
	if b[0][0] != ' ' && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return string(b[0][0])
	}
	if b[0][2] != ' ' && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return string(b[0][2])
	}

	// Check for draw
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b[i][j] == ' ' {
				return ""
			}
		}
	}

	return "N"
}

func createMatrix() *([3][3]rune) {
	var matrix [3][3]rune

	matrix[0][0] = ' '
	matrix[0][1] = ' '
	matrix[0][2] = ' '
	matrix[1][0] = ' '
	matrix[1][1] = ' '
	matrix[1][2] = ' '
	matrix[2][0] = ' '
	matrix[2][1] = ' '
	matrix[2][2] = ' '

	return &matrix
}

func ClearBoard(fromId int, toId int) {
	key := GetKey(fromId, toId)
	delete(gameMaps, key)
}

func printBoard(matrix [3][3]rune) {
	fmt.Println("┌───┬───┬───┐")
	for i := 0; i < 3; i++ {
		fmt.Print("│")
		for j := 0; j < 3; j++ {
			if matrix[i][j] == ' ' {
				fmt.Print(" D ")
			} else {
				fmt.Printf(" %c ", matrix[i][j])
			}
			fmt.Print("│")
		}
		fmt.Println()
		if i < 2 {
			fmt.Println("├───┼───┼───┤")

		}
	}
	fmt.Println("└───┴───┴───┘")
}

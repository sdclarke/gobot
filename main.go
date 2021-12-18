package main

import (
	"fmt"
	"log"

	"github.com/sdclarke/gobot/pkg/game"
)

func sendMessage(msg string) {
	fmt.Print(msg)
}

func recvMessage() string {
	message := ""
	fmt.Scanln(&message)
	return message
}

func main() {
	board := game.NewBoard(7, 7)
	log.Printf("%#v", board)
	side := game.South
	canSwap := true
	minimax := game.NewMinimax(side)

	for {
		message := recvMessage()
		msgType, err := game.GetMessageType(message)
		if err != nil {
			log.Fatalf("Error %v", err)
		}
		switch msgType {
		case game.Start:
			south, err := game.InterpretStartMessage(message)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			if south {
				canSwap = false
				sendMessage(game.CreateMoveMessage(1))
			} else {
				side = game.North
				minimax.UpdateSide(side)
			}
			log.Printf("Can swap: %v, side: %v", canSwap, side)
			break
		case game.State:
			move, err := game.InterpretStateMessage(message, board)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			if move.Move == -1 {
				side = side.Opposite()
				minimax.UpdateSide(side)
			}
			if move.Again {
				if canSwap {
					side = side.Opposite()
					minimax.UpdateSide(side)
					sendMessage(game.CreateSwapMessage())
				} else {
					hole, err := minimax.GetBestMove(board)
					if err != nil {
						log.Fatalf("Error getting best move: %v", err)
					}
					if hole < 1 {
						log.Fatalf("Something went wrong")
					}
					sendMessage(game.CreateMoveMessage(hole))
				}
				canSwap = false
			}
			break
		case game.End:
			return
		}
	}
}

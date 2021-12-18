package game

import (
	"errors"
	"fmt"
	"strings"
)

type MessageType int

const (
	Start MessageType = iota
	State
	End
	Error
)

type moveTurn struct {
	End   bool
	Again bool
	Move  int
}

func GetMessageType(message string) (MessageType, error) {
	if strings.HasPrefix(message, "START") {
		return Start, nil
	} else if strings.HasPrefix(message, "CHANGE") {
		return State, nil
	} else if strings.HasPrefix(message, "END") {
		return End, nil
	}
	return Error, errors.New("Unknown message type")
}

func InterpretStartMessage(message string) (bool, error) {
	if strings.HasSuffix(message, "South") {
		return true, nil
	} else if strings.HasSuffix(message, "North") {
		return false, nil
	}
	return false, errors.New("Unknown start message format")
}

func InterpretStateMessage(message string, b *board) (*moveTurn, error) {
	m := &moveTurn{}
	s := strings.Split(message, ";")
	if len(s) != 4 {
		return nil, errors.New("Incorrect state message")
	}
	if s[1] == "SWAP" {
		m.Move = -1
	} else {
		_, err := fmt.Sscanf(s[1], "%d", &m.Move)
		if err != nil {
			return nil, err
		}
	}

	boardParts := strings.Split(s[2], ",")
	if len(boardParts) != 2*(b.getNoOfHoles()+1) {
		return nil, errors.New("Incorrect length of board in state message")
	}
	for i := 0; i < b.getNoOfHoles(); i++ {
		seeds := 0
		_, err := fmt.Sscanf(boardParts[i], "%d", &seeds)
		if err != nil {
			return nil, err
		}
		b.setSeeds(North, i+1, seeds)
	}
	seeds := 0
	_, err := fmt.Sscanf(boardParts[b.getNoOfHoles()], "%d", &seeds)
	if err != nil {
		return nil, err
	}
	b.setSeedsInStore(North, seeds)
	for i := 0; i < b.getNoOfHoles(); i++ {
		seeds := 0
		_, err := fmt.Sscanf(boardParts[i+b.getNoOfHoles()+1], "%d", &seeds)
		if err != nil {
			return nil, err
		}
		b.setSeeds(South, i+1, seeds)
	}
	_, err = fmt.Sscanf(boardParts[2*b.getNoOfHoles()+1], "%d", &seeds)
	if err != nil {
		return nil, err
	}
	b.setSeedsInStore(South, seeds)

	m.End = false
	if s[3] == "YOU" {
		m.Again = true
	} else if s[3] == "OPP" {
		m.Again = false
	} else if s[3] == "END" {
		m.End = true
		m.Again = false
	} else {
		return nil, errors.New("Incorrect end of state message")
	}
	return m, nil
}

func CreateMoveMessage(hole int) string {
	return fmt.Sprintf("MOVE;%d\n", hole)
}

func CreateSwapMessage() string {
	return "SWAP\n"
}

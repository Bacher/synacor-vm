package main

import (
	"bufio"
	"log"
	"os"
)

var command []byte = nil

func run(m *Memory) {
	for {
		opCode := m.readOpCode()

		if opCode == 4 || opCode == 5 || opCode >= 9 && opCode <= 13 {
			a := m.readRegister()
			b := m.readValue()
			c := m.readValue()

			var res uint16

			switch opCode {
			case 4:
				res = bool2uint16(b == c)
			case 5:
				res = bool2uint16(b > c)
			case 9:
				res = b + c
			case 10:
				res = b * c
			case 11:
				res = b % c
			case 12:
				res = b & c
			case 13:
				res = b | c
			}

			m.setRegister(a, res)

		} else {
			switch opCode {
			case 1:
				a := m.readRegister()
				b := m.readValue()

				m.setRegister(a, b)
			case 2:
				a := m.readValue()

				m.stackPush(a)
			case 3:
				a := m.readRegister()

				m.setRegister(a, m.stackPop())
			case 6:
				a := m.readAddress()
				m.goTo(a)
			case 7:
				a := m.readValue()
				b := m.readAddress()

				if a > 0 {
					m.goTo(b)
				}
			case 8:
				a := m.readValue()
				b := m.readAddress()

				if a == 0 {
					m.goTo(b)
				}
			case 14:
				a := m.readRegister()
				b := m.readValue()

				m.setRegister(a, b^32767)
			case 15:
				a := m.readRegister()
				b := m.readAddress()

				m.setRegister(a, m.getValueByAddress(b))
			case 16:
				a := m.readAddress()
				b := m.readValue()

				m.setValueByAddress(a, b)

			case 17:
				a := m.readAddress()

				m.stackPush(uint16(m.index))
				m.goTo(a)
			case 18:
				address := Address(m.stackPop())
				m.goTo(address)
			case 19:
				a := m.readValue()

				os.Stdout.WriteString(string(a))
			case 20:
				a := m.readRegister()

				if len(command) == 0 {
					os.Stdout.WriteString("> ")

					reader := bufio.NewReader(os.Stdin)

					text, err := reader.ReadString('\n')

					if err != nil {
						log.Fatalf("Read line failed %v", err)
					}

					command = []byte(text)
				}

				m.setRegister(a, uint16(command[0]))
				command = command[1:]
			case 21:
				// noop
			default:
				log.Panicf("Unsupported OpCode {%d}", opCode)
			}
		}
	}
}

func bool2uint16(val bool) uint16 {
	if val {
		return 1
	} else {
		return 0
	}
}

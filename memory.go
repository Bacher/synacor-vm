package main

import (
	"encoding/binary"
	"io/ioutil"
	"log"
)

const (
	MAX_NUMBER          = 32767
	NUMBER_BOUND        = 32768
	FIRST_REGISTER      = 32768
	LAST_REGISTER       = 32775
	INVALID_VALUE_START = 32776
	MAX_ADDRESS         = 32767
)

type OpCode uint16
type Register uint16
type Address uint16

type Memory struct {
	data      []uint16
	registers []uint16
	stack     []uint16
	stackI    int
	index     Address
}

func loadMemoryFromFile(fileName string) *Memory {
	application, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Panicf("Program not found. %v", err)
		return nil
	}

	data := make([]uint16, 32768)

	for i := 0; i < len(application)/2; i++ {
		shift := i * 2

		data[i] = binary.LittleEndian.Uint16(application[shift : shift+2])
	}

	memory := &Memory{data, []uint16{0, 0, 0, 0, 0, 0, 0, 0}, make([]uint16, 100), 0, 0}

	return memory
}

func (m *Memory) readInt() uint16 {
	value := m.data[m.index]
	m.index++

	return value
}

func (m *Memory) readOpCode() OpCode {
	opCode := m.readInt()

	if opCode > 22 {
		log.Panicf("Invalid opCode %d", opCode)
	}

	return OpCode(opCode)
}

func (m *Memory) readValue() uint16 {
	value := m.getValueByAddress(m.index)

	m.index++

	return value
}

func (m *Memory) readAddress() Address {
	address := m.readValue()

	if address > MAX_ADDRESS {
		log.Fatalf("Invalid address: {%d}", address)
	}

	return Address(address)
}

func (m *Memory) readRegister() Register {
	register := m.readInt()

	if register < FIRST_REGISTER || register > LAST_REGISTER {
		log.Fatalf("Invalid register {%d}", register)
	}

	return Register(register - FIRST_REGISTER)
}

func (m *Memory) setRegister(a Register, b uint16) {
	m.registers[a] = b % NUMBER_BOUND
}

func (m *Memory) goTo(address Address) {
	m.index = address
}

func (m *Memory) getValueByAddress(address Address) uint16 {
	value := m.data[address]

	if value >= INVALID_VALUE_START {
		log.Panicf("Invalid value {%v}", value)
	}

	if value <= MAX_NUMBER {
		return value
	} else {
		return m.registers[value-FIRST_REGISTER]
	}
}

func (m *Memory) setValueByAddress(address Address, value uint16) {
	m.data[address] = value
}

func (m *Memory) stackPush(value uint16) {
	if m.stackI == len(m.stack) {
		newStack := make([]uint16, len(m.stack)*2)
		copy(newStack, m.stack)
		m.stack = newStack
	}

	m.stack[m.stackI] = value
	m.stackI++
}

func (m *Memory) stackPop() uint16 {
	if m.stackI == 0 {
		log.Fatal("Drain from empty stack")
	}

	m.stackI--

	return m.stack[m.stackI]
}

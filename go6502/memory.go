package go6502

type MemoryMapEntry interface {
	WithinRange(address uint16) bool
	Get(address uint16) uint8
	Set(address uint16, value uint8)
}

type RAM struct {
	data         []uint8
	addressStart uint16
	size         uint16
}

func (R *RAM) WithinRange(address uint16) bool {
	return address >= R.addressStart && uint32(address) < uint32(R.addressStart)+uint32(R.size)
}

func (R *RAM) Get(address uint16) uint8 {
	value := R.data[address-R.addressStart]
	//log.Printf("Reading from address %04X. Value: %02X", address, value)
	return value
}

func (R *RAM) Set(address uint16, value uint8) {
	//log.Printf("Writing to address %04X. Value: %02X", address, value)
	R.data[address-R.addressStart] = value
}

func NewRAM(addressStart uint16, size uint16) *RAM {
	return &RAM{
		data:         make([]uint8, size),
		addressStart: addressStart,
		size:         size,
	}
}

type Memory struct {
	entries []MemoryMapEntry
}

func (m *Memory) Get(address uint16) uint8 {
	for _, entry := range m.entries {
		// TODO: Refactor memory, so only it is aware about memory mapping
		// Mapped memory like RAM or Screen should get only "internal" address
		// Which will be address to get minus mapping address start
		if entry.WithinRange(address) {
			return entry.Get(address)
		}
	}
	return 0xFF
}

func (m *Memory) Set(address uint16, value ...uint8) {
	for i := 0; i < len(value); i++ {
		valueAddress := address + uint16(i)
		for _, entry := range m.entries {
			if entry.WithinRange(valueAddress) {
				entry.Set(valueAddress, value[i])
				break
			}
		}
	}
}

func DefaultMemory() *Memory {
	return &Memory{
		entries: []MemoryMapEntry{
			NewRAM(0, 65535),
		},
	}
}

func NewMemory(entries ...MemoryMapEntry) *Memory {
	return &Memory{
		entries: entries,
	}
}

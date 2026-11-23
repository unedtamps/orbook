package orderbook

import "math/bits"

const (
	orderMapMinCap = 2048
	fibMul64       = 11400714819323198485 // floor(2^64 / phi)
)

// mapSlot is 16 bytes and contains zero Go pointers (GC noscan).
type mapSlot struct {
	key OrderID // 8 B (0 = empty)
	idx uint32  // 4 B (arena index >= 1)
	_   uint32  // pad to 16 B
}

// orderMap is a flat open-addressed linear-probing hash table using Fibonacci
// hashing and backward-shift deletion (zero tombstones, zero pointers).
type orderMap struct {
	slots []mapSlot
	mask  uint64
	shift uint8
	count int
}

func newOrderMap(initCap int) orderMap {
	capPow2 := orderMapMinCap
	for capPow2 < initCap {
		capPow2 <<= 1
	}
	return orderMap{
		slots: make([]mapSlot, capPow2),
		mask:  uint64(capPow2 - 1),
		shift: uint8(64 - bits.TrailingZeros64(uint64(capPow2))),
	}
}

func (m *orderMap) hash(id OrderID) uint64 {
	x := uint64(id)
	return ((x ^ (x >> 32)) * fibMul64) >> m.shift
}

func (m *orderMap) Has(id OrderID) bool {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			return true
		}
		if k == 0 {
			return false
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) Get(id OrderID) (uint32, bool) {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			return m.slots[i].idx, true
		}
		if k == 0 {
			return 0, false
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) Set(id OrderID, idx uint32) {
	if m.count*2 >= len(m.slots) {
		m.grow()
	}
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == 0 {
			m.slots[i] = mapSlot{key: id, idx: idx}
			m.count++
			return
		}
		if k == id {
			m.slots[i].idx = idx
			return
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) Delete(id OrderID) (uint32, bool) {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			idx := m.slots[i].idx
			m.count--
			hole := i
			j := (i + 1) & mask
			for {
				jk := m.slots[j].key
				if jk == 0 {
					m.slots[hole] = mapSlot{}
					return idx, true
				}
				home := m.hash(jk)
				if ((hole - home) & mask) < ((j - home) & mask) {
					m.slots[hole] = m.slots[j]
					hole = j
				}
				j = (j + 1) & mask
			}
		}
		if k == 0 {
			return 0, false
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) grow() {
	oldSlots := m.slots
	newCap := len(oldSlots) * 2
	m.slots = make([]mapSlot, newCap)
	m.mask = uint64(newCap - 1)
	m.shift = uint8(64 - bits.TrailingZeros64(uint64(newCap)))

	mask := m.mask
	for _, s := range oldSlots {
		if s.key == 0 {
			continue
		}
		i := m.hash(s.key)
		for m.slots[i].key != 0 {
			i = (i + 1) & mask
		}
		m.slots[i] = s
	}
}

func (m *orderMap) Clear() {
	if m.count == 0 {
		return
	}
	clear(m.slots)
	m.count = 0
}

package book

import "math/bits"

const (
	orderMapInitCap = 2048
	fibMul64        = 11400714819323198485
)

type orderSlot struct {
	key OrderID
	val *Order
}

type orderMap struct {
	slots []orderSlot
	mask  uint64
	shift uint8
	count int
}

func newOrderMap(initCap int) orderMap {
	capPow2 := orderMapInitCap
	for capPow2 < initCap {
		capPow2 <<= 1
	}
	return orderMap{
		slots: make([]orderSlot, capPow2),
		mask:  uint64(capPow2 - 1),
		shift: uint8(64 - bits.TrailingZeros64(uint64(capPow2))),
	}
}

func (m *orderMap) hash(id OrderID) uint64 {
	x := uint64(id)
	return ((x ^ (x >> 32)) * fibMul64) >> m.shift
}

func (m *orderMap) Len() int {
	return m.count
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

func (m *orderMap) Get(id OrderID) (*Order, bool) {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			return m.slots[i].val, true
		}
		if k == 0 {
			return nil, false
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) Set(id OrderID, val *Order) {
	if m.count*2 >= len(m.slots) {
		m.grow()
	}
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == 0 {
			m.slots[i] = orderSlot{key: id, val: val}
			m.count++
			return
		}
		if k == id {
			m.slots[i].val = val
			return
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) Delete(id OrderID) (*Order, bool) {
	mask := m.mask
	i := m.hash(id)
	for {
		k := m.slots[i].key
		if k == id {
			val := m.slots[i].val
			m.count--
			hole := i
			j := (i + 1) & mask
			for {
				jk := m.slots[j].key
				if jk == 0 {
					m.slots[hole] = orderSlot{}
					return val, true
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
			return nil, false
		}
		i = (i + 1) & mask
	}
}

func (m *orderMap) grow() {
	oldSlots := m.slots
	newCap := len(oldSlots) * 2
	m.slots = make([]orderSlot, newCap)
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

func (m *orderMap) ForEach(fn func(*Order)) {
	if m.count == 0 {
		return
	}
	for i := range m.slots {
		if m.slots[i].key != 0 {
			fn(m.slots[i].val)
		}
	}
}

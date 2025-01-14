package mem

type MemoryKind string

const (
	MemoryKindRAM  = "RAM"
	MemoryKindSwap = "SWAP"
)

type Memory struct {
	ID       string
	Kind     MemoryKind
	Capacity float32
	Used     float32
}

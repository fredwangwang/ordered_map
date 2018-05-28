package orderedmap

type KeyElement struct {
	Key   string
	Index int
}

type KeyIndices []KeyElement

func (a KeyIndices) Len() int           { return len(a) }
func (a KeyIndices) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a KeyIndices) Less(i, j int) bool { return a[i].Index < a[j].Index }

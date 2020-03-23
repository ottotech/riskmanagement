package listing

type Risk struct {
	ID             string
	RiskMatrixID   int
	Name           string
	Probability    int
	Impact         int
	Classification string
	Strategy       string
}

// ByAge implements sort.Interface for []Person based on
// the Age field.
type ByName []Risk

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

package mapping

// Segmentation is the main entry point to the MCH mapping
type Segmentation interface {
	IsValid(padid int) bool
	FindPadByFEE(dualSampaID int, dualSampaChannel int) (int, error)
	FindPadByPosition(x float64, y float64) (int, error)
}

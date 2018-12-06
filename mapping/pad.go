package mapping

// PadFEELocator retrieves DualSampa ID and channel information
// from a PadUID (within a given detection element).
// Note that mapping.Segmentation interface satisfy PadFEELocator.
type PadFEELocator interface {
	PadDualSampaChannel(paduid PadUID) int
	PadDualSampaID(paduid PadUID) DualSampaID
}

// PadPositioner returns the x and y positions (in cm)
// of a pad, relative to detection element origin.
// Note that mapping.Segmentation interface satisfy PadFEELocator.
type PadPositioner interface {
	PadPositionX(paduid PadUID) float64
	PadPositionY(paduid PadUID) float64
}

// PadSizer returns the x and y sizes (in cm) of a pd.
// Note that mapping.Segmentation interface satisfy PadFEELocator.
type PadSizer interface {
	PadSizeX(paduid PadUID) float64
	PadSizeY(paduid PadUID) float64
}

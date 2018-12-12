package mapping

// PadFEELocator retrieves DualSampa ID and channel information
// from a PadUID (within a given detection element).
// Note that mapping.Segmentation interface satisfy PadFEELocator.
type PadFEELocator interface {
	PadDualSampaChannel(paduid PadUID) DualSampaChannelID
	PadDualSampaID(paduid PadUID) DualSampaID
}

// PadFEELocatorFunc is any function that can create
// a PadFEELocator from a detection element id.
type PadFEELocatorFunc func(deid DEID) PadFEELocator

// PadSizerPositioner is able to return the x and y positions (in cm)
// of a pad, relative to detection element origin,
// as well as its x and y sizes (in cm).
// Note that mapping.Segmentation interface satisfy PadSizerPositioner.
type PadSizerPositioner interface {
	PadPositionX(paduid PadUID) float64
	PadPositionY(paduid PadUID) float64
	PadSizeX(paduid PadUID) float64
	PadSizeY(paduid PadUID) float64
}

// PadByFEEFinder finds a pad by its front-end electronics
// (FEE) identifiers.
type PadByFEEFinder interface {
	FindPadByFEE(DualSampaID, DualSampaChannelID) (PadUID, error)
}

// PadByFEEFinderFunc is any function that can create
// a PadByFEEFinder from a detection element id.
type PadByFEEFinderFunc func(deid DEID) PadByFEEFinder

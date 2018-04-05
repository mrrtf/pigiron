package mapping

// Segmentation is the main entry point to the MCH mapping
type Segmentation interface {
	NofPads() int
	NofDualSampas() int
	DualSampaID(dualSampaIndex int) (int, error)
	IsValid(padid int) bool
	FindPadByFEE(dualSampaID int, dualSampaChannel int) (int, error)
	FindPadByPosition(x float64, y float64) (int, error)
	ForEachPadInDualSampa(dualSampaID int, padHandler func(paduid int))
}

// ForEachDetectionElement loops over all detection elements and call the detElemIdHandler function
// for each of them
func ForEachDetectionElement(detElemIdHandler func(detElemID int)) {
	for _, detElemID := range []int{100, 101, 102, 103, 200, 201, 202, 203, 300, 301, 302, 303, 400, 401, 402, 403, 500, 501,
		502, 503, 504, 505, 506, 507, 508, 509, 510, 511, 512, 513, 514, 515, 516, 517, 600, 601,
		602, 603, 604, 605, 606, 607, 608, 609, 610, 611, 612, 613, 614, 615, 616, 617, 700, 701,
		702, 703, 704, 705, 706, 707, 708, 709, 710, 711, 712, 713, 714, 715, 716, 717, 718, 719,
		720, 721, 722, 723, 724, 725, 800, 801, 802, 803, 804, 805, 806, 807, 808, 809, 810, 811,
		812, 813, 814, 815, 816, 817, 818, 819, 820, 821, 822, 823, 824, 825, 900, 901, 902, 903,
		904, 905, 906, 907, 908, 909, 910, 911, 912, 913, 914, 915, 916, 917, 918, 919, 920, 921,
		922, 923, 924, 925, 1000, 1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010, 1011, 1012, 1013,
		1014, 1015, 1016, 1017, 1018, 1019, 1020, 1021, 1022, 1023, 1024, 1025} {
		detElemIdHandler(detElemID)
	}
}

// ForOneDetectionElementOfEachSegmentationType loops over one detection element per segmentation type
// and call the detElemIdHandler function for each of them
func ForOneDetectionElementOfEachSegmentationType(detElemIdHandler func(detElemID int)) {
	for _, detElemID := range []int{100, 300, 500, 501, 502, 503, 504, 600, 601, 602, 700, 701, 702, 703, 704, 705, 706, 902, 903, 904, 905} {
		detElemIdHandler(detElemID)
	}
}

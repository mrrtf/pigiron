package impl4

import "github.com/mrrtf/pigiron/mapping"

type createSegType17 struct{}

func (seg createSegType17) Build(isBendingPlane bool, deid mapping.DEID) mapping.CathodeSegmentation {
	if isBendingPlane {
		return newCathodeSegmentation(deid, 17, true,
			[]padGroup{
				{1, 3, 0, -120, -20},
				{2, 12, 0, -117.5, -20},
				{3, 6, 0, -112.5, -20},
				{4, 13, 0, -110, -20},
				{5, 4, 0, -105, -20},
				{6, 3, 0, -100, -20},
				{7, 12, 0, -97.5, -20},
				{8, 6, 0, -92.5, -20},
				{9, 13, 0, -90, -20},
				{10, 4, 0, -85, -20},
				{18, 3, 0, -80, -20},
				{19, 12, 0, -77.5, -20},
				{20, 6, 0, -72.5, -20},
				{21, 13, 0, -70, -20},
				{22, 4, 0, -65, -20},
				{23, 3, 0, -60, -20},
				{24, 12, 0, -57.5, -20},
				{25, 6, 0, -52.5, -20},
				{26, 13, 0, -50, -20},
				{27, 4, 0, -45, -20},
				{35, 3, 1, -40, -20},
				{36, 12, 1, -35, -20},
				{37, 6, 1, -25, -20},
				{38, 13, 1, -20, -20},
				{39, 4, 1, -10, -20},
				{103, 8, 2, 80, -20},
				{104, 8, 2, 100, -20},
				{107, 8, 2, 40, -20},
				{108, 8, 2, 60, -20},
				{112, 3, 1, 0, -20},
				{113, 12, 1, 5, -20},
				{114, 6, 1, 15, -20},
				{115, 13, 1, 20, -20},
				{116, 4, 1, 30, -20},
				{201, 0, 2, 100, -4},
				{202, 7, 2, 90, 4},
				{203, 5, 2, 80, -4},
				{206, 0, 2, 60, -4},
				{207, 7, 2, 50, 4},
				{208, 5, 2, 40, -4},
				{211, 2, 1, 30, 0},
				{212, 11, 1, 20, 0},
				{213, 9, 1, 15, 4},
				{214, 10, 1, 5, 0},
				{215, 1, 1, 0, 0},
				{308, 2, 0, -85, 0},
				{309, 11, 0, -90, 0},
				{310, 9, 0, -92.5, 4},
				{311, 10, 0, -97.5, 0},
				{312, 1, 0, -100, 0},
				{313, 2, 0, -105, 0},
				{314, 11, 0, -110, 0},
				{315, 9, 0, -112.5, 4},
				{316, 10, 0, -117.5, 0},
				{317, 1, 0, -120, 0},
				{325, 2, 0, -45, 0},
				{326, 11, 0, -50, 0},
				{327, 9, 0, -52.5, 4},
				{328, 10, 0, -57.5, 0},
				{329, 1, 0, -60, 0},
				{330, 2, 0, -65, 0},
				{331, 11, 0, -70, 0},
				{332, 9, 0, -72.5, 4},
				{333, 10, 0, -77.5, 0},
				{334, 1, 0, -80, 0},
				{338, 2, 1, -10, 0},
				{339, 11, 1, -20, 0},
				{340, 9, 1, -25, 4},
				{341, 10, 1, -35, 0},
				{342, 1, 1, -40, 0},
			},
			[]padGroupType{
				/* L10 */ NewPadGroupType(2, 48, []int{15, 16, 14, 17, 13, 18, 12, 19, 11, 20, 10, 21, 9, 22, 8, 23, 7, 24, 6, 25, 5, 26, 4, 27, 3, 28, 2, 29, 1, 30, 0, 31, -1, 48, -1, 49, -1, 50, -1, 51, -1, 52, -1, 53, -1, 54, -1, 55, -1, 56, -1, 57, -1, 58, -1, 59, -1, 60, -1, 61, -1, 62, -1, 63, -1, 32, -1, 33, -1, 34, -1, 35, -1, 36, -1, 37, -1, 38, -1, 39, -1, 40, -1, 41, -1, 42, -1, 43, -1, 44, -1, 45, -1, 46, -1, 47}),
				/* L5 */ NewPadGroupType(2, 40, []int{55, 56, 54, 57, 53, 58, 52, 59, 51, 60, 50, 61, 49, 62, 48, 63, 31, 32, 30, 33, 29, 34, 28, 35, 27, 36, 26, 37, 25, 38, 24, 39, 23, 40, 22, 41, 21, 42, 20, 43, 19, 44, 18, 45, 17, 46, 16, 47, 15, -1, 14, -1, 13, -1, 12, -1, 11, -1, 10, -1, 9, -1, 8, -1, 7, -1, 6, -1, 5, -1, 4, -1, 3, -1, 2, -1, 1, -1, 0, -1}),
				/* L6 */ NewPadGroupType(2, 40, []int{23, 24, 22, 25, 21, 26, 20, 27, 19, 28, 18, 29, 17, 30, 16, 31, 15, 48, 14, 49, 13, 50, 12, 51, 11, 52, 10, 53, 9, 54, 8, 55, 7, 56, 6, 57, 5, 58, 4, 59, 3, 60, 2, 61, 1, 62, 0, 63, -1, 32, -1, 33, -1, 34, -1, 35, -1, 36, -1, 37, -1, 38, -1, 39, -1, 40, -1, 41, -1, 42, -1, 43, -1, 44, -1, 45, -1, 46, -1, 47}),
				/* L7 */ NewPadGroupType(2, 40, []int{47, -1, 46, -1, 45, -1, 44, -1, 43, -1, 42, -1, 41, -1, 40, -1, 39, -1, 38, -1, 37, -1, 36, -1, 35, -1, 34, -1, 33, -1, 32, -1, 63, 0, 62, 1, 61, 2, 60, 3, 59, 4, 58, 5, 57, 6, 56, 7, 55, 8, 54, 9, 53, 10, 52, 11, 51, 12, 50, 13, 49, 14, 48, 15, 31, 16, 30, 17, 29, 18, 28, 19, 27, 20, 26, 21, 25, 22, 24, 23}),
				/* L8 */ NewPadGroupType(2, 40, []int{-1, 0, -1, 1, -1, 2, -1, 3, -1, 4, -1, 5, -1, 6, -1, 7, -1, 8, -1, 9, -1, 10, -1, 11, -1, 12, -1, 13, -1, 14, -1, 15, 47, 16, 46, 17, 45, 18, 44, 19, 43, 20, 42, 21, 41, 22, 40, 23, 39, 24, 38, 25, 37, 26, 36, 27, 35, 28, 34, 29, 33, 30, 32, 31, 63, 48, 62, 49, 61, 50, 60, 51, 59, 52, 58, 53, 57, 54, 56, 55}),
				/* L9 */ NewPadGroupType(2, 48, []int{63, 32, 62, 33, 61, 34, 60, 35, 59, 36, 58, 37, 57, 38, 56, 39, 55, 40, 54, 41, 53, 42, 52, 43, 51, 44, 50, 45, 49, 46, 48, 47, 31, -1, 30, -1, 29, -1, 28, -1, 27, -1, 26, -1, 25, -1, 24, -1, 23, -1, 22, -1, 21, -1, 20, -1, 19, -1, 18, -1, 17, -1, 16, -1, 15, -1, 14, -1, 13, -1, 12, -1, 11, -1, 10, -1, 9, -1, 8, -1, 7, -1, 6, -1, 5, -1, 4, -1, 3, -1, 2, -1, 1, -1, 0, -1}),
				/* O10 */ NewPadGroupType(2, 32, []int{48, 31, 49, 30, 50, 29, 51, 28, 52, 27, 53, 26, 54, 25, 55, 24, 56, 23, 57, 22, 58, 21, 59, 20, 60, 19, 61, 18, 62, 17, 63, 16, 32, 15, 33, 14, 34, 13, 35, 12, 36, 11, 37, 10, 38, 9, 39, 8, 40, 7, 41, 6, 42, 5, 43, 4, 44, 3, 45, 2, 46, 1, 47, 0}),
				/* O11 */ NewPadGroupType(2, 32, []int{31, 48, 30, 49, 29, 50, 28, 51, 27, 52, 26, 53, 25, 54, 24, 55, 23, 56, 22, 57, 21, 58, 20, 59, 19, 60, 18, 61, 17, 62, 16, 63, 15, 32, 14, 33, 13, 34, 12, 35, 11, 36, 10, 37, 9, 38, 8, 39, 7, 40, 6, 41, 5, 42, 4, 43, 3, 44, 2, 45, 1, 46, 0, 47}),
				/* O12 */ NewPadGroupType(2, 32, []int{47, 0, 46, 1, 45, 2, 44, 3, 43, 4, 42, 5, 41, 6, 40, 7, 39, 8, 38, 9, 37, 10, 36, 11, 35, 12, 34, 13, 33, 14, 32, 15, 63, 16, 62, 17, 61, 18, 60, 19, 59, 20, 58, 21, 57, 22, 56, 23, 55, 24, 54, 25, 53, 26, 52, 27, 51, 28, 50, 29, 49, 30, 48, 31}),
				/* O9 */ NewPadGroupType(2, 32, []int{0, 47, 1, 46, 2, 45, 3, 44, 4, 43, 5, 42, 6, 41, 7, 40, 8, 39, 9, 38, 10, 37, 11, 36, 12, 35, 13, 34, 14, 33, 15, 32, 16, 63, 17, 62, 18, 61, 19, 60, 20, 59, 21, 58, 22, 57, 23, 56, 24, 55, 25, 54, 26, 53, 27, 52, 28, 51, 29, 50, 30, 49, 31, 48}),
				/* Z1 */ NewPadGroupType(3, 40, []int{-1, 39, 40, -1, 38, 41, -1, 37, 42, -1, 36, 43, -1, 35, 44, -1, 34, 45, -1, 33, 46, -1, 32, 47, -1, 63, -1, -1, 62, -1, -1, 61, -1, -1, 60, -1, -1, 59, -1, -1, 58, -1, -1, 57, -1, -1, 56, -1, -1, 55, -1, -1, 54, -1, -1, 53, -1, -1, 52, -1, -1, 51, -1, -1, 50, -1, -1, 49, -1, -1, 48, -1, 0, 31, -1, 1, 30, -1, 2, 29, -1, 3, 28, -1, 4, 27, -1, 5, 26, -1, 6, 25, -1, 7, 24, -1, 8, 23, -1, 9, 22, -1, 10, 21, -1, 11, 20, -1, 12, 19, -1, 13, 18, -1, 14, 17, -1, 15, 16, -1}),
				/* Z2 */ NewPadGroupType(3, 40, []int{7, 8, -1, 6, 9, -1, 5, 10, -1, 4, 11, -1, 3, 12, -1, 2, 13, -1, 1, 14, -1, 0, 15, -1, -1, 16, -1, -1, 17, -1, -1, 18, -1, -1, 19, -1, -1, 20, -1, -1, 21, -1, -1, 22, -1, -1, 23, -1, -1, 24, -1, -1, 25, -1, -1, 26, -1, -1, 27, -1, -1, 28, -1, -1, 29, -1, -1, 30, -1, -1, 31, -1, -1, 48, 47, -1, 49, 46, -1, 50, 45, -1, 51, 44, -1, 52, 43, -1, 53, 42, -1, 54, 41, -1, 55, 40, -1, 56, 39, -1, 57, 38, -1, 58, 37, -1, 59, 36, -1, 60, 35, -1, 61, 34, -1, 62, 33, -1, 63, 32}),
				/* Z3 */ NewPadGroupType(3, 40, []int{32, 63, -1, 33, 62, -1, 34, 61, -1, 35, 60, -1, 36, 59, -1, 37, 58, -1, 38, 57, -1, 39, 56, -1, 40, 55, -1, 41, 54, -1, 42, 53, -1, 43, 52, -1, 44, 51, -1, 45, 50, -1, 46, 49, -1, 47, 48, -1, -1, 31, -1, -1, 30, -1, -1, 29, -1, -1, 28, -1, -1, 27, -1, -1, 26, -1, -1, 25, -1, -1, 24, -1, -1, 23, -1, -1, 22, -1, -1, 21, -1, -1, 20, -1, -1, 19, -1, -1, 18, -1, -1, 17, -1, -1, 16, -1, -1, 15, 0, -1, 14, 1, -1, 13, 2, -1, 12, 3, -1, 11, 4, -1, 10, 5, -1, 9, 6, -1, 8, 7}),
				/* Z4 */ NewPadGroupType(3, 40, []int{-1, 16, 15, -1, 17, 14, -1, 18, 13, -1, 19, 12, -1, 20, 11, -1, 21, 10, -1, 22, 9, -1, 23, 8, -1, 24, 7, -1, 25, 6, -1, 26, 5, -1, 27, 4, -1, 28, 3, -1, 29, 2, -1, 30, 1, -1, 31, 0, -1, 48, -1, -1, 49, -1, -1, 50, -1, -1, 51, -1, -1, 52, -1, -1, 53, -1, -1, 54, -1, -1, 55, -1, -1, 56, -1, -1, 57, -1, -1, 58, -1, -1, 59, -1, -1, 60, -1, -1, 61, -1, -1, 62, -1, -1, 63, -1, 47, 32, -1, 46, 33, -1, 45, 34, -1, 44, 35, -1, 43, 36, -1, 42, 37, -1, 41, 38, -1, 40, 39, -1}),
			},
			[]padSize{
				{2.5, 0.5},
				{5, 0.5},
				{10, 0.5},
			})
	}
	return newCathodeSegmentation(deid, 17, false,
		[]padGroup{
			{1035, 5, 0, -85.7142868, -20},
			{1036, 5, 0, -91.42857361, -20},
			{1037, 5, 0, -97.14286041, -20},
			{1038, 5, 0, -102.8571396, -20},
			{1039, 5, 0, -108.5714264, -20},
			{1040, 5, 0, -114.2857132, -20},
			{1041, 5, 0, -120, -20},
			{1052, 5, 0, -45.7142868, -20},
			{1053, 5, 0, -51.42856979, -20},
			{1054, 5, 0, -57.1428566, -20},
			{1055, 5, 0, -62.8571434, -20},
			{1056, 5, 0, -68.57142639, -20},
			{1057, 5, 0, -74.2857132, -20},
			{1058, 5, 0, -80, -20},
			{1064, 15, 1, -10, -20},
			{1065, 19, 1, -20, -20},
			{1066, 18, 1, -31.4285717, -20},
			{1067, 14, 1, -40, -20},
			{1125, 11, 2, 100, -20},
			{1126, 10, 2, 80, -20},
			{1129, 11, 2, 60, -20},
			{1130, 10, 2, 40, -20},
			{1133, 1, 1, 25.7142849, -20},
			{1134, 6, 1, 14.28571415, -20},
			{1135, 0, 1, 4.440892099e-15, -20},
			{1228, 8, 2, 80, 0},
			{1229, 9, 2, 100, 0},
			{1233, 8, 2, 40, 0},
			{1234, 9, 2, 60, 0},
			{1240, 16, 1, -7.105427358e-15, -5},
			{1241, 12, 1, 8.571428299, -5},
			{1242, 13, 1, 20, -5},
			{1243, 17, 1, 30, -5},
			{1325, 4, 0, -120, 0},
			{1326, 4, 0, -114.2857132, 0},
			{1327, 4, 0, -108.5714264, 0},
			{1328, 4, 0, -102.8571396, 0},
			{1329, 4, 0, -97.14286041, 0},
			{1330, 4, 0, -91.42857361, 0},
			{1331, 4, 0, -85.7142868, 0},
			{1342, 4, 0, -80, 0},
			{1343, 4, 0, -74.2857132, 0},
			{1344, 4, 0, -68.57142639, 0},
			{1345, 4, 0, -62.8571434, 0},
			{1346, 4, 0, -57.1428566, 0},
			{1347, 4, 0, -51.42856979, 0},
			{1348, 4, 0, -45.7142868, 0},
			{1359, 2, 1, -40, 0},
			{1360, 7, 1, -25.7142849, 0},
			{1361, 3, 1, -14.28571415, 0},
		},
		[]padGroupType{
			/* L1 */ NewPadGroupType(20, 4, []int{3, 7, 11, 15, 18, 21, 24, 27, 30, 49, 52, 55, 58, 61, 32, 35, 38, 41, 44, 47, 2, 6, 10, 14, 17, 20, 23, 26, 29, 48, 51, 54, 57, 60, 63, 34, 37, 40, 43, 46, 1, 5, 9, 13, 16, 19, 22, 25, 28, 31, 50, 53, 56, 59, 62, 33, 36, 39, 42, 45, 0, 4, 8, 12, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}),
			/* L2 */ NewPadGroupType(20, 4, []int{2, 5, 8, 11, 14, 17, 20, 23, 26, 29, 48, 51, 54, 57, 60, 63, 35, 39, 43, 47, 1, 4, 7, 10, 13, 16, 19, 22, 25, 28, 31, 50, 53, 56, 59, 62, 34, 38, 42, 46, 0, 3, 6, 9, 12, 15, 18, 21, 24, 27, 30, 49, 52, 55, 58, 61, 33, 37, 41, 45, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 32, 36, 40, 44}),
			/* L3 */ NewPadGroupType(20, 4, []int{44, 40, 36, 32, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 45, 41, 37, 33, 61, 58, 55, 52, 49, 30, 27, 24, 21, 18, 15, 12, 9, 6, 3, 0, 46, 42, 38, 34, 62, 59, 56, 53, 50, 31, 28, 25, 22, 19, 16, 13, 10, 7, 4, 1, 47, 43, 39, 35, 63, 60, 57, 54, 51, 48, 29, 26, 23, 20, 17, 14, 11, 8, 5, 2}),
			/* L4 */ NewPadGroupType(20, 4, []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 12, 8, 4, 0, 45, 42, 39, 36, 33, 62, 59, 56, 53, 50, 31, 28, 25, 22, 19, 16, 13, 9, 5, 1, 46, 43, 40, 37, 34, 63, 60, 57, 54, 51, 48, 29, 26, 23, 20, 17, 14, 10, 6, 2, 47, 44, 41, 38, 35, 32, 61, 58, 55, 52, 49, 30, 27, 24, 21, 18, 15, 11, 7, 3}),
			/* O1 */ NewPadGroupType(8, 8, []int{40, 32, 56, 48, 24, 16, 8, 0, 41, 33, 57, 49, 25, 17, 9, 1, 42, 34, 58, 50, 26, 18, 10, 2, 43, 35, 59, 51, 27, 19, 11, 3, 44, 36, 60, 52, 28, 20, 12, 4, 45, 37, 61, 53, 29, 21, 13, 5, 46, 38, 62, 54, 30, 22, 14, 6, 47, 39, 63, 55, 31, 23, 15, 7}),
			/* O2 */ NewPadGroupType(8, 8, []int{7, 15, 23, 31, 55, 63, 39, 47, 6, 14, 22, 30, 54, 62, 38, 46, 5, 13, 21, 29, 53, 61, 37, 45, 4, 12, 20, 28, 52, 60, 36, 44, 3, 11, 19, 27, 51, 59, 35, 43, 2, 10, 18, 26, 50, 58, 34, 42, 1, 9, 17, 25, 49, 57, 33, 41, 0, 8, 16, 24, 48, 56, 32, 40}),
			/* O3 */ NewPadGroupType(16, 4, []int{3, 7, 11, 15, 19, 23, 27, 31, 51, 55, 59, 63, 35, 39, 43, 47, 2, 6, 10, 14, 18, 22, 26, 30, 50, 54, 58, 62, 34, 38, 42, 46, 1, 5, 9, 13, 17, 21, 25, 29, 49, 53, 57, 61, 33, 37, 41, 45, 0, 4, 8, 12, 16, 20, 24, 28, 48, 52, 56, 60, 32, 36, 40, 44}),
			/* O4 */ NewPadGroupType(16, 4, []int{44, 40, 36, 32, 60, 56, 52, 48, 28, 24, 20, 16, 12, 8, 4, 0, 45, 41, 37, 33, 61, 57, 53, 49, 29, 25, 21, 17, 13, 9, 5, 1, 46, 42, 38, 34, 62, 58, 54, 50, 30, 26, 22, 18, 14, 10, 6, 2, 47, 43, 39, 35, 63, 59, 55, 51, 31, 27, 23, 19, 15, 11, 7, 3}),
			/* O5 */ NewPadGroupType(28, 2, []int{47, 45, 43, 41, 39, 37, 35, 33, 63, 61, 59, 57, 55, 53, 51, 49, 31, 29, 27, 25, 23, 21, 19, 17, 15, 13, 11, 9, 46, 44, 42, 40, 38, 36, 34, 32, 62, 60, 58, 56, 54, 52, 50, 48, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8}),
			/* O6 */ NewPadGroupType(28, 2, []int{39, 37, 35, 33, 63, 61, 59, 57, 55, 53, 51, 49, 31, 29, 27, 25, 23, 21, 19, 17, 15, 13, 11, 9, 7, 5, 3, 1, 38, 36, 34, 32, 62, 60, 58, 56, 54, 52, 50, 48, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8, 6, 4, 2, 0}),
			/* O7 */ NewPadGroupType(28, 2, []int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 48, 50, 52, 54, 56, 58, 60, 62, 32, 34, 36, 38, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 49, 51, 53, 55, 57, 59, 61, 63, 33, 35, 37, 39}),
			/* O8 */ NewPadGroupType(28, 2, []int{8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 48, 50, 52, 54, 56, 58, 60, 62, 32, 34, 36, 38, 40, 42, 44, 46, 9, 11, 13, 15, 17, 19, 21, 23, 25, 27, 29, 31, 49, 51, 53, 55, 57, 59, 61, 63, 33, 35, 37, 39, 41, 43, 45, 47}),
			/* P1 */ NewPadGroupType(16, 5, []int{47, 46, 41, 36, 63, 58, 53, 48, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 42, 37, 32, 59, 54, 49, 28, 24, 20, 16, 12, 8, 4, 0, -1, -1, 43, 38, 33, 60, 55, 50, 29, 25, 21, 17, 13, 9, 5, 1, -1, -1, 44, 39, 34, 61, 56, 51, 30, 26, 22, 18, 14, 10, 6, 2, -1, -1, 45, 40, 35, 62, 57, 52, 31, 27, 23, 19, 15, 11, 7, 3}),
			/* P2 */ NewPadGroupType(16, 5, []int{-1, -1, -1, -1, -1, -1, -1, -1, 27, 22, 17, 12, 7, 2, 1, 0, 44, 40, 36, 32, 60, 56, 52, 48, 28, 23, 18, 13, 8, 3, -1, -1, 45, 41, 37, 33, 61, 57, 53, 49, 29, 24, 19, 14, 9, 4, -1, -1, 46, 42, 38, 34, 62, 58, 54, 50, 30, 25, 20, 15, 10, 5, -1, -1, 47, 43, 39, 35, 63, 59, 55, 51, 31, 26, 21, 16, 11, 6, -1, -1}),
			/* P3 */ NewPadGroupType(14, 5, []int{3, 7, 11, 15, 20, 25, 30, 51, 56, 61, 34, 39, 43, 47, 2, 6, 10, 14, 19, 24, 29, 50, 55, 60, 33, 38, 42, 46, 1, 5, 9, 13, 18, 23, 28, 49, 54, 59, 32, 37, 41, 45, 0, 4, 8, 12, 17, 22, 27, 48, 53, 58, 63, 36, 40, 44, -1, -1, -1, -1, 16, 21, 26, 31, 52, 57, 62, 35, -1, -1}),
			/* P4 */ NewPadGroupType(14, 5, []int{3, 7, 12, 17, 22, 27, 48, 53, 58, 63, 35, 39, 43, 47, 2, 6, 11, 16, 21, 26, 31, 52, 57, 62, 34, 38, 42, 46, 1, 5, 10, 15, 20, 25, 30, 51, 56, 61, 33, 37, 41, 45, 0, 4, 9, 14, 19, 24, 29, 50, 55, 60, 32, 36, 40, 44, -1, -1, 8, 13, 18, 23, 28, 49, 54, 59, -1, -1, -1, -1}),
			/* Q1 */ NewPadGroupType(14, 5, []int{-1, -1, -1, -1, 59, 54, 49, 28, 23, 18, 13, 8, -1, -1, 44, 40, 36, 32, 60, 55, 50, 29, 24, 19, 14, 9, 4, 0, 45, 41, 37, 33, 61, 56, 51, 30, 25, 20, 15, 10, 5, 1, 46, 42, 38, 34, 62, 57, 52, 31, 26, 21, 16, 11, 6, 2, 47, 43, 39, 35, 63, 58, 53, 48, 27, 22, 17, 12, 7, 3}),
			/* Q2 */ NewPadGroupType(14, 5, []int{-1, -1, 35, 62, 57, 52, 31, 26, 21, 16, -1, -1, -1, -1, 44, 40, 36, 63, 58, 53, 48, 27, 22, 17, 12, 8, 4, 0, 45, 41, 37, 32, 59, 54, 49, 28, 23, 18, 13, 9, 5, 1, 46, 42, 38, 33, 60, 55, 50, 29, 24, 19, 14, 10, 6, 2, 47, 43, 39, 34, 61, 56, 51, 30, 25, 20, 15, 11, 7, 3}),
			/* Q3 */ NewPadGroupType(16, 5, []int{-1, -1, 6, 11, 16, 21, 26, 31, 51, 55, 59, 63, 35, 39, 43, 47, -1, -1, 5, 10, 15, 20, 25, 30, 50, 54, 58, 62, 34, 38, 42, 46, -1, -1, 4, 9, 14, 19, 24, 29, 49, 53, 57, 61, 33, 37, 41, 45, -1, -1, 3, 8, 13, 18, 23, 28, 48, 52, 56, 60, 32, 36, 40, 44, 0, 1, 2, 7, 12, 17, 22, 27, -1, -1, -1, -1, -1, -1, -1, -1}),
			/* Q4 */ NewPadGroupType(16, 5, []int{3, 7, 11, 15, 19, 23, 27, 31, 52, 57, 62, 35, 40, 45, -1, -1, 2, 6, 10, 14, 18, 22, 26, 30, 51, 56, 61, 34, 39, 44, -1, -1, 1, 5, 9, 13, 17, 21, 25, 29, 50, 55, 60, 33, 38, 43, -1, -1, 0, 4, 8, 12, 16, 20, 24, 28, 49, 54, 59, 32, 37, 42, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, 48, 53, 58, 63, 36, 41, 46, 47}),
		},
		[]padSize{
			{0.714285714, 2.5},
			{0.714285714, 5},
			{0.714285714, 10},
		})
}

func init() {
	mapping.RegisterCathodeSegmentationBuilder(17, createSegType17{})
}

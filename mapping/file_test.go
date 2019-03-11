package mapping_test

// All the tests that require an input file

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/mrrtf/pigiron/geo"
	"github.com/mrrtf/pigiron/mapping"
)

func assertInputForDetElemIsCorrect(m TestChannelList) {
	const NDES int = 156
	const NFEC int = 16828
	const NPADS int = 1064008

	if len(m.DetectionElements) != NDES {
		log.Fatalf("Expected %d detection elements, got %d", NDES, len(m.DetectionElements))
	}

	nfec := 0
	npads := 0
	for _, de := range m.DetectionElements {
		nfec += len(de.FECs)
		for _, manu := range de.FECs {
			n := len(manu.Channels)
			npads += n
		}
	}

	// sanity check before starting our tests
	if nfec != NFEC {
		log.Fatalf("Expected %d FEC, got %d", NFEC, nfec)
	}

	if npads != NPADS {
		log.Fatalf("Expected %d pads, got %d", NPADS, npads)
	}

}

func checkSegmentationCreation(t *testing.T, deid mapping.DEID) {
	for _, isBendingPlane := range []bool{true, false} {
		seg := mapping.NewCathodeSegmentation(deid, isBendingPlane)
		if seg == nil {
			t.Errorf("Could not create segmentation for de %d isbending %v", deid, isBendingPlane)
		}
	}
}

// check all fecs are present
// must be called after the checkSegmentationCreation
// method as we do not recheck here that we can
// actually create the segmentation, we assume we can.
func checkAllFECIDs(t *testing.T, de DetectionElement) {
	fecids := de.fecIDs()
	var dualSampas []mapping.DualSampaID
	for _, isBendingPlane := range []bool{true, false} {
		seg := mapping.NewCathodeSegmentation(de.ID, isBendingPlane)
		for dsi := 0; dsi < seg.NofDualSampas(); dsi++ {
			dsid, _ := seg.DualSampaID(dsi)
			dualSampas = append(dualSampas, dsid)
		}

	}
	if !containSameDualSampaIDs(dualSampas, fecids) {
		t.Errorf("DE %d does not have the right FEC ids : %v vs %v", de.ID, dualSampas, fecids)
	}
}

func channelFromOneDualSampa(cseg mapping.CathodeSegmentation, dsid mapping.DualSampaID) ChannelInfoSlice {
	var channels ChannelInfoSlice

	cseg.ForEachPadInDualSampa(dsid, func(padcid mapping.PadCID) {
		if cseg.PadDualSampaID(padcid) != dsid {
			log.Fatalf("actual %d != expected %d for padcid %d", cseg.PadDualSampaID(padcid), dsid, padcid)
		}
		channels = append(channels, ChannelInfo{dsid, cseg.PadDualSampaChannel(padcid)})
	})

	return channels
}

// check all channels are present
// must be called after the checkSegmentationCreation
// method as we do not recheck here that we can
// actually create the segmentation, we assume we can.
func checkAllChannels(t *testing.T, de DetectionElement) {
	deChannels := de.channelIDs()
	var channels ChannelInfoSlice
	for _, isBendingPlane := range []bool{true, false} {
		seg := mapping.NewCathodeSegmentation(de.ID, isBendingPlane)
		for dsi := 0; dsi < seg.NofDualSampas(); dsi++ {
			dsid, _ := seg.DualSampaID(dsi)
			tchannels := channelFromOneDualSampa(seg, dsid)
			expectedChannels := de.channels(dsid)
			channels = append(channels, tchannels...)
			if !containSameChannelElements(tchannels, expectedChannels) {
				t.Fatalf("Different channels for DE %d DUALSAMPA %d : %v vs %v",
					de.ID, dsid, tchannels, expectedChannels)
			}
		}

	}
	if !containSameChannelElements(channels, deChannels) {
		t.Errorf("DE %d does not have the right channels", de.ID)
	}

}

func TestDetectionElementChannels(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	path := filepath.Join("testdata", "test_channel_list.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("could not read test file")
	}

	detElems, err := UnmarshalTestChannelList(data)
	if err != nil {
		log.Fatal("could not decode test file")
	}

	// check our input makes sense
	assertInputForDetElemIsCorrect(detElems)

	// now for the real tests ...
	for _, de := range detElems.DetectionElements {
		t.Run(fmt.Sprintf("checkSegmentationCreation(%d)", de.ID), func(t *testing.T) {
			checkSegmentationCreation(t, de.ID)
		})
		t.Run(fmt.Sprintf("checkAllFECIDs(%d)", de.ID), func(t *testing.T) {
			checkAllFECIDs(t, de)
		})
		t.Run(fmt.Sprintf("checkAllChannels(%d)", de.ID), func(t *testing.T) {
			checkAllChannels(t, de)
		})
	}
}

func testOnePosition(cseg mapping.CathodeSegmentation, tp Testposition) error {

	padcid, err := cseg.FindPadByPosition(tp.X, tp.Y)

	if err != nil && err != mapping.ErrInvalidPadCID {
		return fmt.Errorf("Unexpected error:%s", err)
	}

	if cseg.IsValid(padcid) && tp.isOutside() {
		return fmt.Errorf("found a pad at position where there should not be one : %v", tp)
	}

	if !cseg.IsValid(padcid) && !tp.isOutside() {
		return fmt.Errorf("did not find a pad at position where there should be one : %v", tp)
	}

	if cseg.IsValid(padcid) && (!geo.EqualFloat(cseg.PadPositionX(padcid), tp.PX) ||
		!geo.EqualFloat(cseg.PadPositionY(padcid), tp.PY)) {

		buf := new(bytes.Buffer)
		PrintCathodePad(buf, cseg, padcid)
		s := buf.String()

		return fmt.Errorf("\nExpected %v\nGot %v", tp.String(), s)
	}
	return nil
}

// TestPositions reads in test position from an external
// json file and checks that the FindPadByPosition function
// agrees with the results in that file.
func TestPositions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	path := filepath.Join("testdata", "test_random_pos.json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	testfile, err := UnmarshalTestRandomPos(data)
	if err != nil {
		log.Fatal("could not decode test file")
	}

	var notok int
	for _, tp := range testfile.Testpositions {
		seg := mapping.NewCathodeSegmentation(tp.De, tp.isBendingPlane())
		err = testOnePosition(seg, tp)
		if err != nil {
			t.Log(err)
			notok++
		}
	}
}

// TestNeighbours reads in an external json file containing
// for each pad the list of its neighbours and checks that
// the GetNeighbours function agrees with the results in
// that file
func TestNeighbours(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	ntest := 0
	nfail := 0
	for _, deid := range []mapping.DEID{100, 300, 500, 501, 502, 503, 504, 600, 601, 602, 700, 701, 702, 703, 704, 705, 706, 902, 903, 904, 905} {
		testNeighboursOneDE(t, deid, &ntest, &nfail)
	}
	if nfail != 0 {
		t.Errorf("tested:%d failed:%d", ntest, nfail)
	}
}

type dsIDCh struct {
	id mapping.DualSampaID
	ch mapping.DualSampaChannelID
}

func jsonGetNeighbours(tnei TestNeighbourStruct, deid mapping.DEID, dsid mapping.DualSampaID, dsch mapping.DualSampaChannelID) []dsIDCh {
	var neighbours []dsIDCh
	for _, nei := range tnei.Neighbours {
		if nei.Deid != deid {
			continue
		}
		for _, ds := range nei.Ds {
			if ds.ID != dsid {
				continue
			}
			for _, channels := range ds.Channels {
				if channels.Ch != dsch {
					continue
				}
				for _, n := range channels.Nei {
					neighbours = append(neighbours, dsIDCh{n.Dsid, n.Dsch})
				}
			}
		}
	}
	return neighbours
}

func compareNeighbours(nref []dsIDCh, ni []int, nnei int, cseg mapping.CathodeSegmentation) error {
	if len(nref) != nnei {
		return fmt.Errorf("Want %d neighbours - Got %d", len(nref), nnei)
	}

	var n2 []mapping.PadCID
	// convert dsIDCh to padcids
	for _, dsch := range nref {
		padcid, err := cseg.FindPadByFEE(dsch.id, dsch.ch)
		if err != nil {
			return fmt.Errorf("Got an invalid pad for DS %4d CH %2d", dsch.id, dsch.ch)
		}
		n2 = append(n2, padcid)
	}

	var n []mapping.PadCID
	for _, i := range ni {
		n = append(n, mapping.PadCID(i))
	}

	sort.Slice(n2, func(i, j int) bool {
		return int(n2[i]) > int(n2[j])
	})

	sort.Slice(n[:nnei], func(i, j int) bool {
		return int(n[i]) > int(n[j])
	})

	for i := 0; i < nnei; i++ {
		if n[i] != n2[i] {
			return fmt.Errorf("Wanted PadCID %d - Got %d", n2[i], n[i])
		}
	}
	return nil
}

func testNeighboursOneDE(t *testing.T, deid mapping.DEID, ntest, nfail *int) {
	path := filepath.Join("testdata", "test_neighbours_list_"+strconv.Itoa(int(deid))+".json")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	tnei, err := UnmarshalTestNeighbours(data)
	if err != nil {
		log.Fatal("could not decode test file")
	}

	for _, bending := range []bool{true, false} {
		cseg := mapping.NewCathodeSegmentation(deid, bending)
		n := make([]int, 13)
		cseg.ForEachPad(func(padcid mapping.PadCID) {
			*ntest++
			dsid := cseg.PadDualSampaID(padcid)
			dsch := cseg.PadDualSampaChannel(padcid)
			nref := jsonGetNeighbours(tnei, deid, dsid, dsch)
			nnei := cseg.GetNeighbourIDs(padcid, n)
			if len(n) != 13 {
				log.Fatal("ZOB")
			}
			err := compareNeighbours(nref, n, nnei, cseg)
			if err != nil {
				t.Errorf("Problem for DE %4d DS %4d CH %2d : %s", deid, dsid, dsch, err.Error())
				t.Errorf("%v", nref)
				msg := ">"
				for _, p := range n {
					msg += fmt.Sprintf("(%v %v) ", cseg.PadDualSampaID(mapping.PadCID(p)), cseg.PadDualSampaChannel(mapping.PadCID(p)))
				}
				t.Errorf(msg)
				*nfail++
			}
		})
	}
}

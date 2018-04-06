package mapping_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"testing"

	"github.com/aphecetche/pigiron/mapping"
)

func assertInputForDetElemIsCorrect(m Mch) {
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

func checkSegmentationCreation(t *testing.T, deid int) {
	for _, isBendingPlane := range []bool{true, false} {
		seg := mapping.NewSegmentation(deid, isBendingPlane)
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
	var dualSampas []int
	for _, isBendingPlane := range []bool{true, false} {
		seg := mapping.NewSegmentation(de.ID, isBendingPlane)
		for dsi := 0; dsi < seg.NofDualSampas(); dsi++ {
			dsid, _ := seg.DualSampaID(dsi)
			dualSampas = append(dualSampas, dsid)
		}

	}
	if !containSameIntElements(dualSampas, fecids) {
		sort.Ints(dualSampas)
		sort.Ints(fecids)
		t.Errorf("DE %d does not have the right FEC ids : %v vs %v", de.ID, dualSampas, fecids)
	}
}

func channelFromOneDualSampa(seg mapping.Segmentation, dualSampaID int) ChannelInfoSlice {
	var channels ChannelInfoSlice

	seg.ForEachPadInDualSampa(dualSampaID, func(paduid int) {
		if seg.PadDualSampaID(paduid) != dualSampaID {
			log.Fatalf("actual %d != expected %d for paduid %d", seg.PadDualSampaID(paduid), dualSampaID, paduid)
		}
		channels = append(channels, ChannelInfo{dualSampaID, seg.PadDualSampaChannel(paduid)})
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
		seg := mapping.NewSegmentation(de.ID, isBendingPlane)
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
	path := filepath.Join("testdata", "test_de.json")
	fmt.Print(path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("could not read test file")
	}

	detElems, err := UnmarshalMch(data)
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

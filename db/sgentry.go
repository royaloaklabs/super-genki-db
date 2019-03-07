package db

import (
	"strings"

	"github.com/Xsixteen/super-genki-db/freq"
	"github.com/Xsixteen/super-genki-db/jmdict"
	"github.com/gojp/kana"
)

type SGEntry struct {
	Id        int
	Japanese  string
	Furigana  string
	English   string
	Romanji   string
	Frequency float64
}

func NewSGEntryFromJMDict(jme *jmdict.Entry) *SGEntry {
	entry := &SGEntry{}

	entry.Id = jme.EntSeq

	if len(jme.KEle) == 0 {
		// no kanji representation, default to first kana
		entry.Japanese = jme.Rele[0].Reb
	} else {
		// has kanji representation, default w/ kanji and kana
		entry.Japanese = jme.KEle[0].Keb
		entry.Furigana = jme.Rele[0].Reb
	}

	// join all Gloss from Sense and combine all Sense into one string
	var tempGloss []string
	var tempSense []string
	for _, sense := range jme.Sense {
		for _, gloss := range sense.Gloss {
			tempGloss = append(tempGloss, gloss.Value)
		}
		tempSense = append(tempSense, strings.Join(tempGloss, GlossDelimiter))
	}
	entry.English = strings.Join(tempSense, SenseDelimiter)

	entry.Romanji = kana.KanaToRomaji(jme.Rele[0].Reb)

	if frequency, ok := freq.DataTable[entry.Japanese]; ok {
		// frequency data found for Japanese entry (kanji or kana)
		entry.Frequency = frequency
	} else if frequency, ok := freq.DataTable[entry.Furigana]; ok {
		// frequency data found for furigana but not kanji
		// can be multiple (for example ない) be more precise later on
		entry.Frequency = frequency
	} else {
		// nothing set frequency to -1
		entry.Frequency = -1.0
	}
	return entry
}

package db

import (
	"strings"

	"github.com/gojp/kana"
	"github.com/royaloaklabs/super-genki-db/freq"
	"github.com/royaloaklabs/super-genki-db/jmdict"
)

const (
	Delimiter      = " "
	SenseDelimiter = "; "
	GlossDelimiter = ", "
)

type SGEntry struct {
	// base data for FTS4 table
	Id        int
	Japanese  string
	Furigana  string
	English   string
	Romaji    string
	Frequency float64

	// extra data for readings table
	KanjiAlt   []string
	ReadingAlt []string

	// extra data for definitions table
	Sense []Sense
}

type Sense struct {
	POS   string
	Gloss string
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

		// add alternative kanji readings (if any)
		if len(jme.KEle) > 1 {
			for i := 1; i < len(jme.KEle); i++ {
				entry.KanjiAlt = append(entry.KanjiAlt, jme.KEle[i].Keb)
			}
		}
	}

	// add alternative kanji readings (if any)
	if len(jme.Rele) > 1 {
		for i := 1; i < len(jme.Rele); i++ {
			entry.ReadingAlt = append(entry.ReadingAlt, jme.Rele[i].Reb)
		}
	}

	// join all Gloss from Sense and combine all Sense into one string
	var tempGloss []string
	var tempSense []string
	for _, sense := range jme.Sense {
		for _, gloss := range sense.Gloss {
			tempGloss = append(tempGloss, gloss.Value)
		}

		// used for definitions table
		entry.Sense = append(entry.Sense, Sense{
			POS:   strings.Join(sense.Pos, "; "),
			Gloss: strings.Join(tempGloss, ";;"),
		})

		tempSense = append(tempSense, strings.Join(tempGloss, GlossDelimiter))
		tempGloss = make([]string, 0)
	}
	entry.English = strings.Join(tempSense, SenseDelimiter)

	entry.Romaji = kana.KanaToRomaji(jme.Rele[0].Reb)

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

package keys

import "fmt"

var allMiscKeys []string

func MiscKeys() []string {
	if len(allMiscKeys) == 0 {

		var formats = []string{
			"%[1]vleading-%[2]v",
			"%[1]v/leading-%[2]v-path",
			"%[1]v/%[1]v/leading-multiple-%[2]v-path",
			"trailing-%[2]v%[1]v",
			"trailing-%[2]v-path/%[1]v",
			"trailing-multiple-%[2]v-path/%[1]v/%[1]v",
			"inter%[1]v[2]%v",
			"inter/%[1]v/%[2]v-path",
			"inter-multiple/%[1]v/%[1]v/%[2]v-path",
		}

		formatKeys := func(substr string, name string) []string {
			keys := make([]string, len(formats))
			for index, format := range formats {
				key := fmt.Sprintf(format, substr, name)
				keys[index] = key
			}
			return keys
		}

		var allKeySets = [][]string{
			formatKeys(".", "dot"),
			formatKeys("..", "double-dot"),
			formatKeys("...", "triple-dot"),

			formatKeys("/", "slash"),
			formatKeys("//", "double-slash"),
			formatKeys("///", "triple-slash"),

			formatKeys("\\", "backslash"),
			formatKeys("\\\\", "double-backslash"),
			formatKeys("\\\\\\", "triple-backslash"),

			formatKeys("Ĺąŧıņ", "latin-extended-a"),
			formatKeys("Ƚȃțȉƞ", "latin-extended-b"),
			formatKeys("Ɫⱥⱦin", "latin-extended-c"),
			formatKeys("ɪpɐ", "ipa-extensions"),
			formatKeys("ʰʸˀˈː", "spacing-modifiers"),
			formatKeys("̘̀̈̐", "combining-diacriticals"),

			formatKeys("ελληνικά", "greek"),
			formatKeys("ⲧⲙⲛ̄ⲧⲣⲙ̄ⲛ̄ⲕⲏⲙⲉ", "coptic"),
			formatKeys("Кириллица", "cyrillic"),
			formatKeys("ԀԈԐԘ", "cyrillic-supplement"),
			formatKeys("հայերէն", "armenian"),
			formatKeys("עִבְרִית", "hebrew"),
			formatKeys("العَرَبِيَّة", "arabic"),
			formatKeys("فارسی", "persian"),
			formatKeys("اُردُو", "urdu"),
			formatKeys("ܣܘܪܝܝܐ", "syriac"),
			formatKeys("ތާނަ", "thaana"),
			formatKeys("ߒߞߏ", "n'ko"),
			formatKeys("ࠀࠈࠐ࠘", "samaritan"),
			formatKeys("ࡀࡄࡈࡌ", "mandaic"),
			formatKeys("हिन्दी", "devanagari"),
			formatKeys("বাংলা", "bengali"),
			formatKeys("ਗੁਰਮੁਖੀ", "gurmukhi"),
			formatKeys("ગુજરાતી", "gujarati"),
			formatKeys("ଓଡ଼ିଆ", "odia"),
			formatKeys("தமிழ்", "tamil"),
			formatKeys("తెలుగు", "telugu"),
			formatKeys("ಕನ್ನಡ", "kannada"),
			formatKeys("മലയാളം", "malayalam"),
			formatKeys("සිංහල", "sinhala"),
			formatKeys("ภาษาไทย", "thai"),
			formatKeys("ພາສາລາວ", "lao"),
			formatKeys("བོད་སྐད།", "tibetan"),
			formatKeys("မြန်မာစာ", "myanmar"),
			formatKeys("მხედრული", "georgian"),
			formatKeys("ᄀᄄᄈᄌ", "hangul-jamo"),
			formatKeys("አማርኛ", "ethiopic"),
			formatKeys("ᎀᎄᎈᎌ", "ethiopic-supplement"),
			formatKeys("ᏣᎳᎩ", "cherokee"),
			formatKeys("ᓄᓇᕗᑦ", "ucas"),
			formatKeys(" ᚄᚈᚌ", "ogham"),
			formatKeys("ᚠᚤᚨᚬ", "runic"),
			formatKeys("ᜊᜌ᜔ᜊᜌᜒᜈ᜔", "tagalog"),
			formatKeys("ᜱᜨᜳᜨᜳᜢ", "hanunoo"),
			formatKeys("ᝊᝓᝑᝒ", "buhid"),
			formatKeys("ᝦᝪᝯ", "tagbanwa"),
			formatKeys("ខេមរភាសា", "khmer"),
			formatKeys("ᠮᠣᠩᠭᠣᠯ ᠬᠡᠯᠡ", "mongolian"),
			formatKeys("\u18b0\u18b8\u18a0\u18a8", "ucas-extended"),
			formatKeys("ᤕᤠᤰᤌᤢᤱ", "limbu"),
			formatKeys("ᥖᥭᥰᥘᥫᥴ", "tai-le"),

			formatKeys("汉语", "chinese-simplified"),
			formatKeys("漢語", "chinese-traditional"),
			formatKeys("한국어", "korean"),
			formatKeys("日本語", "japanese"),

			formatKeys("😀😈😨😸💛💣🤷👩‍🌾", "emoji"),
		}

		for _, ks := range allKeySets {
			allMiscKeys = append(allMiscKeys, ks...)
		}
	}
	return allMiscKeys
}

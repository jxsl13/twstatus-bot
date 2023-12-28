package model

import (
	"fmt"

	"github.com/jxsl13/twstatus-bot/utils"
)

var (

	// https://github.com/teeworlds/teeworlds/blob/a1911c8f7d8458fb4076ef8e7651e8ef5e91ab3e/datasrc/countryflags/index.json#L30
	// https://github.com/ddnet/ddnet/blob/e0052c3aec74490c328ef6205917016974d790d4/data/countryflags/index.txt#L32
	flags = map[int]Flag{
		737: {737, "ss", ":flag_ss:"},
		901: {901, "gb-eng", ":england:"},  // XEN - England
		902: {902, "gb-nir", ":flag_je:"},  // XNI - Northern Ireland
		903: {903, "gb-sct", ":scotland:"}, // XSC - Scotland
		904: {904, "gb-wls", ":wales:"},    // XWA - Wales
		905: {905, "eu", ":flag_eu:"},
		906: {906, "es-ct", ":flag_es:"},
		907: {907, "es-ga", ":flag_es:"},
		950: {950, "xbz", ":flag_es:"},       // XBZ - Balearic Islands
		951: {951, "xca", ":flag_es:"},       // XCA - Catalonia
		952: {952, "xes", ":flag_es:"},       // XES - Spain
		953: {953, "xga", ":flag_es:"},       // XGA - Galicia
		-1:  {-1, "default", ":flag_black:"}, // default flag
		20:  {20, "ad", ":flag_ad:"},
		784: {784, "ae", ":flag_ae:"},
		4:   {4, "af", ":flag_af:"},
		28:  {28, "ag", ":flag_ag:"},
		660: {660, "ai", ":flag_ai:"},
		8:   {8, "al", ":flag_al:"},
		51:  {51, "am", ":flag_am:"},
		24:  {24, "ao", ":flag_ao:"},
		32:  {32, "ar", ":flag_ar:"},
		16:  {16, "as", ":flag_as:"},
		40:  {40, "at", ":flag_at:"},
		36:  {36, "au", ":flag_au:"},
		533: {533, "aw", ":flag_aw:"},
		248: {248, "ax", ":flag_ax:"},
		31:  {31, "az", ":flag_az:"},
		70:  {70, "ba", ":flag_ba:"},
		52:  {52, "bb", ":flag_bb:"},
		50:  {50, "bd", ":flag_bd:"},
		56:  {56, "be", ":flag_be:"},
		854: {854, "bf", ":flag_bf:"},
		100: {100, "bg", ":flag_bg:"},
		48:  {48, "bh", ":flag_bh:"},
		108: {108, "bi", ":flag_bi:"},
		204: {204, "bj", ":flag_bj:"},
		652: {652, "bl", ":flag_bl:"},
		60:  {60, "bm", ":flag_bm:"},
		96:  {96, "bn", ":flag_bn:"},
		68:  {68, "bo", ":flag_bo:"},
		76:  {76, "br", ":flag_br:"},
		44:  {44, "bs", ":flag_bs:"},
		64:  {64, "bt", ":flag_bt:"},
		72:  {72, "bw", ":flag_bw:"},
		112: {112, "by", ":flag_by:"},
		84:  {84, "bz", ":flag_bz:"},
		124: {124, "ca", ":flag_ca:"},
		166: {166, "cc", ":flag_cc:"},
		180: {180, "cd", ":flag_cd:"},
		140: {140, "cf", ":flag_cf:"},
		178: {178, "cg", ":flag_cg:"},
		756: {756, "ch", ":flag_ch:"},
		384: {384, "ci", ":flag_ci:"},
		184: {184, "ck", ":flag_ck:"},
		152: {152, "cl", ":flag_cl:"},
		120: {120, "cm", ":flag_cm:"},
		156: {156, "cn", ":flag_cn:"},
		170: {170, "co", ":flag_co:"},
		188: {188, "cr", ":flag_cr:"},
		192: {192, "cu", ":flag_cu:"},
		132: {132, "cv", ":flag_cv:"},
		531: {531, "cw", ":flag_cw:"},
		162: {162, "cx", ":flag_cx:"},
		196: {196, "cy", ":flag_cy:"},
		203: {203, "cz", ":flag_cz:"},
		276: {276, "de", ":flag_de:"},
		262: {262, "dj", ":flag_dj:"},
		208: {208, "dk", ":flag_dk:"},
		212: {212, "dm", ":flag_dm:"},
		214: {214, "do", ":flag_do:"},
		12:  {12, "dz", ":flag_dz:"},
		218: {218, "ec", ":flag_ec:"},
		233: {233, "ee", ":flag_ee:"},
		818: {818, "eg", ":flag_eg:"},
		732: {732, "eh", ":flag_eh:"},
		232: {232, "er", ":flag_er:"},
		724: {724, "es", ":flag_es:"},
		231: {231, "et", ":flag_et:"},
		246: {246, "fi", ":flag_fi:"},
		242: {242, "fj", ":flag_fj:"},
		238: {238, "fk", ":flag_fk:"},
		583: {583, "fm", ":flag_fm:"},
		234: {234, "fo", ":flag_fo:"},
		250: {250, "fr", ":flag_fr:"},
		266: {266, "ga", ":flag_ga:"},
		826: {826, "gb", ":flag_gb:"},
		308: {308, "gd", ":flag_gd:"},
		268: {268, "ge", ":flag_ge:"},
		254: {254, "gf", ":flag_gf:"},
		831: {831, "gg", ":flag_gg:"},
		288: {288, "gh", ":flag_gh:"},
		292: {292, "gi", ":flag_gi:"},
		304: {304, "gl", ":flag_gl:"},
		270: {270, "gm", ":flag_gm:"},
		324: {324, "gn", ":flag_gn:"},
		312: {312, "gp", ":flag_gp:"},
		226: {226, "gq", ":flag_gq:"},
		300: {300, "gr", ":flag_gr:"},
		239: {239, "gs", ":flag_gs:"},
		320: {320, "gt", ":flag_gt:"},
		316: {316, "gu", ":flag_gu:"},
		624: {624, "gw", ":flag_gw:"},
		328: {328, "gy", ":flag_gy:"},
		344: {344, "hk", ":flag_hk:"},
		340: {340, "hn", ":flag_hn:"},
		191: {191, "hr", ":flag_hr:"},
		332: {332, "ht", ":flag_ht:"},
		348: {348, "hu", ":flag_hu:"},
		360: {360, "id", ":flag_id:"},
		372: {372, "ie", ":flag_ie:"},
		376: {376, "il", ":flag_il:"},
		833: {833, "im", ":flag_im:"},
		356: {356, "in", ":flag_in:"},
		86:  {86, "io", ":flag_io:"},
		368: {368, "iq", ":flag_iq:"},
		364: {364, "ir", ":flag_ir:"},
		352: {352, "is", ":flag_is:"},
		380: {380, "it", ":flag_it:"},
		832: {832, "je", ":flag_je:"},
		388: {388, "jm", ":flag_jm:"},
		400: {400, "jo", ":flag_jo:"},
		392: {392, "jp", ":flag_jp:"},
		404: {404, "ke", ":flag_ke:"},
		417: {417, "kg", ":flag_kg:"},
		116: {116, "kh", ":flag_kh:"},
		296: {296, "ki", ":flag_ki:"},
		174: {174, "km", ":flag_km:"},
		659: {659, "kn", ":flag_kn:"},
		408: {408, "kp", ":flag_kp:"},
		410: {410, "kr", ":flag_kr:"},
		414: {414, "kw", ":flag_kw:"},
		136: {136, "ky", ":flag_ky:"},
		398: {398, "kz", ":flag_kz:"},
		418: {418, "la", ":flag_la:"},
		422: {422, "lb", ":flag_lb:"},
		662: {662, "lc", ":flag_lc:"},
		438: {438, "li", ":flag_li:"},
		144: {144, "lk", ":flag_lk:"},
		430: {430, "lr", ":flag_lr:"},
		426: {426, "ls", ":flag_ls:"},
		440: {440, "lt", ":flag_lt:"},
		442: {442, "lu", ":flag_lu:"},
		428: {428, "lv", ":flag_lv:"},
		434: {434, "ly", ":flag_ly:"},
		504: {504, "ma", ":flag_ma:"},
		492: {492, "mc", ":flag_mc:"},
		498: {498, "md", ":flag_md:"},
		499: {499, "me", ":flag_me:"},
		663: {663, "mf", ":flag_mf:"},
		450: {450, "mg", ":flag_mg:"},
		584: {584, "mh", ":flag_mh:"},
		807: {807, "mk", ":flag_mk:"},
		466: {466, "ml", ":flag_ml:"},
		104: {104, "mm", ":flag_mm:"},
		496: {496, "mn", ":flag_mn:"},
		446: {446, "mo", ":flag_mo:"},
		580: {580, "mp", ":flag_mp:"},
		474: {474, "mq", ":flag_mq:"},
		478: {478, "mr", ":flag_mr:"},
		500: {500, "ms", ":flag_ms:"},
		470: {470, "mt", ":flag_mt:"},
		480: {480, "mu", ":flag_mu:"},
		462: {462, "mv", ":flag_mv:"},
		454: {454, "mw", ":flag_mw:"},
		484: {484, "mx", ":flag_mx:"},
		458: {458, "my", ":flag_my:"},
		508: {508, "mz", ":flag_mz:"},
		516: {516, "na", ":flag_na:"},
		540: {540, "nc", ":flag_nc:"},
		562: {562, "ne", ":flag_ne:"},
		574: {574, "nf", ":flag_nf:"},
		566: {566, "ng", ":flag_ng:"},
		558: {558, "ni", ":flag_ni:"},
		528: {528, "nl", ":flag_nl:"},
		578: {578, "no", ":flag_no:"},
		524: {524, "np", ":flag_np:"},
		520: {520, "nr", ":flag_nr:"},
		570: {570, "nu", ":flag_nu:"},
		554: {554, "nz", ":flag_nz:"},
		512: {512, "om", ":flag_om:"},
		591: {591, "pa", ":flag_pa:"},
		604: {604, "pe", ":flag_pe:"},
		258: {258, "pf", ":flag_pf:"},
		598: {598, "pg", ":flag_pg:"},
		608: {608, "ph", ":flag_ph:"},
		586: {586, "pk", ":flag_pk:"},
		616: {616, "pl", ":flag_pl:"},
		666: {666, "pm", ":flag_pm:"},
		612: {612, "pn", ":flag_pn:"},
		630: {630, "pr", ":flag_pr:"},
		275: {275, "ps", ":flag_ps:"},
		620: {620, "pt", ":flag_pt:"},
		585: {585, "pw", ":flag_pw:"},
		600: {600, "py", ":flag_py:"},
		634: {634, "qa", ":flag_qa:"},
		638: {638, "re", ":flag_re:"},
		642: {642, "ro", ":flag_ro:"},
		688: {688, "rs", ":flag_rs:"},
		643: {643, "ru", ":flag_ru:"},
		646: {646, "rw", ":flag_rw:"},
		682: {682, "sa", ":flag_sa:"},
		90:  {90, "sb", ":flag_sb:"},
		690: {690, "sc", ":flag_sc:"},
		736: {736, "sd", ":flag_sd:"},
		752: {752, "se", ":flag_se:"},
		702: {702, "sg", ":flag_sg:"},
		654: {654, "sh", ":flag_sh:"},
		705: {705, "si", ":flag_si:"},
		703: {703, "sk", ":flag_sk:"},
		694: {694, "sl", ":flag_sl:"},
		674: {674, "sm", ":flag_sm:"},
		686: {686, "sn", ":flag_sn:"},
		706: {706, "so", ":flag_so:"},
		740: {740, "sr", ":flag_sr:"},
		678: {678, "st", ":flag_st:"},
		222: {222, "sv", ":flag_sv:"},
		534: {534, "sx", ":flag_sx:"},
		760: {760, "sy", ":flag_sy:"},
		748: {748, "sz", ":flag_sz:"},
		796: {796, "tc", ":flag_tc:"},
		148: {148, "td", ":flag_td:"},
		260: {260, "tf", ":flag_tf:"},
		768: {768, "tg", ":flag_tg:"},
		764: {764, "th", ":flag_th:"},
		762: {762, "tj", ":flag_tj:"},
		772: {772, "tk", ":flag_tk:"},
		626: {626, "tl", ":flag_tl:"},
		795: {795, "tm", ":flag_tm:"},
		788: {788, "tn", ":flag_tn:"},
		776: {776, "to", ":flag_to:"},
		792: {792, "tr", ":flag_tr:"},
		780: {780, "tt", ":flag_tt:"},
		798: {798, "tv", ":flag_tv:"},
		158: {158, "tw", ":flag_tw:"},
		834: {834, "tz", ":flag_tz:"},
		804: {804, "ua", ":flag_ua:"},
		800: {800, "ug", ":flag_ug:"},
		840: {840, "us", ":flag_us:"},
		858: {858, "uy", ":flag_uy:"},
		860: {860, "uz", ":flag_uz:"},
		670: {670, "va", ":flag_va:"},
	}

	sortedFlagIds = utils.SortedMapKeys(flags)
)

type Flag struct {
	ID    int
	Abbr  string
	Emoji string
}

func (f Flag) String() string {
	return fmt.Sprintf("%s `%s`", f.Abbr, f.Emoji)
}

func Flags() []Flag {
	result := make([]Flag, 0, len(sortedFlagIds))
	for _, key := range sortedFlagIds {
		result = append(result, flags[key])
	}
	return result
}

func KnownFlag(id int) bool {
	_, found := flags[id]
	return found
}

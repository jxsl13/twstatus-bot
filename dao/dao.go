package dao

import (
	"context"
	"database/sql"
	"errors"

	"modernc.org/sqlite"
)

const (
	UniqueConstraintViolation = 1555
)

func IsUniqueConstraintErr(err error) bool {
	serr, ok := err.(*sqlite.Error)
	if ok {
		return serr.Code() == UniqueConstraintViolation
	}
	return false
}

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Conn interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type TxConn interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Commit() error
	Rollback() error
}

func InitDatabase(ctx context.Context, conn Conn) error {
	_, err := conn.ExecContext(ctx, `
	PRAGMA foreign_keys=on;

	CREATE TABLE IF NOT EXISTS guild (
		guild_id INTEGER NOT NULL PRIMARY KEY,
		description TEXT NOT NULL DEFAULT ""
	);

	CREATE TABLE IF NOT EXISTS channel (
		guild_id INTEGER NOT NULL,
		channel_id INTEGER NOT NULL,
		message_id INTEGER NOT NULL,
		running INTEGER NOT NULL DEFAULT 0,
		PRIMARY KEY (guild_id, channel_id),
    	FOREIGN KEY (guild_id)
    		REFERENCES guild (guild_id)
    		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS flag (
		flag_id INTEGER NOT NULL,
		channel_id INTEGER NOT NULL,
		abbr TEXT NOT NULL,
		symbol TEXT NOT NULL,
		PRIMARY KEY (flag_id, channel_id),
    	FOREIGN KEY (channel_id)
    		REFERENCES channel (channel_id)
    		ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS tw_server (
		channel_id INTEGER NOT NULL,
		address TEXT NOT NULL,
		PRIMARY KEY (channel_id, address),
    	FOREIGN KEY (channel_id)
    		REFERENCES channel (channel_id)
    		ON DELETE CASCADE
	);
`)
	return err
}

var flags = map[int][]string{
	737: {"ss", ":flag_ss:"},
	901: {"xen", ":england:"},        // XEN - England
	902: {"xni", ":flag_je:"},        // XNI - Northern Ireland
	903: {"xsc", ":scotland:"},       // XSC - Scotland
	904: {"xwa", ":wales:"},          // XWA - Wales
	950: {"xbz", ":flag_es:"},        // XBZ - Balearic Islands
	951: {"xca", ":flag_es:"},        // XCA - Catalonia
	952: {"xes", ":flag_es:"},        // XES - Spain
	953: {"xga", ":flag_xga:"},       // XGA - Galicia
	-1:  {"default", ":flag_black:"}, // default flag
	20:  {"ad", ":flag_ad:"},
	784: {"ae", ":flag_ae:"},
	4:   {"af", ":flag_af:"},
	28:  {"ag", ":flag_ag:"},
	660: {"ai", ":flag_ai:"},
	8:   {"al", ":flag_al:"},
	51:  {"am", ":flag_am:"},
	24:  {"ao", ":flag_ao:"},
	32:  {"ar", ":flag_ar:"},
	16:  {"as", ":flag_as:"},
	40:  {"at", ":flag_at:"},
	36:  {"au", ":flag_au:"},
	533: {"aw", ":flag_aw:"},
	248: {"ax", ":flag_ax:"},
	31:  {"az", ":flag_az:"},
	70:  {"ba", ":flag_ba:"},
	52:  {"bb", ":flag_bb:"},
	50:  {"bd", ":flag_bd:"},
	56:  {"be", ":flag_be:"},
	854: {"bf", ":flag_bf:"},
	100: {"bg", ":flag_bg:"},
	48:  {"bh", ":flag_bh:"},
	108: {"bi", ":flag_bi:"},
	204: {"bj", ":flag_bj:"},
	652: {"bl", ":flag_bl:"},
	60:  {"bm", ":flag_bm:"},
	96:  {"bn", ":flag_bn:"},
	68:  {"bo", ":flag_bo:"},
	76:  {"br", ":flag_br:"},
	44:  {"bs", ":flag_bs:"},
	64:  {"bt", ":flag_bt:"},
	72:  {"bw", ":flag_bw:"},
	112: {"by", ":flag_by:"},
	84:  {"bz", ":flag_bz:"},
	124: {"ca", ":flag_ca:"},
	166: {"cc", ":flag_cc:"},
	180: {"cd", ":flag_cd:"},
	140: {"cf", ":flag_cf:"},
	178: {"cg", ":flag_cg:"},
	756: {"ch", ":flag_ch:"},
	384: {"ci", ":flag_ci:"},
	184: {"ck", ":flag_ck:"},
	152: {"cl", ":flag_cl:"},
	120: {"cm", ":flag_cm:"},
	156: {"cn", ":flag_cn:"},
	170: {"co", ":flag_co:"},
	188: {"cr", ":flag_cr:"},
	192: {"cu", ":flag_cu:"},
	132: {"cv", ":flag_cv:"},
	531: {"cw", ":flag_cw:"},
	162: {"cx", ":flag_cx:"},
	196: {"cy", ":flag_cy:"},
	203: {"cz", ":flag_cz:"},
	276: {"de", ":flag_de:"},
	262: {"dj", ":flag_dj:"},
	208: {"dk", ":flag_dk:"},
	212: {"dm", ":flag_dm:"},
	214: {"do", ":flag_do:"},
	12:  {"dz", ":flag_dz:"},
	218: {"ec", ":flag_ec:"},
	233: {"ee", ":flag_ee:"},
	818: {"eg", ":flag_eg:"},
	732: {"eh", ":flag_eh:"},
	232: {"er", ":flag_er:"},
	724: {"es", ":flag_es:"},
	231: {"et", ":flag_et:"},
	246: {"fi", ":flag_fi:"},
	242: {"fj", ":flag_fj:"},
	238: {"fk", ":flag_fk:"},
	583: {"fm", ":flag_fm:"},
	234: {"fo", ":flag_fo:"},
	250: {"fr", ":flag_fr:"},
	266: {"ga", ":flag_ga:"},
	826: {"gb", ":flag_gb:"},
	308: {"gd", ":flag_gd:"},
	268: {"ge", ":flag_ge:"},
	254: {"gf", ":flag_gf:"},
	831: {"gg", ":flag_gg:"},
	288: {"gh", ":flag_gh:"},
	292: {"gi", ":flag_gi:"},
	304: {"gl", ":flag_gl:"},
	270: {"gm", ":flag_gm:"},
	324: {"gn", ":flag_gn:"},
	312: {"gp", ":flag_gp:"},
	226: {"gq", ":flag_gq:"},
	300: {"gr", ":flag_gr:"},
	239: {"gs", ":flag_gs:"},
	320: {"gt", ":flag_gt:"},
	316: {"gu", ":flag_gu:"},
	624: {"gw", ":flag_gw:"},
	328: {"gy", ":flag_gy:"},
	344: {"hk", ":flag_hk:"},
	340: {"hn", ":flag_hn:"},
	191: {"hr", ":flag_hr:"},
	332: {"ht", ":flag_ht:"},
	348: {"hu", ":flag_hu:"},
	360: {"id", ":flag_id:"},
	372: {"ie", ":flag_ie:"},
	376: {"il", ":flag_il:"},
	833: {"im", ":flag_im:"},
	356: {"in", ":flag_in:"},
	86:  {"io", ":flag_io:"},
	368: {"iq", ":flag_iq:"},
	364: {"ir", ":flag_ir:"},
	352: {"is", ":flag_is:"},
	380: {"it", ":flag_it:"},
	832: {"je", ":flag_je:"},
	388: {"jm", ":flag_jm:"},
	400: {"jo", ":flag_jo:"},
	392: {"jp", ":flag_jp:"},
	404: {"ke", ":flag_ke:"},
	417: {"kg", ":flag_kg:"},
	116: {"kh", ":flag_kh:"},
	296: {"ki", ":flag_ki:"},
	174: {"km", ":flag_km:"},
	659: {"kn", ":flag_kn:"},
	408: {"kp", ":flag_kp:"},
	410: {"kr", ":flag_kr:"},
	414: {"kw", ":flag_kw:"},
	136: {"ky", ":flag_ky:"},
	398: {"kz", ":flag_kz:"},
	418: {"la", ":flag_la:"},
	422: {"lb", ":flag_lb:"},
	662: {"lc", ":flag_lc:"},
	438: {"li", ":flag_li:"},
	144: {"lk", ":flag_lk:"},
	430: {"lr", ":flag_lr:"},
	426: {"ls", ":flag_ls:"},
	440: {"lt", ":flag_lt:"},
	442: {"lu", ":flag_lu:"},
	428: {"lv", ":flag_lv:"},
	434: {"ly", ":flag_ly:"},
	504: {"ma", ":flag_ma:"},
	492: {"mc", ":flag_mc:"},
	498: {"md", ":flag_md:"},
	499: {"me", ":flag_me:"},
	663: {"mf", ":flag_mf:"},
	450: {"mg", ":flag_mg:"},
	584: {"mh", ":flag_mh:"},
	807: {"mk", ":flag_mk:"},
	466: {"ml", ":flag_ml:"},
	104: {"mm", ":flag_mm:"},
	496: {"mn", ":flag_mn:"},
	446: {"mo", ":flag_mo:"},
	580: {"mp", ":flag_mp:"},
	474: {"mq", ":flag_mq:"},
	478: {"mr", ":flag_mr:"},
	500: {"ms", ":flag_ms:"},
	470: {"mt", ":flag_mt:"},
	480: {"mu", ":flag_mu:"},
	462: {"mv", ":flag_mv:"},
	454: {"mw", ":flag_mw:"},
	484: {"mx", ":flag_mx:"},
	458: {"my", ":flag_my:"},
	508: {"mz", ":flag_mz:"},
	516: {"na", ":flag_na:"},
	540: {"nc", ":flag_nc:"},
	562: {"ne", ":flag_ne:"},
	574: {"nf", ":flag_nf:"},
	566: {"ng", ":flag_ng:"},
	558: {"ni", ":flag_ni:"},
	528: {"nl", ":flag_nl:"},
	578: {"no", ":flag_no:"},
	524: {"np", ":flag_np:"},
	520: {"nr", ":flag_nr:"},
	570: {"nu", ":flag_nu:"},
	554: {"nz", ":flag_nz:"},
	512: {"om", ":flag_om:"},
	591: {"pa", ":flag_pa:"},
	604: {"pe", ":flag_pe:"},
	258: {"pf", ":flag_pf:"},
	598: {"pg", ":flag_pg:"},
	608: {"ph", ":flag_ph:"},
	586: {"pk", ":flag_pk:"},
	616: {"pl", ":flag_pl:"},
	666: {"pm", ":flag_pm:"},
	612: {"pn", ":flag_pn:"},
	630: {"pr", ":flag_pr:"},
	275: {"ps", ":flag_ps:"},
	620: {"pt", ":flag_pt:"},
	585: {"pw", ":flag_pw:"},
	600: {"py", ":flag_py:"},
	634: {"qa", ":flag_qa:"},
	638: {"re", ":flag_re:"},
	642: {"ro", ":flag_ro:"},
	688: {"rs", ":flag_rs:"},
	643: {"ru", ":flag_ru:"},
	646: {"rw", ":flag_rw:"},
	682: {"sa", ":flag_sa:"},
	90:  {"sb", ":flag_sb:"},
	690: {"sc", ":flag_sc:"},
	736: {"sd", ":flag_sd:"},
	752: {"se", ":flag_se:"},
	702: {"sg", ":flag_sg:"},
	654: {"sh", ":flag_sh:"},
	705: {"si", ":flag_si:"},
	703: {"sk", ":flag_sk:"},
	694: {"sl", ":flag_sl:"},
	674: {"sm", ":flag_sm:"},
	686: {"sn", ":flag_sn:"},
	706: {"so", ":flag_so:"},
	740: {"sr", ":flag_sr:"},
	678: {"st", ":flag_st:"},
	222: {"sv", ":flag_sv:"},
	534: {"sx", ":flag_sx:"},
	760: {"sy", ":flag_sy:"},
	748: {"sz", ":flag_sz:"},
	796: {"tc", ":flag_tc:"},
	148: {"td", ":flag_td:"},
	260: {"tf", ":flag_tf:"},
	768: {"tg", ":flag_tg:"},
	764: {"th", ":flag_th:"},
	762: {"tj", ":flag_tj:"},
	772: {"tk", ":flag_tk:"},
	626: {"tl", ":flag_tl:"},
	795: {"tm", ":flag_tm:"},
	788: {"tn", ":flag_tn:"},
	776: {"to", ":flag_to:"},
	792: {"tr", ":flag_tr:"},
	780: {"tt", ":flag_tt:"},
	798: {"tv", ":flag_tv:"},
	158: {"tw", ":flag_tw:"},
	834: {"tz", ":flag_tz:"},
	804: {"ua", ":flag_ua:"},
	800: {"ug", ":flag_ug:"},
	840: {"us", ":flag_us:"},
	858: {"uy", ":flag_uy:"},
	860: {"uz", ":flag_uz:"},
	336: {"va", ":flag_va:"},
	670: {"vc", ":flag_vc:"},
	862: {"ve", ":flag_ve:"},
	92:  {"vg", ":flag_vg:"},
	850: {"vi", ":flag_vi:"},
	704: {"vn", ":flag_vn:"},
	548: {"vu", ":flag_vu:"},
	876: {"wf", ":flag_wf:"},
	882: {"ws", ":flag_ws:"},
	887: {"ye", ":flag_ye:"},
	710: {"za", ":flag_za:"},
	894: {"zm", ":flag_zm:"},
	716: {"zw", ":flag_zw:"},
	10:  {"aq", ":flag_aq:"},
	535: {"bq", ":flag_bq:"},
	74:  {"bv", ":flag_bv:"},
	334: {"hm", ":flag_hm:"},
	744: {"sj", ":flag_sj:"},
	581: {"um", ":flag_um:"},
	175: {"yt", ":flag_yt:"},
}

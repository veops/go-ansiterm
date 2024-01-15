package _const

const (
	SP    = " "
	NUL   = "\x00"
	BEL   = "\x07"
	BS    = "\x08"
	HT    = "\x09"
	LF    = "\n"
	VT    = "\x0b"
	FF    = "\x0c"
	CR    = "\r"
	SO    = "\x0e"
	SI    = "\x0f"
	CAN   = "\x18"
	SUB   = "\x1a"
	ESC   = "\x1b"
	DEL   = "\x7f"
	CSIC0 = ESC + "["
	CSIC1 = "\x9b"
	CSI   = CSIC0
	STC0  = ESC + "\\"
	STC1  = "\x9c"
	ST    = STC0
	OSCC0 = ESC + "]"
	OSCC1 = "\x9d"
	OSC   = OSCC0
)

const (
	RegSP    = ` `
	RegNUL   = `\x00`
	RegBEL   = `\x07`
	RegBS    = `\x08`
	RegHT    = `\x09`
	RegLF    = `\n`
	RegVT    = `\x0b`
	RegFF    = `\x0c`
	RegCR    = `\r`
	RegSO    = `\x0e`
	RegSI    = `\x0f`
	RegCAN   = `\x18`
	RegSUB   = `\x1a`
	RegESC   = `\x1b`
	RegDEL   = `\x7f`
	RegCSIC0 = RegESC + `[`
	RegCSIC1 = `\x9b`
	RegCSI   = RegCSIC0
	RegSTC0  = RegESC + `\\`
	RegSTC1  = `\x9c`
	RegST    = RegSTC0
	RegOSCC0 = RegESC + `]`
	RegOSCC1 = `\x9d`
	RegOSC   = RegOSCC0
)

#include "textflag.h"

// func load(p unsafe.Pointer) i8x32
// Requires: AVX
TEXT ·load(SB), NOSPLIT, $0-38
	MOVQ p_base+0(FP), AX
	VMOVUPS (AX), Y0
	VMOVUPS Y0, ret+8(FP)
    //VZEROUPPER
    RET

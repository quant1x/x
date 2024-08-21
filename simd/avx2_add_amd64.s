#include "textflag.h"

// func AddFloat32x4(a, b Float32x4) Float32x4
//TEXT ·AddFloat32x4(SB), $0-48
//	MOVUPD b+0(FP), X0 // b => X0
//	MOVUPD a+16(FP), X1 // a => X1
//	ADDPS X1, X0
//	MOVUPD X0, ret+32(FP) // X0 => ret
//	RET

// func avx2_f32x8_add(x, y f32x8) f32x8
// Requires: AVX
TEXT ·avx2_f32x8_add(SB), $0-96
	VMOVUPS x+0(FP), Y0 // x => Y0
	VMOVUPS y+32(FP), Y1 // y => Y1
	VADDPS  Y1, Y0, Y0
	VMOVUPS Y0, ret+64(FP)
	VZEROUPPER
	RET

// func _mm256_add_ps(x, y, z []float32)
// Requires: AVX
TEXT ·_mm256_add_ps1(SB), $0-24
    MOVQ x+0(FP), DI
    MOVQ y+8(FP), SI
    MOVQ z+16(FP), DX

	VMOVUPS (DI), Y0 // x => Y0
	VMOVUPS (SI), Y1 // y => Y1
	VADDPS  Y1, Y0, Y0
	VMOVUPS Y0, (DX)
	VZEROUPPER
	RET

// func _mm256_add_ps(a []float32, b []float32, c []float32) int
// Requires: AVX
TEXT ·_mm256_add_ps(SB), NOSPLIT, $0-80
	MOVQ a_base+0(FP), AX
	MOVQ b_base+24(FP), CX
	MOVQ c_base+48(FP), DX
	MOVQ a_len+8(FP), BX

loop:
	CMPQ    BX, $0x00000008
	JL      done
	VMOVUPS (AX), Y0
	VMOVUPS (CX), Y1
	VADDPS  Y1, Y0, Y0
	VMOVUPS Y0, (DX)
	ADDQ    $0x00000020, AX
	ADDQ    $0x00000020, CX
	ADDQ    $0x00000020, DX
	SUBQ    $0x00000008, BX
	JMP     loop

done:
    MOVQ BX, ret+72(FP) // BX => ret
    VZEROUPPER
	RET

// func Add_AVX2_F32(x []float32, y []float32)
// Requires: AVX
TEXT ·Add_AVX2_F32(SB), NOSPLIT, $0-48
	MOVQ  x_base+0(FP), DI
	MOVQ  y_base+24(FP), SI
	MOVQ  x_len+8(FP), DX
	TESTQ DX, DX
	JE    LBB1_7
	CMPQ  DX, $0x20
	JAE   LBB1_3
	XORL  AX, AX
	JMP   LBB1_6

LBB1_3:
	MOVQ DX, AX
	ANDQ $-32, AX
	XORL CX, CX

LBB1_4:
	VMOVUPS (DI)(CX*4), Y0
	VMOVUPS 32(DI)(CX*4), Y1
	VMOVUPS 64(DI)(CX*4), Y2
	VMOVUPS 96(DI)(CX*4), Y3
	VADDPS  (SI)(CX*4), Y0, Y0
	VADDPS  32(SI)(CX*4), Y1, Y1
	VADDPS  64(SI)(CX*4), Y2, Y2
	VADDPS  96(SI)(CX*4), Y3, Y3
	VMOVUPS Y0, (DI)(CX*4)
	VMOVUPS Y1, 32(DI)(CX*4)
	VMOVUPS Y2, 64(DI)(CX*4)
	VMOVUPS Y3, 96(DI)(CX*4)
	ADDQ    $0x20, CX
	CMPQ    AX, CX
	JNE     LBB1_4
	CMPQ    AX, DX
	JE      LBB1_7

LBB1_6:
	VMOVSS (DI)(AX*4), X0
	VADDSS (SI)(AX*4), X0, X0
	VMOVSS X0, (DI)(AX*4)
	ADDQ   $0x01, AX
	CMPQ   DX, AX
	JNE    LBB1_6

LBB1_7:
	VZEROUPPER
	RET

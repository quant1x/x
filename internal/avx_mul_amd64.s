//+build !noasm !appengine
// AUTO-GENERATED BY C2GOASM -- DO NOT EDIT

TEXT ·_avx2_mm256_float32_add(SB), $0-32

    MOVQ a+0(FP), DI
    MOVQ b+8(FP), SI
    MOVQ c+16(FP), DX
    MOVQ n+24(FP), CX

    LONG $0x07418d48             // lea    rax, [rcx + 7]
    WORD $0x8548; BYTE $0xc9     // test    rcx, rcx
    LONG $0xc1490f48             // cmovns    rax, rcx
    WORD $0x8949; BYTE $0xc0     // mov    r8, rax
    LONG $0x03e8c149             // shr    r8, 3
    LONG $0xf8e08348             // and    rax, -8
    WORD $0x2948; BYTE $0xc1     // sub    rcx, rax
    WORD $0x8545; BYTE $0xc0     // test    r8d, r8d
	JLE LBB0_6
    WORD $0x8944; BYTE $0xc0     // mov    eax, r8d
    WORD $0xe083; BYTE $0x03     // and    eax, 3
    LONG $0x04f88341             // cmp    r8d, 4
	JB LBB0_4
    LONG $0xfce08141; WORD $0xffff; BYTE $0x7f // and    r8d, 2147483644
LBB0_3:
    QUAD $0x000000008610fcc5     // vmovups    ymm0, yword [rsi]
    QUAD $0x000000008758fcc5     // vaddps    ymm0, ymm0, yword [rdi]
    QUAD $0x000000008211fcc5     // vmovups    yword [rdx], ymm0
    QUAD $0x000000008610fcc5     // vmovups    ymm0, yword [rsi + 32]
    QUAD $0x000000008758fcc5     // vaddps    ymm0, ymm0, yword [rdi + 32]
    QUAD $0x000000008211fcc5     // vmovups    yword [rdx + 32], ymm0
    QUAD $0x000000008610fcc5     // vmovups    ymm0, yword [rsi + 64]
    QUAD $0x000000008758fcc5     // vaddps    ymm0, ymm0, yword [rdi + 64]
    QUAD $0x000000008211fcc5     // vmovups    yword [rdx + 64], ymm0
    QUAD $0x000000008610fcc5     // vmovups    ymm0, yword [rsi + 96]
    QUAD $0x000000008758fcc5     // vaddps    ymm0, ymm0, yword [rdi + 96]
    QUAD $0x000000008211fcc5     // vmovups    yword [rdx + 96], ymm0
    LONG $0x80ef8348             // sub    rdi, -128
    LONG $0x80ee8348             // sub    rsi, -128
    LONG $0x80ea8348             // sub    rdx, -128
    LONG $0xfcc08341             // add    r8d, -4
	JNE LBB0_3
LBB0_4:
    WORD $0xc085                 // test    eax, eax
	JE LBB0_6
LBB0_5:
    QUAD $0x000000008610fcc5     // vmovups    ymm0, yword [rsi]
    QUAD $0x000000008758fcc5     // vaddps    ymm0, ymm0, yword [rdi]
    QUAD $0x000000008211fcc5     // vmovups    yword [rdx], ymm0
    LONG $0x20c78348             // add    rdi, 32
    LONG $0x20c68348             // add    rsi, 32
    LONG $0x20c28348             // add    rdx, 32
    WORD $0xc8ff                 // dec    eax
	JNE LBB0_5
LBB0_6:
    WORD $0xc985                 // test    ecx, ecx
	JLE LBB0_19
    WORD $0xc889                 // mov    eax, ecx
    LONG $0x20f88348             // cmp    rax, 32
	JAE LBB0_9
    WORD $0x3145; BYTE $0xc0     // xor    r8d, r8d
	JMP LBB0_14
LBB0_9:
    WORD $0x8949; BYTE $0xd1     // mov    r9, rdx
    WORD $0x2949; BYTE $0xf9     // sub    r9, rdi
    WORD $0x3145; BYTE $0xc0     // xor    r8d, r8d
    LONG $0x80f98149; WORD $0x0000; BYTE $0x00 // cmp    r9, 128
	JB LBB0_14
    WORD $0x8949; BYTE $0xd1     // mov    r9, rdx
    WORD $0x2949; BYTE $0xf1     // sub    r9, rsi
    LONG $0x80f98149; WORD $0x0000; BYTE $0x00 // cmp    r9, 128
	JB LBB0_14
    WORD $0x8941; BYTE $0xc9     // mov    r9d, ecx
    LONG $0x1fe18341             // and    r9d, 31
    WORD $0x8949; BYTE $0xc0     // mov    r8, rax
    WORD $0x294d; BYTE $0xc8     // sub    r8, r9
    WORD $0x3145; BYTE $0xd2     // xor    r10d, r10d
LBB0_12:
    QUAD $0x00009684107ca1c4; WORD $0x0000 // vmovups    ymm0, yword [rsi + 4*r10]
    QUAD $0x0000968c107ca1c4; WORD $0x0000 // vmovups    ymm1, yword [rsi + 4*r10 + 32]
    QUAD $0x00009694107ca1c4; WORD $0x0000 // vmovups    ymm2, yword [rsi + 4*r10 + 64]
    QUAD $0x0000969c107ca1c4; WORD $0x0000 // vmovups    ymm3, yword [rsi + 4*r10 + 96]
    QUAD $0x00009784587ca1c4; WORD $0x0000 // vaddps    ymm0, ymm0, yword [rdi + 4*r10]
    QUAD $0x0000978c5874a1c4; WORD $0x0000 // vaddps    ymm1, ymm1, yword [rdi + 4*r10 + 32]
    QUAD $0x00009794586ca1c4; WORD $0x0000 // vaddps    ymm2, ymm2, yword [rdi + 4*r10 + 64]
    QUAD $0x0000979c5864a1c4; WORD $0x0000 // vaddps    ymm3, ymm3, yword [rdi + 4*r10 + 96]
    QUAD $0x00009284117ca1c4; WORD $0x0000 // vmovups    yword [rdx + 4*r10], ymm0
    QUAD $0x0000928c117ca1c4; WORD $0x0000 // vmovups    yword [rdx + 4*r10 + 32], ymm1
    QUAD $0x00009294117ca1c4; WORD $0x0000 // vmovups    yword [rdx + 4*r10 + 64], ymm2
    QUAD $0x0000929c117ca1c4; WORD $0x0000 // vmovups    yword [rdx + 4*r10 + 96], ymm3
    LONG $0x20c28349             // add    r10, 32
    WORD $0x394d; BYTE $0xd0     // cmp    r8, r10
	JNE LBB0_12
    WORD $0x854d; BYTE $0xc9     // test    r9, r9
	JE LBB0_19
LBB0_14:
    WORD $0x2944; BYTE $0xc1     // sub    ecx, r8d
    WORD $0x894d; BYTE $0xc1     // mov    r9, r8
    WORD $0xe183; BYTE $0x03     // and    ecx, 3
	JE LBB0_17
    WORD $0x894d; BYTE $0xc1     // mov    r9, r8
LBB0_16:
    LONG $0x107aa1c4; WORD $0x8e44; BYTE $0x04 // vmovss    xmm0, dword [rsi + 4*r9]
    LONG $0x587aa1c4; WORD $0x8f44; BYTE $0x04 // vaddss    xmm0, xmm0, dword [rdi + 4*r9]
    LONG $0x117aa1c4; WORD $0x8a44; BYTE $0x04 // vmovss    dword [rdx + 4*r9], xmm0
    WORD $0xff49; BYTE $0xc1     // inc    r9
    WORD $0xff48; BYTE $0xc9     // dec    rcx
	JNE LBB0_16
LBB0_17:
    WORD $0x2949; BYTE $0xc0     // sub    r8, rax
    LONG $0xfcf88349             // cmp    r8, -4
	JA LBB0_19
LBB0_18:
    LONG $0x107aa1c4; WORD $0x8e44; BYTE $0x04 // vmovss    xmm0, dword [rsi + 4*r9]
    LONG $0x587aa1c4; WORD $0x8f44; BYTE $0x04 // vaddss    xmm0, xmm0, dword [rdi + 4*r9]
    LONG $0x117aa1c4; WORD $0x8a44; BYTE $0x04 // vmovss    dword [rdx + 4*r9], xmm0
    LONG $0x107aa1c4; WORD $0x8e44; BYTE $0x08 // vmovss    xmm0, dword [rsi + 4*r9 + 4]
    LONG $0x587aa1c4; WORD $0x8f44; BYTE $0x08 // vaddss    xmm0, xmm0, dword [rdi + 4*r9 + 4]
    LONG $0x117aa1c4; WORD $0x8a44; BYTE $0x08 // vmovss    dword [rdx + 4*r9 + 4], xmm0
    LONG $0x107aa1c4; WORD $0x8e44; BYTE $0x0c // vmovss    xmm0, dword [rsi + 4*r9 + 8]
    LONG $0x587aa1c4; WORD $0x8f44; BYTE $0x0c // vaddss    xmm0, xmm0, dword [rdi + 4*r9 + 8]
    LONG $0x117aa1c4; WORD $0x8a44; BYTE $0x0c // vmovss    dword [rdx + 4*r9 + 8], xmm0
    LONG $0x107aa1c4; WORD $0x8e44; BYTE $0x10 // vmovss    xmm0, dword [rsi + 4*r9 + 12]
    LONG $0x587aa1c4; WORD $0x8f44; BYTE $0x10 // vaddss    xmm0, xmm0, dword [rdi + 4*r9 + 12]
    LONG $0x117aa1c4; WORD $0x8a44; BYTE $0x10 // vmovss    dword [rdx + 4*r9 + 12], xmm0
    LONG $0x04c18349             // add    r9, 4
    WORD $0x394c; BYTE $0xc8     // cmp    rax, r9
	JNE LBB0_18
LBB0_19:
    VZEROUPPER
    RET

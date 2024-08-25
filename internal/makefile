SRC_AVX2 = src/avx_mul.c
OBJS_AVX2 = $(SRC_AVX2:.c=.asm)
SRC_AVX512 = src/floats_avx512.c
OBJS_AVX512 = $(SRC_AVX512:.c=.asm)

#SOURCES = src/floats_avx.c src/floats_avx512.c
SOURCES = src/avx_mul.c
OBJS_ASM = $(SOURCES:.c=.asm)
OBJS_ASM1 = $(notdir $(OBJS_ASM))
OBJS_PLAN9 = $(OBJS_ASM1:.asm=_amd64.s)
#X86_FLAGS = -S -O3 -DENABLE_AVX2 -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -fno-builtin -ffast-math
X86_FLAGS := -S -O3 -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -fno-builtin -ffast-math
X86_AVX2_FLAGS = $(X86_FLAGS) -DENABLE_AVX2 -mavx -mfma
X86_AVX512_FLAGS = $(X86_FLAGS) -DENABLE_AVX512 -mavx -mavx512f -mavx512dq

.PHONY : all
all: asm plan9

asm: $(OBJS_AVX2) #$(OBJS_AVX512)
	@echo "asm..."

plan9:$(OBJS_PLAN9)
	@echo "plan9..."
	@echo $(OBJS_PLAN9)
	@echo "success"

$(OBJS_AVX2):%.asm : %.c
	@echo "clang avx2..."
	clang $(X86_AVX2_FLAGS) $< -o $@

#$(OBJS_AVX512):%.asm : %.c
#	@echo "clang avx512..."
#	clang $(X86_AVX512_FLAGS) $< -o $@

$(OBJS_PLAN9):$(OBJS_ASM)
	@echo "c2goasm..."
	asmfile=$(patsubst %_amd64.s,src/%.asm,$@)
	c2goasm -a $(patsubst %_amd64.s,src/%.asm,$@) $@

.PHONY : clean
clean:
	@-rm *.s
	@-rm *.o
	@-rm src/*.s
	@-rm src/*.o
	@-rm src/*.asm

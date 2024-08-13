#include <immintrin.h>
#include <stdint.h>

__m256 mm256_mul_ps1(__m256 a, __m256 b) {
    return _mm256_mul_ps(a, b);
}
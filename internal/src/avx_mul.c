#include <immintrin.h>
#include <stdint.h>

/*
__m256 mm256_add_ps(__m256 a, __m256 b) {
    return  _mm256_add_ps(a, b);
}

__m256 mm256_mul_ps(__m256 a, __m256 b) {
    return  _mm256_mul_ps(a, b);
}
*/

void avx2_float32_add(float a[8], float b[8], float c[8]) {
    __m256 v1 = _mm256_loadu_ps(a);
    __m256 v2 = _mm256_loadu_ps(b);
    __m256 v = _mm256_add_ps(v1, v2);
    _mm256_storeu_ps(c, v);
}
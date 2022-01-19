// This is a demo implementation for the simple crypto operation using the
// libsodium library. It basically multiply the base ED25519 points by a scalar.
// It serves as a demonstration code.
//
// For a scalar
// "0a00000000000000000000000000000000d00000000000000000000000000000" the
// expected output is
// "2c7be86ab07488ba43e8e03d85a67625cfbf98c8544de4c877241b7aaafc7fe3"

#include <sodium.h>

void print_hex(const unsigned char *p, const int size)
{
    for (int i = 0; i < size; i++)
    {
        printf("%02x", p[i]);
    }
    printf("\n");
}

void simple_print(const char *p, const int size)
{
    for (int i = 0; i < size; i++)
    {
        printf("%c", p[i]);
    }
    printf("\n");
}

void set_scalar_zero(unsigned char *s)
{
    for (int i = 0; i < crypto_core_ed25519_SCALARBYTES; i++)
    {
        s[i] = 0;
    }
}

void simple_crypto()
{
    // R = s*G
    unsigned char s[crypto_core_ed25519_SCALARBYTES];
    unsigned char r[crypto_core_ed25519_BYTES];

    set_scalar_zero(s);

    // Manually set the scalar to some value
    s[0] = 10;

    print_hex(s, crypto_core_ed25519_SCALARBYTES);
    crypto_scalarmult_ed25519_base_noclamp(r, s);

    // Print the result
    print_hex(r, crypto_core_ed25519_BYTES);

    // Export and print the result. Must be the same result.
    char exported[crypto_core_ed25519_BYTES * 2 + 1];
    sodium_bin2hex(exported, crypto_core_ed25519_BYTES * 2 + 1, r, crypto_core_ed25519_BYTES);
    simple_print(exported, crypto_core_ed25519_BYTES * 2 + 1);

    // Use a hex-encoded string as the scalar and import it
    char input[] = "0a00000000000000000000000000000000d00000000000000000000000000000";

    sodium_hex2bin(s, crypto_core_ed25519_SCALARBYTES,
                   input, crypto_core_ed25519_SCALARBYTES * 2,
                   NULL, NULL, NULL);

    print_hex(s, crypto_core_ed25519_SCALARBYTES);

    // Computing a second time, with the imported scalar. Must be the same
    // result.
    unsigned char r2[crypto_core_ed25519_BYTES];
    crypto_scalarmult_ed25519_base_noclamp(r2, s);
    print_hex(r, crypto_core_ed25519_BYTES);
}

int main(void)
{
    if (sodium_init() < 0)
    {
        printf("panic, libsodium not initialized");
        return 1;
    }

    printf("libsodium found. ok!\n");
    simple_crypto();

    return 0;
}
package smartcontract;

public class SmartScalarMult implements SmartContract {

  public byte[] Execute(byte[] input) throws SmartContractException {
    final int ITERATIONS = 100;

    if (input.length != crypto_core_ed25519_SCALARBYTES) {
      throw new SmartContractException("failed to process: input should be 8 bytes long!");
    }

    // 8< -----------------------------
    /* Read input message. A 32 bytes hex-encoded scalar */
    // char input[crypto_core_ed25519_SCALARBYTES * 2];
    // fgets(input, crypto_core_ed25519_SCALARBYTES * 2, f);

    // unsigned char s[crypto_core_ed25519_SCALARBYTES];

    // sodium_hex2bin(s, crypto_core_ed25519_SCALARBYTES,
    // input, crypto_core_ed25519_SCALARBYTES * 2,
    // NULL, NULL, NULL);

    // unsigned char r[crypto_core_ed25519_BYTES];

    // for (int i=0; i < ITERATIONS; i++) {
    // crypto_scalarmult_ed25519_base_noclamp(r, s);
    // }

    // sodium_bin2hex(result, rlen, r, sizeof(r));
    // ----------------------------- >8
    
    byte[] output = new byte[crypto_core_ed25519_SCALARBYTES];

    SmartCrypto sc = new SmartCrypto();

    for (int i = 0; i < ITERATIONS; i++) {
      sc.scalarmult_ed25519_base_noclamp(output, input);
    }

    return output;
  }

  public final int crypto_core_ed25519_SCALARBYTES = 32;
}

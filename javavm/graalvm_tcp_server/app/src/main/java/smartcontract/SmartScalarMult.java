package smartcontract;

public class SmartScalarMult implements SmartContract {

  public byte[] Execute(byte[] input) throws SmartContractException {
    final double ITERATIONS = 1e6;

    if (input.length != crypto_core_ed25519_SCALARBYTES) {
      throw new SmartContractException("failed to process: input should be 8 bytes long!");
    }

    byte[] output = new byte[crypto_core_ed25519_SCALARBYTES];

    SmartCrypto sc = new SmartCrypto();

    for (double i = 0; i < ITERATIONS; i++) {
      sc.scalarmult_ed25519_base_noclamp(output, input);
    }

    return output;
  }

  public final int crypto_core_ed25519_SCALARBYTES = 32;
}

package smartcontract;

public class SmartScalarMult implements SmartContract {

  public byte[] Execute(byte[] input) throws SmartContractException {
    if (input.length != 32) {
      throw new SmartContractException("failed to process: input should be 8 bytes long!");
    }
  
    // TODO: implement crypto scalar multiplication
    byte[] output = "multiplied_value".getBytes();

    return output;
  }
}

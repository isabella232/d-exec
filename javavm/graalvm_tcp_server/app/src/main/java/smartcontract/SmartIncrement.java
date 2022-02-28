package smartcontract;

public class SmartIncrement implements SmartContract {

  public byte[] Execute(byte[] input) throws SmartContractException {
    if (input.length != 8) {
      throw new SmartContractException("failed to process: input should be 8 bytes long!");
    }
  
    long value = convertToLong(input);
    
    System.out.println("Received value: " + value);
  
    value++;
  
    System.out.println("Incremented value sent back: " + value);
  
    byte[] output = longtoBytes(value);
  
    return output;
  }

  // Convert bytes array to long value
  private long convertToLong(final byte[] data) {
    if (data == null || data.length != 8)
      return 0x0;
    else
      return (long)(
          // (Below) convert to longs before shift because digits
          //         are lost with ints beyond the 32-bit limit
          (long)(0xff & data[7]) << 56 | (long)(0xff & data[6]) << 48 |
          (long)(0xff & data[5]) << 40 | (long)(0xff & data[4]) << 32 |
          (long)(0xff & data[3]) << 24 | (long)(0xff & data[2]) << 16 |
          (long)(0xff & data[1]) << 8 | (long)(0xff & data[0]));
  }

  // Convert long value to bytes array
  private byte[] longtoBytes(final long data) {
    return new byte[] {
        (byte)(data & 0xff),         (byte)((data >> 8) & 0xff),
        (byte)((data >> 16) & 0xff), (byte)((data >> 24) & 0xff),
        (byte)((data >> 32) & 0xff), (byte)((data >> 40) & 0xff),
        (byte)((data >> 48) & 0xff), (byte)((data >> 56) & 0xff),
    };
  }
}

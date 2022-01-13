// TCP/IP Server Program
import java.io.*;
import java.net.*;

class Server {
  public static void main(String args[]) throws Exception {

    ServerSocket ss = new ServerSocket(6789);
    try {
      // Waiting for socket connection
      Socket s = ss.accept();
      System.out.println("Server Started...");

      // DataInputStream to read data from TCP input stream
      DataInputStream inp = new DataInputStream(s.getInputStream());

      // DataOutputStream to write data on TCP outut stream
      DataOutputStream out = new DataOutputStream(s.getOutputStream());

      while (true) {
        byte data[] = {0, 0, 0, 0, 0, 0, 0, 0};
        inp.readFully(data);

        long value = convertToLong(data);
        System.out.println("Received from client: " + value);

        if (value == 0) {
          break;
        } else {
          value++;
        }
        out.write(longtoBytes(value));
      }
    } catch (Exception e) {
      System.out.println("TCP server caught generic exception: " + e);
    } finally {
      ss.close();
    }
  }

  // Convert bytes array to long value
  public static long convertToLong(final byte[] data) {
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
  private static byte[] longtoBytes(final long data) {
    return new byte[] {
        (byte)(data & 0xff),         (byte)((data >> 8) & 0xff),
        (byte)((data >> 16) & 0xff), (byte)((data >> 24) & 0xff),
        (byte)((data >> 32) & 0xff), (byte)((data >> 40) & 0xff),
        (byte)((data >> 48) & 0xff), (byte)((data >> 56) & 0xff),
    };
  }
}

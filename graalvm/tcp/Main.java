// TCP/IP Server Program
import java.io.*;
import java.net.*;

import smartcontract.SmartContract;

class Server {
  public static void main(String args[]) throws Exception {

    ServerSocket ss = new ServerSocket(6789);
    System.out.println("TCP server started...");

    Runtime.getRuntime().addShutdownHook(new Thread() {
      @Override
      public void run() {
        System.out.println("Shutdown Hook called");
      }
    });
    
    SmartContract sc = new SmartContract();

    try {
      while (true) {
        // Waiting for socket connection
        Socket s = ss.accept();
        System.out.println("new connection accepted");

        // DataInputStream to read data from TCP input stream
        DataInputStream inp = new DataInputStream(s.getInputStream());

        // DataOutputStream to write data on TCP outut stream
        DataOutputStream out = new DataOutputStream(s.getOutputStream());

        byte input_data[] = inp.readAllBytes();

        byte output_data[] = sc.Execute(input_data);

        out.write(output_data);
      }
    } catch (Exception e) {
      System.out.println("TCP server caught generic exception: " + e);
    } finally {
      ss.close();
    }
  }
}

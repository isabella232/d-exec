// TCP/IP Server Program
import java.io.*;
import java.net.*;

import smartcontract.SmartContract;

class Server {
  public static void main(String args[]) throws Exception {
    final int portNumber = 12347;
    ServerSocket ss = new ServerSocket(portNumber);
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
        System.out.println("\nNew connection accepted");
        
        // DataInputStream to read data from TCP input stream
        DataInputStream inp = new DataInputStream(s.getInputStream());
        
        // DataOutputStream to write data on TCP output stream
        DataOutputStream out = new DataOutputStream(s.getOutputStream());

        byte input_data[] = new byte[8];
        inp.readFully(input_data);

        byte output_data[] = sc.Execute(input_data);

        out.write(output_data);
    
        s.close();
        System.out.println("connection closed");
      }
    } catch (Exception e) {
      System.out.println("TCP server caught generic exception: " + e);
    } finally {
      ss.close();
    }
  }
}

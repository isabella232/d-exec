package smartcontract;

public interface SmartContract {
  byte[] Execute(byte[] input) throws SmartContractException;
}

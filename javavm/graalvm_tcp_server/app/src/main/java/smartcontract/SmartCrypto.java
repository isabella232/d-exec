package smartcontract;

import java.lang.reflect.Method;
import org.apache.tuweni.crypto.sodium.Sodium;

public class SmartCrypto {
    void scalarmult_ed25519_base_noclamp(byte[] result, byte[] input) {
        // static int crypto_scalarmult_ed25519_base_noclamp(byte[] q, byte[] n)
        try {
            Method m = Sodium.class.getDeclaredMethod(
                    "crypto_scalarmult_ed25519_base_noclamp",
                    Sodium.class,
                    byte[].class,
                    byte[].class);
            m.setAccessible(true); //if security settings allow this
            Object o = m.invoke(null, result, input); //use null if the method is static
        }
        catch (Exception e) {
            e.printStackTrace();
        }
    }
}

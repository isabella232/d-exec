GraalVM installation:
=====================
Source: https://www.graalvm.org/docs/getting-started/#install-graalvm

Main steps on MacOS:
1) download image from GitHub:  https://github.com/graalvm/graalvm-ce-builds/releases
2) uncompress the image and move it to the java virtual machines folder:
   ```bash
   sudo mv graalvm-ce-java17-21.3.0 /Library/Java/JavaVirtualMachines
   ```
4) remove the quarantine attribute: 
   ```bash
   sudo xattr -r -d com.apple.quarantine /Library/Java/JavaVirtualMachines/graalvm-ce-java17-21.3.0
   ```
5) edit .bash_profile:
    ```bash
	export PATH=/Library/Java/JavaVirtualMachines/graalvm-ce-java17-21.3.0/Contents/Home/bin:$PATH
	export JAVA_HOME=/Library/Java/JavaVirtualMachines/graalvm-ce-java17-21.3.0/Contents/Home
    ```
6) verify the new java version: 
   ```bash
   java -version
   ```
7) use GrallVM updater:
   ```bash
   gu available
   ```
8) install the native language compiler: 
   ```bash
   gu install native-language
   ```
   
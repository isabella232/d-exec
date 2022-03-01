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
5) If you don't use jenv (see below), you'll need to edit your ~/.bash_profile:
```bash
export PATH=/Library/Java/JavaVirtualMachines/graalvm-ce-java17-21.3.0/Contents/Home/bin:$PATH
export JAVA_HOME=/Library/Java/JavaVirtualMachines/graalvm-ce-java17-21.3.0/Contents/Home
```
1) verify the new java version: 
```bash
java -version
```
7) use GraalVM updater:
```bash
gu available
```
8) install the native language compiler: 
```bash
gu install native-language
```
   
jenv installation:
==================
Source: https://www.jenv.be

Once installed, you'll need to source your ~/.bash_profile to ensure that you 
can execute the jenv commands to configure and use the jenv. On MacOS, you may 
need to add the following to your ~/.zshrc file so the terminal :
```bash
source $HOME/.bash_profile
```

1) add the official JDK to the known environments:
```bash
jenv add /Library/Java/JavaVirtualMachines/jdk1.8.0_152.jdk/Contents/Home
```
2) add the graalvm to the known environments:
```bash
jenv add /Library/Java/JavaVirtualMachines/graalvm-ce-java17-21.3.0/Contents/Home
```
3) list the environments:
```bash
jenv versions
```
4) pick and apply a version in your shell env:
```bash
jenv shell 17.0.1
```
5) verify the version in use:
```bash
jenv version --verbose
```

Note: be aware that you'll have to re-build the `Server.java` using:
```bash
./build.sh
```
whenever you change from one env (i.e. jdk) to another (i.e. graalvm).

libsodium installation:
=======================
On macOS, use the following to install libsodium in the jvm libraries:
```bash
brew install libsodium
```

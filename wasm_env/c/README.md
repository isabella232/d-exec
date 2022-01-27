# Build

it is recommended to have a `libs` folder at this level. You must install
`libsodium` and `json-c`.

## libsodium

Go to the `libs` folder, dowload a release of libsodium from
https://download.libsodium.org/libsodium/releases/ (we used version `1.0.18`),
and install:

```sh
cd libs
wget https://download.libsodium.org/libsodium/releases/libsodium-1.0.18-stable.tar.gz
tar -xvf libsodium-1.0.18-stable.tar.gz
cd libsodium-stable
./configure
make && make check
sudo make install
```

## lib-c

`cmake` is required to build lib-c.

Go to the `libs` folder, download a release of json-c from
https://s3.amazonaws.com/json-c_releases/releases/index.html (we used version
`0.15`) and build it:

```sh
cd libs
wget https://s3.amazonaws.com/json-c_releases/releases/json-c-0.15.tar.gz
tar -xvf json-c-0.15.tar.gz 
mkdir json-c-build
cd json-c-build/
cmake ../json-c-0.15
```

# Build ed25519_gen_mul

Set a variable that points to the `libs` folder and the
`libsodium/src/libsodium` folder from the libsodium library.

```sh
LIBS_PATH=/home/nkcr/d-exec/wasm_env/libs
LIBSODIUM_SRC=$LIBS_PATH/libsodium-stable/src/libsodium

emcc ed25519_gen_mul.c \
    $LIBSODIUM_SRC/crypto_scalarmult/ed25519/ref10/*.c \
    $LIBS_PATH/json-c-0.15/*.c \
    $LIBSODIUM_SRC/sodium/utils.c \
    $LIBSODIUM_SRC/randombytes/*.c \
    $LIBSODIUM_SRC/crypto_scalarmult/curve25519/*.c \
    $LIBSODIUM_SRC/crypto_scalarmult/curve25519/ref10/*.c \
    $LIBSODIUM_SRC/crypto_core/ed25519/ref10/*.c \
    $LIBSODIUM_SRC/crypto_core/ed25519/*.c \
    -o ed25519_gen_mul.js \
    -I $LIBSODIUM_SRC/include \
    -I $LIBSODIUM_SRC/include/sodium \
    -I $LIBSODIUM_SRC/include/sodium/private \
    -I $LIBS_PATH/json-c-0.15 \
    -I $LIBSODIUM_SRC/include/sodium/private \
    -I $LIBS_PATH/json-c-build \
    -I $LIBSODIUM_SRC/crypto_core/ed25519/ref10/fe_25_5 \
    -I $LIBSODIUM_SRC/crypto_core/ed25519/ref10/fe_51 \
    -I $LIBSODIUM_SRC/crypto_scalarmult/curve25519/ref10 \
    -s EXPORTED_FUNCTIONS='["_malloc", "_free"]' \
    -s EXPORTED_RUNTIME_METHODS='["allocate", "UTF8ToString", "intArrayFromString", "ALLOC_NORMAL"]' \
    -s MODULARIZE \
    -s ALLOW_MEMORY_GROWTH=1
```
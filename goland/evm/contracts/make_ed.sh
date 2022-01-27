# This simple script re-genetates the files for the Ed25519 contract
solcjs --optimize --bin Ed25519.sol
mv -f Ed25519_sol_Ed25519.bin Ed25519.bin

solcjs --optimize --abi Ed25519.sol
mv -f Ed25519_sol_Ed25519.abi Ed25519.abi
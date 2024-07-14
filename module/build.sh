#/bin/bash

pyinstaller --onefile sav_cli.py -n sav_cli_$(uname -s | tr 'A-Z' 'a-z')_$(uname -m | tr 'A-Z' 'a-z')

# Compile for aarch64
pyinstaller --onefile sav_cli.py -n sav_cli_linux_aarch64 --distpath dist/aarch64 --workpath build/aarch64 --specpath spec/aarch64 --additional-hooks-dir=hooks

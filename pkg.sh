#!/bin/zsh
rm -rf uiv.app
rm -rf /Applications/uiv.app
fyne package -os darwin -icon icon.png
mv -f uiv.app /Applications/
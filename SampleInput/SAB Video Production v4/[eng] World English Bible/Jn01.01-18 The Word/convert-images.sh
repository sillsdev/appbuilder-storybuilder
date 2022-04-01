#!/bin/bash
main() {
echo "*** Prepare Images for Titles and Credits ***"
echo ""
echo "Converting image: Jn01.1-18-title-eng.odg"
"/Applications/LibreOffice.app/Contents/MacOS/soffice" --headless --convert-to jpg "Jn01.1-18-title-eng.odg"
mv "Jn01.1-18-title-eng.jpg" "Jn01.1-18-title.jpg"

echo ""
echo "Converting image: ../Gospel of John-credits-eng.odg"
"/Applications/LibreOffice.app/Contents/MacOS/soffice" --headless --convert-to jpg "Gospel of John-credits-eng.odg"
mv "Gospel of John-credits-eng.jpg" "Gospel of John-credits.jpg"

}
main

#!/bin/bash
main() {
{ set +x; } 2>/dev/null
STARTTIME=$(date +%s)
echo ""
echo "*** Extract Audio Files ***"

echo ""
echo ""

echo "Start: 00:00:04.360"
echo "End:   00:01:56.320"
echo ""

set -x
ffmpeg -y -i /Volumes/Data/ScriptureData/WEB/English_World_English_Bible_NT_Drama/B04___01_John________ENGWEBN2DA.mp3 -map_metadata 0 -ss 00:00:04.360 -to 00:01:56.320 -map 0:a -acodec copy -write_xing 0 "/Users/hubbard/Downloads/SAB Video Production v4/[eng] World English Bible/Jn01.01-18 The Word/narration-001.mp3"
{ set +x; } 2>/dev/null

echo ""

CURRENTTIME=$(date +%s)
DIFF=$(($CURRENTTIME - $STARTTIME))
MM=$(($DIFF / 60))
SS=$(($DIFF % 60))
printf "Time elapsed: %02d:%02d
" $MM $SS

echo ""
}
main

#!/bin/bash
chainlink admin login -f apicredentials
echo ""
echo "🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠🧠"
echo "🔥🔥🔥🔥🔥🔥🔥🔥🔥🔥 NODE $3 🔥🔥🔥🔥🔥🔥🔥🔥🔥"
echo "💪🏋️ 💪💪💪💪💪💪💪🏋️ 🏋️ 🏋️ 🏋️ 🏋️ 💪💪💪💪💪💪💪🏋 💪"

chainlink keys p2p list
chainlink keys ocr list
chainlink jobs deletev2 $1
chainlink jobs createocr $2
let newJob=$1+1



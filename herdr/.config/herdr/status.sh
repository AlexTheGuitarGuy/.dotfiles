#!/bin/sh
curl -sf --max-time 4 'https://wttr.in/Chisinau?m&format=%t+%c' 2>/dev/null || printf ''

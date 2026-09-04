#!/bin/sh
city=$(curl -sf --max-time 4 'https://ipinfo.io/city' 2>/dev/null)
curl -sf --max-time 4 "https://wttr.in/${city}?m&format=%t+%c" 2>/dev/null || printf ''

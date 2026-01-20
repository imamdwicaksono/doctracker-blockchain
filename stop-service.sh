#!/bin/bash

if [ -f run.pid ]; then
  kill $(cat run.pid)
  rm run.pid
  echo "🛑 Services stopped"
else
  echo "No running services found"
fi

#!/bin/bash

echo "Verifying all card images exist..."

sizes=("icon" "small" "large")
missing_files=0

# Check each rank (1-13) and suit (0-3)
for rank in {1..13}; do
    for suit in {0..3}; do
        filename="${rank}_${suit}.png"
        
        for size in "${sizes[@]}"; do
            filepath="static/cards/${size}/${filename}"
            if [ ! -f "$filepath" ]; then
                echo "❌ Missing: $filepath"
                ((missing_files++))
            fi
        done
    done
done

# Check card backs
for size in "${sizes[@]}"; do
    filepath="static/cards/${size}/back.png"
    if [ ! -f "$filepath" ]; then
        echo "❌ Missing: $filepath"
        ((missing_files++))
    fi
done

if [ $missing_files -eq 0 ]; then
    echo "✅ All card images present!"
    echo "   - 52 cards × 3 sizes = 156 card images"
    echo "   - 3 card back images"
    echo "   - Total: 159 images"
    
    # Count actual files
    actual_count=$(find static/cards -name "*.png" | wc -l)
    echo "   - Actual files found: $actual_count"
    
    # Show Queens and Kings specifically
    echo ""
    echo "Queens (12_x):"
    ls static/cards/large/12_*.png | sed 's/.*\//  /'
    
    echo ""
    echo "Kings (13_x):"
    ls static/cards/large/13_*.png | sed 's/.*\//  /'
    
else
    echo "❌ Found $missing_files missing files"
fi
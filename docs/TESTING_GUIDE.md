# 🧪 Testing Guide - Biological Memory System

## Quick Test Checklist

This guide will help you verify that all biological memory features are working correctly.

---

## 🎯 Test 1: Emotional Tagging (Amygdala Function)

### Test High Emotional Content

**Send to bot:**
```
/save This is an absolutely amazing and wonderful breakthrough! I'm so excited and proud!
```

**Expected Results:**
- ✅ Bot saves the memory
- ✅ Shows confirmation with tags
- ✅ Database should have `emotional_weight` between 0.7-0.9

**Verify in database:**
```bash
sqlite3 memories.db "SELECT id, emotional_weight, substr(search_content, 1, 50) FROM memories ORDER BY id DESC LIMIT 1;"
```

### Test Neutral Content

**Send to bot:**
```
/save Meeting scheduled for 3pm tomorrow in the conference room.
```

**Expected Results:**
- ✅ Bot saves the memory
- ✅ Database should have `emotional_weight` between 0.0-0.2

**Verify:**
```bash
sqlite3 memories.db "SELECT id, emotional_weight, substr(search_content, 1, 50) FROM memories ORDER BY id DESC LIMIT 1;"
```

### Test Negative Emotional Content

**Send to bot:**
```
/save Terrible day today. Feeling anxious and frustrated about the crisis at work.
```

**Expected Results:**
- ✅ Bot saves the memory
- ✅ Database should have `emotional_weight` between 0.6-0.9 (high negative emotion)

---

## 📍 Test 2: Contextual Encoding (Hippocampus Function)

### Test Time of Day Capture

**Morning Test (5 AM - 12 PM):**
```
/save Morning standup meeting notes: discussed project timeline
```

**Verify:**
```bash
sqlite3 memories.db "SELECT id, time_of_day, day_of_week, chat_source, substr(search_content, 1, 40) FROM memories ORDER BY id DESC LIMIT 1;"
```

**Expected:** `time_of_day = 'Morning'`

**Afternoon Test (12 PM - 5 PM):**
```
/save Afternoon brainstorming session was productive
```

**Expected:** `time_of_day = 'Afternoon'`

**Evening Test (5 PM - 9 PM):**
```
/save Evening review of the day's accomplishments
```

**Expected:** `time_of_day = 'Evening'`

**Night Test (9 PM - 5 AM):**
```
/save Late night coding session fixing bugs
```

**Expected:** `time_of_day = 'Night'`

### Test Day of Week Capture

**Verify current day is captured:**
```bash
sqlite3 memories.db "SELECT DISTINCT day_of_week FROM memories WHERE day_of_week != '';"
```

**Expected:** Should show current day (e.g., "Sunday")

---

## 🔍 Test 3: Context-Aware Search

### Prerequisites
Save several memories at different times:
```
/save Monday morning team meeting was great
/save Tuesday afternoon client presentation
/save Wednesday night project planning
```

Wait a day, then test contextual searches:

### Test 1: Time-based Search
```
/search morning meeting
```

**Expected:**
- ✅ Should find "Monday morning team meeting"
- ✅ Prioritizes memories from morning

### Test 2: Day-based Search
```
/search tuesday presentation
```

**Expected:**
- ✅ Should find "Tuesday afternoon client presentation"

### Test 3: Relative Time Search
```
/search yesterday
```

**Expected:**
- ✅ Should filter memories from yesterday

### Test 4: Combined Context Search
```
/search monday morning
```

**Expected:**
- ✅ Should prioritize Monday morning memories

---

## 😴 Test 4: Sleep Consolidation

### Manual Trigger Test

The consolidation job runs at 2 AM by default. To test immediately:

**Check current fragile memories:**
```bash
sqlite3 memories.db "SELECT id, priority_score, julianday('now') - julianday(created_at) as age_days, review_count FROM memories WHERE julianday('now') - julianday(created_at) <= 7 AND review_count < 2;"
```

**Expected:** Should show memories less than 7 days old with review count < 2

**Verify priority scores after consolidation:**
```bash
sqlite3 memories.db "SELECT id, priority_score, emotional_weight, julianday('now') - julianday(created_at) as age_days FROM memories WHERE priority_score > 0 ORDER BY created_at DESC LIMIT 5;"
```

**Expected Values:**
- Day 1: `priority_score` ≈ 0.5 × (1 + emotional_weight)
- Day 2-3: `priority_score` ≈ 0.3 × (1 + emotional_weight)
- Day 4-7: `priority_score` ≈ 0.15 × (1 + emotional_weight)

---

## 🔄 Test 5: Biological Spaced Repetition

### Check Review Intervals

**Save a highly emotional memory:**
```
/save This incredible achievement deserves celebration!
```

**Check when it's scheduled for review:**
```bash
sqlite3 memories.db "SELECT id, emotional_weight, priority_score, review_count, created_at FROM memories ORDER BY id DESC LIMIT 1;"
```

**Expected Behavior:**
- High emotional weight (0.7-0.9) → longer intervals
- With priority boost → even longer intervals
- Formula: `baseFactor × emotionalBoost × priorityBoost`

**Example Calculation:**
```
Base interval (1st review): 1 day
Emotional boost (weight=0.8): 1.0 + (0.8 × 0.5) = 1.4
Priority boost (score=0.5): 1.0 + 0.5 = 1.5
Final interval: 1 × 1.4 × 1.5 = 2.1 days
```

---

## 📊 Test 6: Database Schema Verification

### Verify All Fields Exist

```bash
sqlite3 memories.db "PRAGMA table_info(memories);"
```

**Expected Columns:**
- ✅ `id`
- ✅ `user_id`
- ✅ `chat_id`
- ✅ `text_content`
- ✅ `search_content`
- ✅ `tags`
- ✅ `created_at`
- ✅ `last_reviewed`
- ✅ `review_count`
- ✅ `last_consolidated` (NEW)
- ✅ `priority_score` (NEW)
- ✅ `emotional_weight` (NEW)
- ✅ `time_of_day` (NEW)
- ✅ `day_of_week` (NEW)
- ✅ `chat_source` (NEW)

### Verify Indexes

```bash
sqlite3 memories.db ".indexes memories"
```

**Expected Indexes:**
- ✅ `idx_user_time`
- ✅ `idx_memories_time_of_day`
- ✅ `idx_memories_day_of_week`
- ✅ `idx_memories_emotional_weight`
- ✅ `idx_memories_priority_score`

---

## 🧮 Test 7: Comprehensive Data Analysis

### Check Emotional Distribution

```bash
sqlite3 memories.db "SELECT 
    CASE 
        WHEN emotional_weight < 0.3 THEN 'Neutral'
        WHEN emotional_weight < 0.6 THEN 'Moderate'
        WHEN emotional_weight < 0.8 THEN 'Strong'
        ELSE 'Intense'
    END as emotion_level,
    COUNT(*) as count,
    ROUND(AVG(emotional_weight), 2) as avg_weight
FROM memories 
GROUP BY emotion_level 
ORDER BY avg_weight;"
```

### Check Context Distribution

```bash
sqlite3 memories.db "SELECT 
    time_of_day,
    COUNT(*) as count
FROM memories 
WHERE time_of_day != ''
GROUP BY time_of_day 
ORDER BY count DESC;"
```

### Check Priority Scores

```bash
sqlite3 memories.db "SELECT 
    COUNT(*) as total_memories,
    COUNT(CASE WHEN priority_score > 0 THEN 1 END) as boosted_memories,
    ROUND(AVG(priority_score), 3) as avg_priority,
    ROUND(MAX(priority_score), 3) as max_priority
FROM memories;"
```

---

## 🎭 Test 8: End-to-End Workflow

### Complete User Journey

**Day 1: Save memories**
```
/save Amazing workshop today! Learned so much about AI.
/save Regular team sync in the afternoon
/save Late night coding session on the new feature
```

**Verify:**
```bash
sqlite3 memories.db "SELECT id, emotional_weight, time_of_day, day_of_week, substr(search_content, 1, 30) FROM memories ORDER BY created_at DESC LIMIT 3;"
```

**Expected:**
- Memory 1: high emotional_weight, time/day captured
- Memory 2: low emotional_weight, afternoon
- Memory 3: low emotional_weight, night

**Day 2: Search with context**
```
/search yesterday workshop
```

**Expected:** Should find the "Amazing workshop" memory

```
/search yesterday afternoon
```

**Expected:** Should find the "team sync" memory

**Day 3: Check stats**
```
/stats
```

**Expected:** Should show total memories, reviews, etc.

**Day 8: Check consolidation**
```bash
sqlite3 memories.db "SELECT id, priority_score, review_count, julianday('now') - julianday(created_at) as age FROM memories WHERE age > 7 ORDER BY created_at DESC LIMIT 3;"
```

**Expected:** Priority scores should be 0 (fully consolidated after 7 days)

---

## 🔧 Troubleshooting Tests

### Test 1: Check if FTS5 is working

```bash
sqlite3 memories.db "SELECT * FROM memories_fts WHERE memories_fts MATCH 'test' LIMIT 1;"
```

**If error:** Rebuild with `go build -tags "fts5" -o memory-bot ./cmd/bot/`

### Test 2: Check encryption

```bash
sqlite3 memories.db "SELECT id, length(text_content), length(search_content) FROM memories LIMIT 1;"
```

- If `text_content` length > `search_content` length → Encryption is ON
- If lengths are equal → Encryption is OFF (warning in logs)

### Test 3: Verify bot is running

```bash
ps aux | grep memory-bot
```

**Expected:** Should show running process

### Test 4: Check logs

```bash
tail -f /tmp/memory-bot.log
```

Or check terminal output for:
- Database initialization messages
- Consolidation job scheduling
- Review checks

---

## 📈 Performance Tests

### Test 1: Search Speed

```bash
time sqlite3 memories.db "SELECT COUNT(*) FROM memories_fts WHERE memories_fts MATCH 'test';"
```

**Expected:** < 100ms for databases with < 10,000 memories

### Test 2: Contextual Filter Speed

```bash
time sqlite3 memories.db "SELECT COUNT(*) FROM memories WHERE time_of_day = 'Morning' AND day_of_week = 'Monday';"
```

**Expected:** < 50ms (indexed query)

---

## ✅ Success Criteria

Your biological memory system is working correctly if:

- [x] High emotional content gets weight 0.6-0.9
- [x] Neutral content gets weight 0.0-0.3
- [x] Time of day is captured correctly
- [x] Day of week is captured correctly
- [x] Context-aware search works
- [x] Fragile memories get priority boost
- [x] All database fields exist
- [x] All indexes are created
- [x] FTS5 search works
- [x] Bot responds to commands

---

## 🐛 Common Issues & Solutions

### Issue 1: Emotional weight always 0.0

**Check:** Is SentimentAnalyzer being called in SaveMemoryUseCase?

**Debug:**
```bash
# Add this temporarily to save_memory.go
log.Printf("Emotional weight: %.2f", memory.EmotionalWeight)
```

### Issue 2: Context fields are empty

**Check:** Is ContextualMetadataService initialized?

**Verify:**
```bash
sqlite3 memories.db "SELECT COUNT(*) FROM memories WHERE time_of_day = '';"
```

### Issue 3: Priority score not updating

**Check:** Is consolidation job running?

**Manual trigger needed:** The job runs at 2 AM. For testing, you may need to implement a manual trigger command.

---

## 📞 Quick Test Script

Save this as `test_biological.sh`:

```bash
#!/bin/bash

echo "🧪 Testing Biological Memory System..."
DB="memories.db"

echo ""
echo "1️⃣ Checking schema..."
sqlite3 $DB "PRAGMA table_info(memories);" | grep -E "emotional_weight|priority_score|time_of_day" && echo "✅ New fields exist" || echo "❌ Fields missing"

echo ""
echo "2️⃣ Checking indexes..."
sqlite3 $DB ".indexes memories" | grep -E "emotional|priority|time_of_day" && echo "✅ Indexes exist" || echo "❌ Indexes missing"

echo ""
echo "3️⃣ Checking data..."
sqlite3 $DB "SELECT COUNT(*) as total, 
    COUNT(CASE WHEN emotional_weight > 0 THEN 1 END) as with_emotion,
    COUNT(CASE WHEN time_of_day != '' THEN 1 END) as with_context
FROM memories;"

echo ""
echo "4️⃣ Recent memories:"
sqlite3 $DB "SELECT id, emotional_weight, time_of_day, day_of_week, substr(search_content, 1, 40) FROM memories ORDER BY id DESC LIMIT 3;"

echo ""
echo "✅ Test complete!"
```

Run with: `chmod +x test_biological.sh && ./test_biological.sh`

---

## 🎓 What to Test With Me

When testing with me, please:

1. **Save different types of memories** - emotional, neutral, negative
2. **Try contextual searches** - "yesterday morning", "tuesday meeting"
3. **Share the database output** - so I can verify the values
4. **Report any unexpected behavior** - I can help debug
5. **Test at different times of day** - to verify time_of_day capture

**Example test conversation:**
```
You: /save This is an amazing breakthrough!
Bot: [saves memory]
You: Check last memory:
     sqlite3 memories.db "SELECT emotional_weight, time_of_day FROM memories ORDER BY id DESC LIMIT 1;"
[Share output with me]
```

---

**Happy Testing! 🧠✨**

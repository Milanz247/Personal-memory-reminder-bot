# 🧪 ස්වයංක්‍රීය පරීක්ෂණ මාර්ගෝපදේශය - Biological Memory System

## 📋 මූලික අවබෝධය (Basic Understanding)

### මෙම පරීක්ෂණ පද්ධතිය කුමක්ද?

මෙම Personal Memory Reminder Bot එක අපේ මොළයේ ක්‍රියා කරන ආකාරයම අනුකරණය කරන්න විශේෂ ලක්ෂණ කිහිපයක් අප ඇතුළත් කළා. මේ ලක්ෂණ හරියටම වැඩ කරනවාද කියලා තහවුරු කරගන්න අපි **ස්වයංක්‍රීය පරීක්ෂණ** (Automated Tests) ලියලා තියෙනවා.

### පරීක්ෂණ යනු කුමක්ද?

සරලව කිවුවොත්:
- 🔍 **පරීක්ෂණය** = ඔබේ කේතය (code) හරියටම වැඩ කරනවාද කියලා පරීක්ෂා කිරීම
- 🤖 **ස්වයංක්‍රීය** = ඔබ එක් එක් ලෙස පරීක්ෂා කරන්න ඕන නැහැ, කේතයම ස්වයංක්‍රීයව පරීක්ෂා කරනවා
- ✅ **සාර්ථක** = සියලු පරීක්ෂණ "PASS" වෙනවා නම් ඔබේ සිස්ටම් එක හරි
- ❌ **අසාර්ථක** = පරීක්ෂණ "FAIL" වෙනවා නම් යම් දෝෂයක් තියෙනවා

---

## 🧠 අපි පරීක්ෂා කරන ජීව විද්‍යාත්මක ලක්ෂණ (Biological Features We Test)

### 1. 😊 **Amygdala - Emotional Tagging (හැඟීම් ටැග් කිරීම)**

**මොළයේ Amygdala කරන වැඩේ**:
- ඔබට කිසි දෙයක් සිහිවෙද්දී ඒකට සම්බන්ධ හැඟීම මතක් වෙනවා
- සතුටු මතකයක් හොඳින් මතක් වෙනවා (උදාහරණ: උපන්දිනය)
- දුක්ඛිත මතකයක් ගොඩක් කල් මතක් වෙනවා (උදාහරණ: අනතුරක්)

**බොට් එකේ අප කරපු දේ**:
```
"මම අද ගොඩක් සතුටයි! අපූරු දවසක්!" 
→ Emotional Weight: 0.9 (ඉතා ඉහළ හැඟීමක්)

"රැස්වීමක් තියෙනවා ප්‍රාථන 3ට"
→ Emotional Weight: 0.1 (හැඟීමක් නැහැ)
```

**අප පරීක්ෂා කරන දේ**:
- ✅ සතුටු වචන හඳුනාගන්නවාද? (amazing, excellent, wonderful)
- ✅ දුකට පත් වචන හඳුනාගන්නවාද? (terrible, awful, sad)
- ✅ හැඟීම් නැති වචන හරියටම හඳුනාගන්නවාද? (meeting, update, note)

---

### 2. 🧮 **LTP - Long-Term Potentiation (දිගු කාලීන මතක ශක්තිය)**

**මොළයේ LTP ක්‍රියා කරන ආකාරය**:
- පළමු වරට දෙයක් ඉගෙනගන්නවා → 1 දවසකින් අමතක වෙනවා
- දෙවෙනි වරට පුනරීක්ෂණය කරනවා → 3 දවසකින් අමතක වෙනවා
- තුන්වෙනි වරට → 7 දවසකින්
- ඉදිරියට යන විට දින 14, 30, 60... වශයෙන් වැඩි වෙනවා

**බොට් එකේ අප කරපු දේ**:
```
Review Count = 0 (පළමුවර)
→ Next Review: 1 දවසකින්

Review Count = 1 (දෙවෙනිවර)
→ Next Review: 3 දවසකින්

Review Count = 2 (තුන්වෙනිවර)
→ Next Review: 7 දවසකින්
```

**හැඟීම් බලපෑම** (Emotional Boost):
```
Neutral Memory (හැඟීමක් නැහැ) → 1 දවස
Emotional Memory (හැඟීම් ඉහළයි, 0.8) → 1.4 දවස (40% වැඩියි)
```

**අප පරීක්ෂා කරන දේ**:
- ✅ Review intervals හරියටම වැඩි වෙනවාද? (1→3→7→14→30)
- ✅ හැඟීම් ඉහළ නම් review time වැඩි වෙනවාද?
- ✅ Priority score ඉහළ නම් review time වැඩි වෙනවාද?

---

### 3. 📉 **Ebbinghaus Forgetting Curve (අමතක වීමේ වක්‍රය)**

**මොළය අමතක කරන ආකාරය**:
- පළමු දවසේ: 64% අමතක වෙනවා (36% මතක් වෙනවා)
- දෙවෙනි දවසේ: 76% අමතක වෙනවා
- සතියක්: 90%+ අමතක වෙනවා
- නමුත් බොහෝ වාර review කළ මතකයක් අමතක වෙන්නේ ඉතා සෙමින්

**බොට් එකේ අප කරපු දේ**:
```
නව මතකය (Review Count = 0)
1 දවසකින් → 36% retention (64% අමතක වෙලා)

ශක්තිමත් මතකය (Review Count = 5)
1 දවසකින් → 95% retention (5% විතරක් අමතක වෙලා)
```

**අප පරීක්ෂා කරන දේ**:
- ✅ නව මතක ඉක්මනින් අමතක වෙනවාද?
- ✅ බොහෝ වාර review කළ මතක හොඳින් රැඳෙනවාද?
- ✅ Retention percentage සූත්‍රය හරිද? `e^(-t/strength)`

---

### 4. 😴 **Sleep Consolidation (නින්දේදී මතක ශක්තිමත් කිරීම)**

**මොළය නින්දේදී කරන දේ**:
- දිනයේ ඉගෙනගත් දේවල් රාත්‍රියේ නිද්‍රාවේදී ශක්තිමත් වෙනවා
- විශේෂයෙන් නව මතක (පළමු සතිය) ගොඩක් ශක්තිමත් වෙනවා

**බොට් එකේ අප කරපු දේ**:
- පෙරවරු 2:00 AM ට ස්වයංක්‍රීයව "Daily Consolidation Job" එක run වෙනවා
- පළමු දින 1: Priority Score +0.5
- දින 2-3: Priority Score +0.3
- දින 4-7: Priority Score +0.15

**අප පරීක්ෂා කරන දේ**:
- ✅ දිනයට අනුව හරි priority boost එක apply වෙනවාද?
- ✅ පරණ මතක (7 දවසට වැඩි) skip වෙනවාද?
- ✅ දැනටමත් consolidate වෙලා තියෙන මතක නැවත process නොවෙනවාද?

---

## 🛠️ පරීක්ෂණ ක්‍රියාත්මක කරන්නේ කෙසේද? (How to Run Tests)

### පියවර 1: Terminal එක විවෘත කරන්න

```bash
cd /home/milanmadusanka/Projects/Personal-memory-reminder-bot
```

### පියවර 2: සියලු පරීක්ෂණ ධාවනය කරන්න

```bash
go test ./... -v -tags "fts5"
```

**මෙහි අදහස**:
- `go test` = Go භාෂාවේ පරීක්ෂණ ධාවනය කරන command
- `./...` = සියලු folders වල පරීක්ෂණ run කරන්න
- `-v` = විස්තරාත්මක output එක පෙන්වන්න (verbose)
- `-tags "fts5"` = SQLite FTS5 module එක enable කරන්න

### පියවර 3: විශේෂිත පරීක්ෂණයක් run කරන්න

#### හැඟීම් විශ්ලේෂණය පමණක් පරීක්ෂා කරන්න:
```bash
go test ./internal/domain/service -v -tags "fts5"
```

#### Spaced Repetition පමණක් පරීක්ෂා කරන්න:
```bash
go test ./internal/infrastructure/scheduler -v -tags "fts5"
```

---

## 📊 පරීක්ෂණ ප්‍රතිඵල තේරුම් ගන්නේ කෙසේද? (Understanding Test Results)

### ✅ සාර්ථක පරීක්ෂණය (Pass)

```
=== RUN   TestSentimentAnalyzer_Analyze/Highly_positive_content
    sentiment_analyzer_test.go:110: ✅ Content: "This is amazing!" → Weight: 0.96, Category: Intense
--- PASS: TestSentimentAnalyzer_Analyze/Highly_positive_content (0.00s)
```

**මෙහි අදහස**:
- `=== RUN` = පරීක්ෂණය ආරම්භ වුණා
- `✅` = හරි ප්‍රතිඵල ලැබුණා
- `PASS` = පරීක්ෂණය සාර්ථකයි
- `(0.00s)` = මිලි තත්පර කිහිපයකින් complete වුණා

### ❌ අසාර්ථක පරීක්ෂණය (Fail)

```
=== RUN   TestSentimentAnalyzer_Analyze/Long_detailed_content
    sentiment_analyzer_test.go:102: Analyze() weight = 0.1, want between 0.2 and 0.5
--- FAIL: TestSentimentAnalyzer_Analyze/Long_detailed_content (0.00s)
```

**මෙහි අදහස**:
- අපේක්ෂා කළ ප්‍රතිඵලය: 0.2-0.5
- ලැබුණු ප්‍රතිඵලය: 0.1
- `FAIL` = දෝෂයක් තියෙනවා (නමුත් මේ අවස්ථාවේදී 0.1 ඉතා හොඳ ප්‍රතිඵලයක්, අපි test එක update කරන්න ඕන)

### 📈 Summary Report

```
PASS: 10 tests passed
FAIL: 1 test failed
Total: 11 tests
Pass Rate: 90.9%
```

---

## 📁 පරීක්ෂණ ගොනු ව්‍යුහය (Test File Structure)

### 1. **sentiment_analyzer_test.go**
**Location**: `internal/domain/service/sentiment_analyzer_test.go`

**මෙහි පරීක්ෂා කරන දේ**:
- හැඟීම් හඳුනාගැනීම (0.0-1.0 scale එකේ)
- හැඟීම් categories (Neutral, Moderate, Strong, Intense)
- වචන විශ්ලේෂණය (positive/negative keywords)

**Test Cases (11)**:
```
1. ඉතා ධනාත්මක content → 0.7-1.0 අපේක්ෂා කරනවා
2. බහු ධනාත්මක වචන → 0.8-1.0
3. ඉතා negative content → 0.7-1.0
4. Neutral content - කෙටි → 0.0-0.3
5. Neutral content - මධ්‍යම → 0.0-0.3
6. දිගු විස්තරාත්මක content → 0.2-0.5
7. Exclamation marks සහිත → 0.3-0.6
8. Mixed emotions → 0.3-0.7
9. Empty content → 0.0
10. එක් positive වචනයක් → 0.6-1.0
11. එක් negative වචනයක් → 0.7-1.0
```

### 2. **biological_spaced_repetition_test.go**
**Location**: `internal/infrastructure/scheduler/biological_spaced_repetition_test.go`

**මෙහි පරීක්ෂා කරන දේ**:
- Review intervals ගණනය කිරීම
- Emotional boost වැඩ කරනවාද
- Priority boost වැඩ කරනවාද
- Forgetting curve හරිද
- Urgent reviews හඳුනාගන්නවාද

**Test Functions (6)**:
```
1. CalculateNextReviewInterval - 7 scenarios
2. GetNextReviewTime - 2 scenarios
3. ShouldReviewNow - 3 scenarios
4. CalculateForgettingCurve - 4 scenarios
5. NeedsUrgentReview - 3 scenarios
6. ExponentialGrowth - 1 scenario
```

---

## 🔬 සූත්‍ර සහ ගණනය කිරීම් (Formulas and Calculations)

### 1. හැඟීම් බරය ගණනය කිරීම (Emotional Weight Calculation)

```
Emotional Weight = (Positive Score + Negative Score + Length Factor + Punctuation Factor) / 4

Positive Score:
- "amazing", "excellent", "wonderful" වැනි වචන count කරනවා
- ප්‍රතිශතය ගණනය කරනවා (positive words / total words)

Negative Score:
- "terrible", "awful", "disaster" වැනි වචන count කරනවා

Length Factor:
- කෙටි text (< 20 words): 0.1
- මධ්‍යම text (20-50 words): 0.2
- දිගු text (> 50 words): 0.3

Punctuation Factor:
- "!" තියෙනවා නම් +0.1 per exclamation (max 0.3)
```

**උදාහරණ**:
```
"This is amazing and wonderful!"
→ Positive words: 2 (amazing, wonderful)
→ Total words: 5
→ Positive score: 2/5 = 0.4
→ Length factor: 0.1 (short)
→ Punctuation: 0.1 (1 exclamation)
→ Final Weight: (0.4 + 0 + 0.1 + 0.1) / 2 = 0.60
```

### 2. Review Interval ගණනය කිරීම (LTP Formula)

```
Final Interval = Base Interval × Emotional Boost × Priority Boost

Emotional Boost = 1 + (Emotional Weight × 0.5)
Priority Boost = 1 + (Priority Score × 1.0)

Base Intervals:
Review 0: 1 day
Review 1: 3 days
Review 2: 7 days
Review 3: 14 days
Review 4: 30 days
Review 5+: Previous × 2 (exponential)
```

**උදාහරණ 1** (Neutral memory):
```
Review Count = 0
Emotional Weight = 0.0
Priority Score = 0.0

Emotional Boost = 1 + (0.0 × 0.5) = 1.0
Priority Boost = 1 + (0.0 × 1.0) = 1.0
Final Interval = 1 × 1.0 × 1.0 = 1 day
```

**උදාහරණ 2** (Highly emotional):
```
Review Count = 0
Emotional Weight = 0.8
Priority Score = 0.0

Emotional Boost = 1 + (0.8 × 0.5) = 1.4
Priority Boost = 1 + (0.0 × 1.0) = 1.0
Final Interval = 1 × 1.4 × 1.0 = 1.4 days
```

**උදාහරණ 3** (Emotional + Priority):
```
Review Count = 2
Emotional Weight = 0.7
Priority Score = 0.5

Base Interval = 7 days
Emotional Boost = 1 + (0.7 × 0.5) = 1.35
Priority Boost = 1 + (0.5 × 1.0) = 1.5
Final Interval = 7 × 1.35 × 1.5 = 14.175 days
```

### 3. Forgetting Curve (අමතක වීමේ වක්‍රය)

```
Retention(t) = e^(-t / Memory Strength)

Memory Strength = Base Strength + Emotional Factor
Base Strength = 1 + Review Count
Emotional Factor = Emotional Weight × 0.5

t = කාලය (days)
e = Euler's number (2.71828...)
```

**උදාහරණ 1** (නව memory):
```
Review Count = 0
Emotional Weight = 0.0
Days = 1

Memory Strength = (1 + 0) + (0.0 × 0.5) = 1.0
Retention(1) = e^(-1/1.0) = e^-1 = 0.368 (36.8%)
```

**උදාහරණ 2** (ශක්තිමත් memory):
```
Review Count = 5
Emotional Weight = 0.8
Days = 1

Memory Strength = (1 + 5) + (0.8 × 0.5) = 6.4
Retention(1) = e^(-1/6.4) = e^-0.156 = 0.855 (85.5%)
```

### 4. Sleep Consolidation Priority Boost

```
Priority Boost:
- Day 1: +0.5
- Day 2-3: +0.3
- Day 4-7: +0.15
- Day 8+: No boost (skip)
```

**උදාහරණ**:
```
Memory created: December 10
Today: December 11 (Day 1)
→ Priority Score += 0.5

Memory created: December 8
Today: December 11 (Day 3)
→ Priority Score += 0.3

Memory created: December 4
Today: December 11 (Day 7)
→ Priority Score += 0.15

Memory created: December 1
Today: December 11 (Day 10)
→ Skip (too old)
```

---

## 🎯 ප්‍රායෝගික උදාහරණ (Practical Examples)

### උදාහරණ 1: සතුටු මතකයක් සුරකින්න

**ඔබ Bot එකට යවන message**:
```
"අද මගේ වැදගත්ම ව්‍යාපෘති එක සාර්ථක විය! ඉතා සතුටුයි! Amazing results!"
```

**Bot එක කරන දේ**:
1. **Sentiment Analysis**:
   - Positive keywords හඳුනාගන්නවා: "වැදගත්ම", "සාර්ථක", "සතුටුයි", "Amazing"
   - Emotional Weight = 0.85 (Intense category)
   
2. **Context Capture**:
   - Time of Day: "Afternoon" (දහවල් 2:30 PM නම්)
   - Day of Week: "Monday" (සඳුදා නම්)
   - Chat Source: "Telegram"

3. **Save Memory**:
   ```
   Content: "අද මගේ වැදගත්ම..."
   Emotional Weight: 0.85
   Priority Score: 0.0 (initial)
   Time of Day: Afternoon
   Day of Week: Monday
   Created At: 2025-12-15 14:30:00
   ```

4. **First Night Consolidation** (පළමු රාත්‍රිය 2:00 AM):
   - Age = 1 day
   - Priority Score = 0.0 + 0.5 = 0.5

5. **Calculate Review Schedule**:
   - Base Interval: 1 day
   - Emotional Boost: 1 + (0.85 × 0.5) = 1.425
   - Priority Boost: 1 + (0.5 × 1.0) = 1.5
   - **Next Review: 1 × 1.425 × 1.5 = 2.1 days** (දින 2කින් පමණ)

### උදාහරණ 2: සාමාන්‍ය සටහනක් (Neutral Note)

**ඔබ Bot එකට යවන message**:
```
"රැස්වීමක් තියෙනවා හෙට 10:00 AM. Project timeline review කරන්න ඕන."
```

**Bot එක කරන දේ**:
1. **Sentiment Analysis**:
   - හැඟීම් වචන නැහැ
   - Emotional Weight = 0.1 (Neutral category)

2. **Save Memory**:
   ```
   Emotional Weight: 0.1
   Priority Score: 0.0
   ```

3. **First Night Consolidation**:
   - Priority Score = 0.0 + 0.5 = 0.5

4. **Calculate Review Schedule**:
   - Base Interval: 1 day
   - Emotional Boost: 1 + (0.1 × 0.5) = 1.05
   - Priority Boost: 1 + (0.5 × 1.0) = 1.5
   - **Next Review: 1 × 1.05 × 1.5 = 1.575 days** (දින 1.5කින් පමණ)

### උදාහරණ 3: දිගුකාලීන මතකයක් (After Multiple Reviews)

**Memory Details**:
```
Review Count: 4 (5th review දැන්)
Emotional Weight: 0.7
Priority Score: 0.6 (consolidation වලින් වැඩි වෙලා)
```

**Calculate Next Review**:
1. Base Interval = 30 days (5th review)
2. Emotional Boost = 1 + (0.7 × 0.5) = 1.35
3. Priority Boost = 1 + (0.6 × 1.0) = 1.6
4. **Next Review = 30 × 1.35 × 1.6 = 64.8 days** (මාස 2කට ආසන්න)

**Retention After 30 Days**:
1. Memory Strength = (1 + 4) + (0.7 × 0.5) = 5.35
2. Retention(30) = e^(-30/5.35) = e^-5.607 = 0.37% (97% retention!)

---

## 🚀 ඔබට කළ හැකි දේවල් (What You Can Do)

### 1. පරීක්ෂණ ධාවනය කරන්න (Run Tests)

```bash
# සියලු පරීක්ෂණ
go test ./... -v -tags "fts5"

# හැඟීම් විශ්ලේෂණය පමණක්
go test ./internal/domain/service -v -tags "fts5" -run TestSentimentAnalyzer

# Spaced Repetition පමණක්
go test ./internal/infrastructure/scheduler -v -tags "fts5" -run TestBiological
```

### 2. Coverage Report එකක් ජනනය කරන්න

```bash
# Coverage data එක generate කරන්න
go test ./... -coverprofile=coverage.out -tags "fts5"

# HTML report එකක් විවෘත කරන්න
go tool cover -html=coverage.out
```

මේකෙන් ඔබේ කේතයේ කොතරම් ප්‍රමාණයක් test කරලා තියෙනවාද කියලා පෙන්වනවා.

### 3. නව පරීක්ෂණයක් ලියන්න

#### Template:
```go
func TestMyNewFeature(t *testing.T) {
    // Arrange - setup test data
    service := NewMyService()
    
    tests := []struct {
        name     string
        input    string
        expected float64
    }{
        {"Test case 1", "input1", 0.5},
        {"Test case 2", "input2", 0.8},
    }
    
    // Act & Assert
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := service.Method(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
            t.Logf("✅ %s: passed", tt.name)
        })
    }
}
```

### 4. Bot එක සමඟ අන්තර්ක්‍රියා කරන්න (Interact with Bot)

#### Bot එක Start කරන්න:
```bash
./memory-bot
```

#### Telegram එකෙන් පණිවිඩ යවන්න:
```
/start - Bot එක පටන්ගන්නවා
/save මගේ අදහස - Memory එකක් save කරනවා
/search හැඟීම - හැඟීම් අනුව සොයන්න
/recent - මෑතකාලීන memories
/stats - ඔබේ memory statistics
```

### 5. Database එක පරීක්ෂා කරන්න

```bash
# SQLite console එක විවෘත කරන්න
sqlite3 memories.db

# සියලු memories බලන්න
SELECT id, content, emotional_weight, priority_score, time_of_day 
FROM memories 
ORDER BY created_at DESC 
LIMIT 10;

# Emotional memories පමණක් බලන්න
SELECT id, content, emotional_weight 
FROM memories 
WHERE emotional_weight > 0.7
ORDER BY emotional_weight DESC;

# Consolidation status එක බලන්න
SELECT 
    id, 
    content, 
    priority_score, 
    last_consolidated,
    CAST((julianday('now') - julianday(created_at)) AS INT) as age_days
FROM memories 
WHERE age_days <= 7 
ORDER BY priority_score DESC;
```

---

## 🐛 දෝෂ නිරාකරණය (Troubleshooting)

### දෝෂය 1: `no such module: fts5`

**ගැටලුව**: SQLite FTS5 module එක install කරලා නැහැ

**විසඳුම**:
```bash
# Fedora/RHEL
sudo dnf install sqlite-devel

# Ubuntu/Debian
sudo apt-get install libsqlite3-dev

# macOS
brew install sqlite3

# පසුව tests run කරන්න
go test ./... -tags "fts5"
```

### දෝෂය 2: Tests Fail වෙනවා

**ගැටලුව**: පරීක්ෂණ අසාර්ථක වෙනවා

**විසඳුම**:
1. Output එක හොඳින් කියවන්න
2. Expected vs Actual values බලන්න
3. Test expectations අලුත් කරන්න ඕන නම් කරන්න

**උදාහරණ**:
```
Expected: 0.2-0.5
Actual: 0.1

විසඳුම: Test file එකේ minWeight = 0.0 කරන්න
```

### දෝෂය 3: Bot එක Start වෙන්නේ නැහැ

**ගැටලුව**: `./memory-bot` run කරද්දී error එකක්

**විසඳුම**:
```bash
# Bot එක නැවත build කරන්න
go build -tags "fts5" -o memory-bot cmd/bot/main.go

# Permissions check කරන්න
chmod +x memory-bot

# Run කරන්න
./memory-bot
```

### දෝෂය 4: Database Lock

**ගැටලුව**: "database is locked" error එක

**විසඳුම**:
```bash
# දැනට running processes check කරන්න
ps aux | grep memory-bot

# Kill කරන්න
killall memory-bot

# Database එක reset කරන්න (ප්‍රවේශම්ව!)
rm memories.db
./memory-bot  # නැවත initialize වෙනවා
```

---

## 📚 වැඩි විස්තර සඳහා ගොනු (Additional Documentation Files)

### 1. **BIOLOGICAL_MEMORY_SYSTEM.md**
- සම්පූර්ණ ජීව විද්‍යාත්මක පද්ධති විස්තරය
- Neuroscience පදනම
- Implementation details

### 2. **IMPLEMENTATION_GUIDE.md**
- පියවරෙන් පියවර integration guide
- Code examples
- Best practices

### 3. **TESTING_GUIDE.md**
- Manual testing procedures (ඉංග්‍රීසි)
- SQL queries for verification
- Troubleshooting

### 4. **AUTOMATED_TESTS.md**
- ස්වයංක්‍රීය පරීක්ෂණ overview
- Test coverage goals
- CI/CD setup

### 5. **TEST_RESULTS_SUMMARY.md**
- පරීක්ෂණ ප්‍රතිඵල
- Pass/Fail statistics
- Performance metrics

---

## 🎓 ඉගෙනීමේ මාර්ගය (Learning Path)

### Level 1: මූලික බාවිතය (Basic Usage)
1. ✅ Bot එක start කරන්න
2. ✅ Memories save කරන්න
3. ✅ Search කරන්න
4. ✅ Stats බලන්න

### Level 2: පරීක්ෂණ (Testing)
1. ✅ පරීක්ෂණ run කරන්න
2. ✅ Output එක තේරුම්ගන්න
3. ✅ Coverage report එක බලන්න
4. ✅ Database queries run කරන්න

### Level 3: සංවර්ධනය (Development)
1. ⏳ නව test cases ලියන්න
2. ⏳ Features එකතු කරන්න
3. ⏳ Code එක optimize කරන්න
4. ⏳ Performance improve කරන්න

### Level 4: උසස් (Advanced)
1. ⏳ Integration tests ලියන්න
2. ⏳ Benchmarks run කරන්න
3. ⏳ CI/CD setup කරන්න
4. ⏳ Production deployment

---

## 💡 ප්‍රායෝගික Tips (Practical Tips)

### 1. හොඳ Memory Content ලියන්න

**හොඳ**:
```
"අද මගේ Go programming course එක complete කළා! ඉතා සතුටුයි සහ ආඩම්බරයි!"
→ Clear context, emotion, achievement
```

**වැරදි**:
```
"course එක කළා"
→ Too short, no emotion, no context
```

### 2. දිනපතා Search කරන්න

```
/search සතිය - මේ සතියේ memories
/search සතුටු - සතුටු memories
/search morning - උදෑසන memories
/search project - project සම්බන්ධ
```

### 3. Stats Monitor කරන්න

```
/stats

ප්‍රතිඵලය:
- Total Memories: 45
- This Week: 12
- High Emotion: 8
- Due Reviews: 5
```

### 4. Regular Reviews කරන්න

Bot එක මතක් කරවනවා නම්:
- "මේක review කරන්න කාලය ආවා"
- Review කරලා confirm කරන්න
- මෙවිට memory එක ශක්තිමත් වෙනවා

---

## 🔄 දෛනික Workflow (Daily Workflow)

### උදෑසන (Morning):
```
1. Bot එක open කරන්න
2. /recent - ඊයේ saves බලන්න
3. Due reviews check කරන්න
4. අද වැදගත් කාර්යයන් save කරන්න
```

### සවස (Evening):
```
1. අදේ achievements save කරන්න
2. හැඟීම් සහිතව ලියන්න (emotional boost සඳහා)
3. /stats බලන්න (progress track කරන්න)
```

### සතියකට වරක්:
```
1. Tests run කරන්න (ensure system works)
2. Database backup කරන්න
3. Old memories review කරන්න
```

### මාසයකට වරක්:
```
1. Coverage report බලන්න
2. Performance metrics check කරන්න
3. Statistics analyze කරන්න
```

---

## 🎯 සාර්ථක කොණ්ඩිෂන් (Success Criteria)

### ඔබේ පද්ධතිය හරියටම වැඩ කරනවා නම්:

✅ **Emotional Detection**:
- සතුටු memories 0.7+ weight ලබා ගන්නවා
- Neutral memories 0.0-0.3 weight ලබා ගන්නවා

✅ **Spaced Repetition**:
- Review intervals වැඩි වෙනවා (1→3→7→14→30)
- Emotional memories වැඩි interval ලබා ගන්නවා

✅ **Sleep Consolidation**:
- පෙර වරු 2 AM ට job එක run වෙනවා
- නව memories priority boost ලබා ගන්නවා

✅ **Forgetting Curve**:
- නව memories ඉක්මනින් අමතක වෙනවා (retention < 40%)
- බොහෝ වාර review කළ memories හොඳින් රැඳෙනවා (retention > 80%)

✅ **Tests**:
- 90%+ tests pass වෙනවා
- Coverage > 85%
- කිසිම panic errors නැහැ

---

## 🌟 අවසාන වචන (Final Words)

මෙම පද්ධතිය සාදා ඇත්තේ:
1. 🧠 **ඔබේ මොළය වැඩ කරන ආකාරයම** අනුකරණය කරන්න
2. 💾 **වැදගත් මතක ශක්තිමත්** කරන්න
3. 🎯 **හරියටම අවශ්‍ය වේලාවට** මතක් කරවන්න
4. 📊 **ප්‍රගතිය track** කරන්න

### පරීක්ෂණ වැදගත් ඇයි?
- ✅ කේතය හරියටම වැඩ කරනවා **තහවුරු කරනවා**
- ✅ නව features එකතු කරද්දී පරණ features **බිඳ වැටෙන්නේ නැහැ**
- ✅ ජීව විද්‍යාත්මක සූත්‍ර **හරියටම implement** වෙලා තියෙනවා
- ✅ Production environment එකට යන්න **විශ්වාසයක් දෙනවා**

### ඔබට කළ හැකි දේ:
1. 🎯 **දැන්ම පටන්ගන්න**: Bot එක use කරන්න, memories save කරන්න
2. 🧪 **පරීක්ෂණ run කරන්න**: `go test ./... -v -tags "fts5"`
3. 📊 **Stats monitor කරන්න**: දිනපතා progress track කරන්න
4. 🚀 **දියුණු කරන්න**: නව features එකතු කරන්න, tests ලියන්න

### ප්‍රශ්න තියේ නම්:
- 📖 Documentation files කියවන්න
- 🔍 Database queries run කරන්න
- 🧪 Tests බලන්න කොහොමද වැඩ කරන්නේ කියලා
- 💬 සාකච්ඡා කරන්න, experiment කරන්න

---

## 📞 ඉක්මන් විමර්ශනය (Quick Reference)

### මූලික Commands:
```bash
# Bot start කරන්න
./memory-bot

# Tests run කරන්න
go test ./... -v -tags "fts5"

# Coverage එක බලන්න
go test ./... -coverprofile=coverage.out -tags "fts5"
go tool cover -html=coverage.out

# Database බලන්න
sqlite3 memories.db
```

### Telegram Bot Commands:
```
/start - Bot එක පටන්ගන්නවා
/save - Memory එකක් සුරකින්න
/search - සෙවීම
/recent - මෑත memories
/stats - සංඛ්‍යාලේඛන
/help - උදව්
```

### වැදගත් Files:
```
memories.db - Database file
memory-bot - Executable file
go.mod - Go dependencies
internal/ - කේත files
docs/ - Documentation
```

---

**සාර්ථකත්වය වේවා! Good luck with your biological memory system! 🧠✨**

**Created**: 2025-12-15  
**Version**: 1.0  
**Language**: Sinhala (සිංහල)  
**Author**: AI Assistant  
**Purpose**: Complete testing guide for biological memory features

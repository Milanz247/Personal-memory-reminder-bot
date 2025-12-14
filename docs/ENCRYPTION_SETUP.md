# Encryption Setup Guide

## Overview

Your Personal Memory Bank now supports **AES-256-GCM encryption** to protect your sensitive memory data. All memory content is encrypted before being stored in the database and automatically decrypted when you retrieve it.

## Security Features

- **AES-256 encryption**: Industry-standard encryption with 256-bit keys
- **GCM mode**: Provides both confidentiality and authenticity
- **Per-memory encryption**: Each memory is encrypted with a unique nonce
- **Base64 encoding**: Encrypted data is safely stored as text

## How to Enable Encryption

### 1. Generate a Strong Encryption Key

Generate a random encryption key (minimum 16 characters, but longer is better):

```bash
# Linux/Mac - Generate a 32-character random key
openssl rand -base64 32

# Or use this simple method
echo "your-very-secure-password-$(date +%s)" | sha256sum | cut -d' ' -f1
```

### 2. Set the Environment Variable

Add your encryption key to your `.env` file:

```bash
ENCRYPTION_KEY=your-generated-key-here
```

**Example `.env` file:**

```env
TELEGRAM_BOT_TOKEN=your_telegram_bot_token_here
DB_PATH=./memories.db
ENCRYPTION_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4y5z6
```

### 3. Restart Your Bot

```bash
./stop.sh
./run.sh
```

When encryption is enabled, you'll see:
```
🔒 Encryption enabled for sensitive memory data
```

When encryption is disabled:
```
⚠️  Warning: Encryption is disabled. Set ENCRYPTION_KEY environment variable to enable encryption.
```

## Important Notes

### ⚠️ Encryption Key Management

1. **Keep your key safe**: Store it securely and never share it
2. **Don't lose your key**: If you lose your encryption key, **you cannot recover your encrypted memories**
3. **Key rotation**: Changing the encryption key will make old memories unreadable
4. **Backup**: Keep a secure backup of your encryption key

### Backward Compatibility

- **Without encryption key**: The bot works normally but stores data in plain text
- **With encryption key**: All new memories are encrypted; old plain text memories remain readable

### What Gets Encrypted

✅ **Encrypted:**
- Memory content (text_content)

❌ **Not Encrypted:**
- User IDs
- Chat IDs
- Tags
- Timestamps
- Metadata

## Migration from Unencrypted Database

If you already have memories stored without encryption, they will continue to work. New memories will be encrypted once you set the `ENCRYPTION_KEY`.

To encrypt existing memories:
1. Export your memories
2. Clear the database
3. Set the `ENCRYPTION_KEY`
4. Re-import your memories

## Testing Encryption

### Test 1: Save and Retrieve

```bash
# Save a memory
/save This is a secret message

# Search for it
/search secret

# View recent memories
/recent
```

Your memory should be returned decrypted and readable.

### Test 2: Check Database

Look at the raw database to verify encryption:

```bash
sqlite3 memories.db "SELECT id, substr(text_content, 1, 50) FROM memories LIMIT 5;"
```

You should see base64-encoded gibberish instead of plain text.

## Troubleshooting

### "Failed to decrypt content" Error

**Cause**: Wrong encryption key or corrupted data

**Solution**:
1. Verify your `ENCRYPTION_KEY` is correct
2. If you changed the key, restore the old key
3. Check database integrity: `sqlite3 memories.db "PRAGMA integrity_check;"`

### Performance Considerations

Encryption adds minimal overhead:
- Encryption time: ~1-2ms per memory
- Decryption time: ~1-2ms per memory
- No impact on search speed (tags are not encrypted)

## Security Best Practices

1. **Use environment variables**: Never hardcode keys in your code
2. **Strong keys**: Use at least 32 random characters
3. **Regular backups**: Keep encrypted backups of your database
4. **Secure storage**: Store your `.env` file with restricted permissions

```bash
chmod 600 .env
```

5. **Don't commit keys**: Add `.env` to `.gitignore`

## Technical Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key derivation**: SHA-256 hash of your encryption key
- **Nonce**: 12 bytes, randomly generated per encryption
- **Authentication**: Built-in with GCM mode
- **Encoding**: Base64 for safe database storage

## Example Use Case

```bash
# Set encryption key
echo 'ENCRYPTION_KEY=my-super-secret-key-2024' >> .env

# Start bot
./run.sh

# Your memories are now encrypted! 🔒
```

---

**Remember**: Your encryption key is the only way to decrypt your memories. Store it securely! 🔐

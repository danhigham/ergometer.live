# Troubleshooting Guide

## Firebase: Error (auth/configuration-not-found)

This error means Google Sign-In is not properly configured in Firebase Console.

### Solution Steps

#### 1. Enable Google Authentication

1. Go to https://console.firebase.google.com
2. Select your project: **ergometer-live**
3. Click **Authentication** in the left sidebar
4. Click **Get Started** (if you see it)
5. Click the **Sign-in method** tab
6. Find **Google** in the providers list
7. Click **Google**
8. Toggle **Enable** to ON
9. Set a **Project support email** (required)
10. Click **Save**

#### 2. Add Authorized Domain

1. Still in **Authentication** → **Settings** tab
2. Scroll to **Authorized domains**
3. Make sure `localhost` is in the list
4. If not, click **Add domain** and add `localhost`

#### 3. Verify Your Configuration

In your browser console, you should see:

```
Firebase Config: {
  apiKey: 'AIzaSyDpRw...',
  authDomain: 'ergometer-live.firebaseapp.com',
  projectId: 'ergometer-live'
}
```

If you see any "MISSING" values, your `.env` file is not being loaded correctly.

#### 4. Restart Vite Dev Server

After any changes to `.env`:

```bash
# Stop and restart
./stop-dev.sh
./start-dev.sh

# Or manually restart the UI pane in tmux
# Switch to UI pane: Ctrl+b then arrow keys
# Stop: Ctrl+C
# Start: npm run dev
```

### Common Issues

#### Environment Variables Not Loading

**Symptom**: Console shows `MISSING` for Firebase config values

**Solution**:
1. Verify `.env` file is in `ui/` directory
2. Check file has no quotes around values:
   ```bash
   # WRONG
   VITE_FIREBASE_API_KEY="AIza..."

   # CORRECT
   VITE_FIREBASE_API_KEY=AIza...
   ```
3. Restart Vite (it only reads `.env` on startup)

#### Wrong Project Selected

**Symptom**: Everything configured but still get errors

**Solution**:
1. Double-check you're in the right Firebase project
2. Verify `projectId` in `.env` matches Firebase Console
3. Make sure Web App is created in Firebase Console:
   - Project Settings → General → Your apps
   - Should have a Web app registered

#### Auth Domain Mismatch

**Symptom**: Popup closes immediately or shows error

**Solution**:
1. Check `authDomain` in `.env` exactly matches Firebase
2. Format should be: `your-project-id.firebaseapp.com`
3. No `https://` prefix needed

### Testing Local Mode

If you want to skip authentication temporarily:

1. Go to http://localhost:5173/login
2. Click **"Continue in Local Mode"** instead of Google sign-in
3. This bypasses Firebase entirely and stores data in browser

### Debug Checklist

- [ ] Google auth enabled in Firebase Console
- [ ] Support email set in Google auth config
- [ ] `localhost` in authorized domains
- [ ] `.env` file exists in `ui/` directory
- [ ] No quotes around values in `.env`
- [ ] Vite dev server restarted after `.env` changes
- [ ] Console shows Firebase config (not MISSING)
- [ ] Browser console shows no other errors
- [ ] Correct Firebase project selected

### Still Not Working?

1. Check browser console for full error message
2. Check Firebase Console → Authentication → Users (should be empty but page should load)
3. Try incognito/private browsing mode
4. Clear browser cache and localStorage
5. Try "Continue in Local Mode" as workaround

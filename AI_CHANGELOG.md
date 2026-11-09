# AI Change Log — usersAppAK

Keeps the latest 15 changes made with Claude (oldest are removed). Newest first.

---

## 26. 2026-11-08 — AI knowledge: "Request sent — pending approval" pop-up

**Backups:** `.ai_backups/2026-11-08_before_request_sent_popup/`.

**Description:** works with showcase #30 and shao #43. After every request the company must approve, the store shows
the pop-up (not confirmed yet, check your email, contact the company). Version "2026-11-08 (usersAppAK #26)".
AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 25. 2026-11-08 — AI knowledge: store booking window, contact fields at the end

**Backups:** `.ai_backups/2026-11-08_before_booking_info_last/`.

**Description:** works with showcase #29. Same text as shao #42. Version "2026-11-08 (usersAppAK #25)".
AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 24. 2026-11-08 — AI knowledge: location in requests and appointments, store booking window

**Backups:** `.ai_backups/2026-11-08_before_appointment_location/`.

**Description:** works with shao #41, p3 #51 and showcase #28. Booking requests and Appointments can choose or change
the location ("Save location"); the store window scrolls to the times, shows "Your choice" and has a ↑ button.
Version "2026-11-08 (usersAppAK #24)". AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 23. 2026-11-08 — AI knowledge: past dates, today selected, preview per location

**Backups:** `.ai_backups/2026-11-08_before_past_dates_preview/`.

**Description:** works with shao #40 and p3 #50. Booking schedules: past dates are history, the calendar opens on today
or the first open date, the preview shows one row per location. Version "2026-11-08 (usersAppAK #23)".
AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 22. 2026-11-08 — AI knowledge: Booking schedules in steps, one mode, toasts

**Backups:** `.ai_backups/2026-11-08_before_schedule_steps/`.

**Description:** works with shao #39 and p3 #49. `bAISystemKnowledge.go`: Booking schedules starts with the service
(the rest is locked until then), uses one mode only (changing it deletes the other mode's hours when saved; booked
appointments are kept), shows toasts, and refuses two ranges in the same location at the same time. Version
"2026-11-08 (usersAppAK #22)". AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 21. 2026-11-08 — AI knowledge: "Customer chooses" places per location, 1-hour new ranges

**Backups:** `.ai_backups/2026-11-08_before_places_per_location/`.

**Description:** works with shao #38 and p3 #48. `bAISystemKnowledge.go`: with "same time" ON each location has the
block's places (OFF = shared); a new time range lasts one appointment; saving says which ranges overlap or are too
short. Version "2026-11-08 (usersAppAK #21)". AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 20. 2026-11-08 — AI knowledge: Booking schedules "Create your locations first" box

**Backups:** `.ai_backups/2026-11-08_before_locations_hint/`.

**Description:** works with p3 #47. `bAISystemKnowledge.go`: Booking schedules explains the box with the button
"Go to Locations"; version "2026-11-08 (usersAppAK #20)". AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 19. 2026-11-08 — AI agent knowledge: "We use locations" + "Customer chooses"

**Backups:** `.ai_backups/2026-11-08_before_locations_options/`.

**Description:** works with shao #35 and p3 #46. `brules/bAIAgent/bAISystemKnowledge.go`: the Locations screen now
explains the main switch "We use locations" (off by default = no locations anywhere; saved locations are kept), the
"same time" switch (only when locations are on), the option "Customer chooses (any active location)" on each block
of hours, and that turning locations on is refused while hours clash. Version "2026-11-08 (usersAppAK #19)".
AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Undo:** restore the `.bak` file from the backup folder.

---

## 18. 2026-11-07 — AI agent knowledge: Locations ("Sedes")

**Backups:** `.ai_backups/2026-11-07_before_company_locations/`.

**Description:** works with shao #34 and p3 #45. `brules/bAIAgent/bAISystemKnowledge.go`: new screen "Locations
[/shao/locations]" (places of the company, map point, the switch "more than one location at the same time",
location on hours and appointments, store grouped by location, disabled instead of deleted) and the path in
`aiScreenPaths`. Version "2026-11-07 (usersAppAK #18)". No endpoint or database change. Tests pass. AI knowledge:
updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), rebuild and restart usersAppAK.

---

## 17. 2026-11-04 — AI agent knowledge: clock in the top bar

**Backups:** `.ai_backups/2026-11-04_before_header_clock/`.

**Description:** works with p3 #44. `brules/bAIAgent/bAISystemKnowledge.go`: the Company page line now says the top bar
shows a clock with the company time (and the computer's time when it is in another zone), and the store shows the
company time too. Version "2026-11-04 (usersAppAK #17)". No endpoint or database change. Tests pass. AI knowledge:
updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), rebuild and restart usersAppAK.

---

## 16. 2026-11-04 — AI agent knowledge: company time zone

**Backups:** `.ai_backups/2026-11-04_before_company_time_zone/`.

**Description:** works with shao #32. `brules/bAIAgent/bAISystemKnowledge.go`: the Company page line now says it has
the company **"Time zone"** (default Costa Rica, GMT−6; every hour of the company uses it; changing it does not move
saved appointments: 14:00 stays 14:00). Version "2026-11-04 (usersAppAK #16)". No endpoint, database or tracking
change. Test of the knowledge (version + paths) passes. AI knowledge: updated.

**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), rebuild and restart usersAppAK.

---

## 15. 2026-11-01 — Request logs: request id, JSON logs with levels, panic -> clean 500

**Backups:** `.ai_backups/2026-11-01_before_request_logs/`.

**Description:**
- New `helpers/logger.go` (+ test): `SetupLogger()` in `main.go` after `.env` (same as shao #31: slog, LOG_FORMAT
  json/text, JSON default in production with `severity` / `message`).
- New `routers/request_log.go` (+ test): `requestLogMiddleware` replaces `loggingMiddleware` and `recoveryMiddleware`
  (removed from `routers/router.go`, with `statusRecorder` and the `time` import). One line per call (request_id,
  method, route template, status, duration_ms, bytes) instead of two lines with the raw path, remote address and
  origin; header `X-Request-ID`; panic logged with stack, JSON 500 with `requestId` and CORS headers.
- `env.example.NEW.txt`: `# LOG_FORMAT=`.
- No database change, no endpoint change. AI knowledge: not needed.

**Created:** helpers/logger.go, helpers/logger_test.go, routers/request_log.go, routers/request_log_test.go
**Modified:** main.go, routers/router.go, env.example.NEW.txt

**Revert:** copy the `.bak` files back ("__" = "/"), delete the created files, rebuild and restart usersAppAK.

---

## 14. 2026-11-01 — AI agent knowledge updated to today's menu, version + path test

**Backups:** `.ai_backups/2026-11-01_before_ai_knowledge_update/`.

**Description:**
- `brules/bAIAgent/bAISystemKnowledge.go`: `aiSystemKnowledge` used the old screen names and paths (Start an order,
  /shao/approvala, /shao/oporders...). Now it has today's menu (Sell, Bookings, Customers, My store, Reports, Team)
  with the new paths, the main new features (Showcase link, Your plan, appointment products, Notify the customer per
  step, customer service email, Contact us buttons, AlephK general email) and sign-in with email, Google or Microsoft.
- New `aiKnowledgeVersion` ("2026-11-01 (usersAppAK #14)") sent to Claude; new `aiScreenPaths` (must match p3
  `assistantRoutes.js`) + `bAISystemKnowledge_test.go` (`TestAIKnowledgePaths`, `TestAIKnowledgeVersion`).
- Rule: same as shao #30 (update the AI knowledge in every change that touches screens or features).
- No endpoint, database or tracking change. AI knowledge: updated.

**Created:** brules/bAIAgent/bAISystemKnowledge_test.go
**Modified:** brules/bAIAgent/bAISystemKnowledge.go

**Revert:** copy the `.bak` file back ("__" = "/"), delete the test file, rebuild and restart usersAppAK.

---

## 13. 2026-10-31 — Domains from .env (aleph-k.com, aleph-k.online): CORS, login cookie, CAPTCHA

**Backups:** `.ai_backups/2026-10-31_before_domains/`.

**Description:** new `helpers/allowed_domains.go` (+ test), same rules as shao. `ALLOWED_DOMAINS=aleph-k.com,aleph-k.online`
= each domain and all its subdomains, http and https:
- CORS (`routers/cors_config.go`, `routers/router.go`) accepts them; domain lines allowed in the origins file.
- Login / logout cookie (`HuserLogin.go`, `securityGen/HuserLogout.go`): Domain = the allowed domain that contains the
  API host (`auth.aleph-k.com` → `aleph-k.com`; `AUTH_COOKIE_DOMAIN` may list several and wins); localhost / LAN =
  host-only as before. Secure when the browser used https (`AUTH_COOKIE_SECURE=true|false` forces it), so login works on
  http and https. Before, a set `AUTH_COOKIE_DOMAIN` always made the cookie Secure (login broken on http) and one value
  could not serve two domains.
- CAPTCHA (`brules/bCaptcha`): `*.domain` entries; empty + production = only `ALLOWED_DOMAINS`.
- Docs: `env.example.NEW.txt`, `cors_allowed_origins.example.txt`.

**Created:** helpers/allowed_domains.go, helpers/allowed_domains_test.go
**Modified:** routers/cors_config.go, routers/router.go, brules/bCaptcha/bCaptchaVerify.go,
handlers/HuserLogin/HuserLogin.go, handlers/securityGen/HuserLogout.go, env.example.NEW.txt, cors_allowed_origins.example.txt

**Revert:** copy the `.bak` files back ("__" = "/"), delete the created files, rebuild usersAppAK.

---

## 12. 2026-10-29 — Login keeps the real sign-in provider (fixes "external provider" on change email/password)

**Backups:** `.ai_backups/2026-10-29_before_provider_fix/`.

**Description:** on every login, when the user already existed, `HuserLogin` saved `ProviderID` and `PhotoURL` from
claims named "ProviderID" / "PhotoURL", which do not exist in the Firebase token, so both were saved **empty**. Change
email / change password check `ProviderID == "password"`, so an email+password user was told the account "uses an
external provider". Now the provider comes from the token (`firebase.sign_in_provider`, e.g. `password`,
`google.com`) and the photo from the `picture` claim; when the token has no value the saved one is kept. Users already
affected are fixed automatically at their next login (sign out and sign in again). No DB change.

**Modified:** handlers/HuserLogin/HuserLogin.go

**Revert:** copy the `.bak` file back, rebuild and restart usersAppAK.

---
